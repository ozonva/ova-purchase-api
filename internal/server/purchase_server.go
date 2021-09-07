package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/ozonva/ova-purchase-api/internal/kafka"
	metrics "github.com/ozonva/ova-purchase-api/internal/metric"
	"github.com/ozonva/ova-purchase-api/internal/purchase"
	"github.com/ozonva/ova-purchase-api/internal/repo"
	"github.com/ozonva/ova-purchase-api/internal/utils"
	api "github.com/ozonva/ova-purchase-api/pkg/ova-purchase-api"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PurchaseServer struct {
	repo      repo.Repo
	metrics   metrics.Metrics
	producer  kafka.KafkaProducer
	batchSize uint
	api.UnimplementedPurchaseServiceServer
}

func NewPurchaseServer(repo repo.Repo, metrics metrics.Metrics, producer kafka.KafkaProducer, batchSize uint) *PurchaseServer {
	return &PurchaseServer{
		repo:                               repo,
		producer:                           producer,
		metrics:                            metrics,
		batchSize:                          batchSize,
		UnimplementedPurchaseServiceServer: api.UnimplementedPurchaseServiceServer{},
	}
}

func (s *PurchaseServer) MultiCreatePurchases(ctx context.Context, request *api.MultiCreatePurchaseRequest) (*api.MultiCreatePurchaseResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	purchases := make([]purchase.Purchase, 0, len(request.Purchases))
	for _, v := range request.GetPurchases() {
		p := purchase.New()
		for _, item := range v.Items {
			_, err := p.Add(purchase.Item{
				Price:    decimal.NewFromFloat(item.Price),
				Name:     item.Name,
				Quantity: uint(item.Quantity),
			})
			if err != nil {
				return nil, err
			}
		}
		purchases = append(purchases, p)
	}
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("MultiCreatePurchases")
	defer span.Finish()

	batches, err := utils.SplitToBulks(purchases, uint(s.batchSize))
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(purchases))
	for _, batch := range batches {
		added, err := s.repo.AddPurchases(ctx, batch)
		if err != nil {
			return nil, err
		}
		ids = append(ids, added...)

		childSpan := tracer.StartSpan(
			"MultiCreatePurchases: batch",
			opentracing.Tag{Key: "batchSize", Value: len(batch)},
			opentracing.ChildOf(span.Context()),
		)
		childSpan.Finish()
	}

	err = s.producer.Send(kafka.NewMessage(kafka.MultiCreatePurchase, ids))
	if err != nil {
		return nil, err
	}

	s.metrics.MultiCreatePurchaseCounterInc()

	return &api.MultiCreatePurchaseResponse{
		Ids: ids,
	}, nil
}

func (s *PurchaseServer) UpdatePurchase(ctx context.Context, request *api.UpdatePurchaseRequest) (*api.DescribePurchaseResponse, error) {
	p, err := s.repo.DescribePurchase(ctx, request.Id)
	if err != nil {
		if errors.Is(err, repo.PurchaseNotFoundError) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
	}
	mapItems := make(map[uint64]purchase.Item)
	for _, item := range p.Items {
		mapItems[item.Id] = item
	}
	for _, item := range request.Items {
		if exist, ok := mapItems[item.Id]; ok {
			if _, err = p.Remove(item.Id); err != nil {
				return nil, err
			}
			if item.Quantity == 0 {
				continue
			}
			exist.Name = item.Name
			exist.Price = decimal.NewFromFloat(item.Price)
			exist.Quantity = uint(item.Quantity)
			if _, err = p.Add(exist); err != nil {
				return nil, err
			}
		} else {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("Item with id = %d not found in this purchases", item.Id))
		}
	}
	if err = s.repo.UpdatePurchase(ctx, *p); err != nil {
		return nil, err
	}
	total, _ := p.Total.Float64()
	items := make([]*api.DescribePurchaseResponse_Item, 0)
	for _, item := range p.Items {
		price, _ := item.Price.Float64()
		items = append(items, &api.DescribePurchaseResponse_Item{
			Id:       item.Id,
			Name:     item.Name,
			Price:    price,
			Quantity: uint32(item.Quantity),
		})
	}
	err = s.producer.Send(kafka.NewMessage(kafka.UpdatePurchase, *p))

	if err != nil {
		log.Error().Err(err).Msg("Failed to send event to Kafka")
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.metrics.UpdatePurchaseCounterInc()
	return &api.DescribePurchaseResponse{
		Id:        p.Id,
		Total:     total,
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
		Status:    api.PurchaseStatus(p.Status),
		Items:     items,
	}, nil
}

