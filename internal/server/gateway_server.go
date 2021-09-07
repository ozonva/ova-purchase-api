package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ozonva/ova-purchase-api/internal/config"
	api "github.com/ozonva/ova-purchase-api/pkg/ova-purchase-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net/http"
	"sync"
)

type GatewayServer struct {
	wg         *sync.WaitGroup
	httpConfig config.EndpointConfiguration
	grpcConfig config.EndpointConfiguration
	server     http.Server
}

func NewGatewayServer(wg *sync.WaitGroup, httpConfig config.EndpointConfiguration, grpcConfig config.EndpointConfiguration) Server {
	return &GatewayServer{
		wg:         wg,
		httpConfig: httpConfig,
		grpcConfig: grpcConfig,
	}
}

func (s *GatewayServer) Run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	endpoint := s.grpcConfig.String()

	err := api.RegisterPurchaseServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		return err
	}
	s.server = http.Server{
		Addr:    s.httpConfig.String(),
		Handler: mux,
	}
	err = s.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	if err != nil {
		log.Error().Err(err).Msg("Gateway error")
	}
	s.wg.Done()
	return err
}

func (s *GatewayServer) Disposal() {
	log.Debug().Msg("Shutdown gateway server ...")
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("Failed to shutdown gateway server")
	}
}
