package server

import (
	"context"
	api "github.com/ozonva/ova-purchase-api/pkg/ova-purchase-api"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PurchaseServer struct {
	api.UnimplementedPurchaseServiceServer
	logger *zerolog.Logger
}

func NewPurchaseServer(logger *zerolog.Logger) *PurchaseServer {
	return &PurchaseServer{
		logger:                             logger,
		UnimplementedPurchaseServiceServer: api.UnimplementedPurchaseServiceServer{},
	}
}

func (s *PurchaseServer) CreatePurchase(context context.Context, request *api.CreatePurchaseRequest) (*api.CreatePurchaseResponse, error) {
	s.logger.Info().Msgf("CreatePurchase request: %v", request)
	return s.UnimplementedPurchaseServiceServer.CreatePurchase(context, request)
}

func (s *PurchaseServer) ListPurchases(context context.Context, request *api.ListPurchasesRequest) (*api.ListPurchasesResponse, error) {
	s.logger.Info().Msgf("ListPurchases request: %v", request)
	return s.UnimplementedPurchaseServiceServer.ListPurchases(context, request)
}

func (s *PurchaseServer) DescribePurchase(context context.Context, request *api.DescribePurchaseRequest) (*api.DescribePurchaseResponse, error) {
	s.logger.Info().Msgf("DescribePurchase request: %v", request)
	return s.UnimplementedPurchaseServiceServer.DescribePurchase(context, request)
}

func (s *PurchaseServer) RemovePurchase(context context.Context, request *api.RemovePurchaseRequest) (*emptypb.Empty, error) {
	s.logger.Info().Msgf("RemovePurchase request: %v", request)
	return s.UnimplementedPurchaseServiceServer.RemovePurchase(context, request)
}