func (s *PurchaseServer) CreatePurchase(context context.Context, request *api.CreatePurchaseRequest) (*api.DescribePurchaseResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	p := purchase.New()
	for _, item := range request.Items {
		_, err := p.Add(purchase.Item{
			Price:    decimal.NewFromFloat(item.Price),
			Name:     item.Name,
			Quantity: uint(item.Quantity),
		})
		if err != nil {
			return nil, err
		}
	}
	id, err := s.repo.AddPurchase(context, p)
	if err != nil {
		return nil, err
	}
	created, err := s.repo.DescribePurchase(context, id)
	if err != nil {
		log.Debug().Msgf("ID %d not found", id)
		return nil, err
	}
	total, _ := created.Total.Float64()
	items := make([]*api.DescribePurchaseResponse_Item, 0)
	for _, item := range created.Items {
		price, _ := item.Price.Float64()
		items = append(items, &api.DescribePurchaseResponse_Item{
			Id:       item.Id,
			Name:     item.Name,
			Price:    price,
			Quantity: uint32(item.Quantity),
		})
	}
	err = s.producer.Send(kafka.NewMessage(kafka.CreatePurchase, *created))
	if err != nil {
		log.Error().Err(err).Msg("Failed to send event to Kafka")
		return nil, status.Error(codes.Internal, err.Error())
	}
	s.metrics.CreatePurchaseCounterInc()
	return &api.DescribePurchaseResponse{
		Id:        id,
		Total:     total,
		CreatedAt: timestamppb.New(created.CreatedAt),
		UpdatedAt: timestamppb.New(created.UpdatedAt),
		Status:    api.PurchaseStatus(created.Status),
		Items:     items,
	}, nil
}

func (s *PurchaseServer) ListPurchases(context context.Context, request *api.ListPurchasesRequest) (*api.ListPurchasesResponse, error) {
	limit := uint(request.Limit)
	if limit <= 0 {
		limit = 20
	}
	list, err := s.repo.ListPurchases(context, limit, uint(request.Offset))
	if err != nil {
		return nil, err
	}
	purchases := make([]*api.Purchase, 0, len(list))
	for _, p := range list {
		total, _ := p.Total.Float64()
		purchaseItems := make([]*api.Purchase_Item, 0)
		for _, item := range p.Items {
			price, _ := item.Price.Float64()
			purchaseItems = append(purchaseItems, &api.Purchase_Item{
				Id:       item.Id,
				Name:     item.Name,
				Price:    price,
				Quantity: uint32(item.Quantity),
			})
		}
		purchases = append(purchases, &api.Purchase{
			Id:        p.Id,
			Total:     total,
			CreatedAt: timestamppb.New(p.CreatedAt),
			UpdatedAt: timestamppb.New(p.UpdatedAt),
			Status:    api.PurchaseStatus(p.Status),
			Items:     purchaseItems,
		})
	}
	return &api.ListPurchasesResponse{
		Purchases: purchases,
	}, nil
}

func (s *PurchaseServer) DescribePurchase(context context.Context, request *api.DescribePurchaseRequest) (*api.DescribePurchaseResponse, error) {
	p, err := s.repo.DescribePurchase(context, request.Id)
	if err != nil {
		if errors.Is(err, repo.PurchaseNotFoundError) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	total, _ := p.Total.Float64()
	items := make([]*api.DescribePurchaseResponse_Item, 0, len(p.Items))
	for _, item := range p.Items {
		price, _ := item.Price.Float64()
		items = append(items, &api.DescribePurchaseResponse_Item{
			Id:       item.Id,
			Name:     item.Name,
			Price:    price,
			Quantity: uint32(item.Quantity),
		})
	}
	return &api.DescribePurchaseResponse{
		Id:        p.Id,
		Total:     total,
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
		Status:    api.PurchaseStatus(p.Status),
		Items:     items,
	}, nil
}

func (s *PurchaseServer) RemovePurchase(context context.Context, request *api.RemovePurchaseRequest) (*emptypb.Empty, error) {
	p, err := s.repo.DescribePurchase(context, request.Id)
	if err != nil {
		if errors.Is(err, repo.PurchaseNotFoundError) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	err = s.repo.RemovePurchase(context, request.Id)
	if err != nil {
		if errors.Is(err, repo.PurchaseNotFoundError) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	err = s.producer.Send(kafka.NewMessage(kafka.RemovePurchase, *p))
	if err != nil {
		log.Error().Err(err).Msg("Failed to send event to Kafka")
		return nil, status.Error(codes.Internal, err.Error())
	}
	s.metrics.RemovePurchaseCounterInc()
	return &emptypb.Empty{}, nil
}
