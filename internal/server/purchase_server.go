package server

import (
	"context"
	api "github.com/ozonva/ova-purchase-api/pkg/ova-purchase-api"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PurchaseServer struct {
	api.UnimplementedPurchaseServiceServer
}

func NewPurchaseServer() *PurchaseServer {
	return &PurchaseServer{
		UnimplementedPurchaseServiceServer: api.UnimplementedPurchaseServiceServer{},
	}
}

func (s *PurchaseServer) CreatePurchase(context context.Context, request *api.CreatePurchaseRequest) (*api.CreatePurchaseResponse, error) {
	return s.UnimplementedPurchaseServiceServer.CreatePurchase(context, request)
}

func (s *PurchaseServer) ListPurchases(context context.Context, request *api.ListPurchasesRequest) (*api.ListPurchasesResponse, error) {
	return s.UnimplementedPurchaseServiceServer.ListPurchases(context, request)
}

func (s *PurchaseServer) DescribePurchase(context context.Context, request *api.DescribePurchaseRequest) (*api.DescribePurchaseResponse, error) {
	return s.UnimplementedPurchaseServiceServer.DescribePurchase(context, request)
}

func (s *PurchaseServer) RemovePurchase(context context.Context, request *api.RemovePurchaseRequest) (*emptypb.Empty, error) {
	return s.UnimplementedPurchaseServiceServer.RemovePurchase(context, request)
}
