package server

import (
	"context"
	"errors"
	"github.com/ozonva/ova-purchase-api/internal/purchase"
	"github.com/ozonva/ova-purchase-api/internal/repo"
	api "github.com/ozonva/ova-purchase-api/pkg/ova-purchase-api"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PurchaseServer struct {
	repo repo.Repo
	api.UnimplementedPurchaseServiceServer
}

func NewPurchaseServer(repo repo.Repo) *PurchaseServer {
	return &PurchaseServer{
		repo:                               repo,
		UnimplementedPurchaseServiceServer: api.UnimplementedPurchaseServiceServer{},
	}
}

func (s *PurchaseServer) CreatePurchase(context context.Context, request *api.CreatePurchaseRequest) (*api.CreatePurchaseResponse, error) {
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
	items := make([]*api.CreatePurchaseResponse_Item, 0)
	for _, item := range created.Items {
		price, _ := item.Price.Float64()
		items = append(items, &api.CreatePurchaseResponse_Item{
			Id:       item.Id,
			Name:     item.Name,
			Price:    price,
			Quantity: uint32(item.Quantity),
		})
	}
	return &api.CreatePurchaseResponse{
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
	err := s.repo.RemovePurchase(context, request.Id)
	if err != nil {
		if errors.Is(err, repo.PurchaseNotFoundError) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
