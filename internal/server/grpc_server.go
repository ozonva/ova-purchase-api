package server

import (
	"context"
	"fmt"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/ozonva/ova-purchase-api/internal/config"
	api "github.com/ozonva/ova-purchase-api/pkg/ova-purchase-api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"os"
	"sync"
	"time"
)

type GrpcServer struct {
	wg      *sync.WaitGroup
	config  config.EndpointConfiguration
	service api.PurchaseServiceServer
	server  *grpc.Server
}

func NewGrpcServer(wg *sync.WaitGroup, config config.EndpointConfiguration, service api.PurchaseServiceServer) Server {
	return &GrpcServer{
		service: service,
		wg:      wg,
		config:  config,
	}
}

func logUnaryInterceptor(log *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)
		if err != nil {
			log.Error().Time("time", startTime).Err(err).Msgf("Executing endpoint %s", info.FullMethod)
		} else {
			log.Info().Time("time", startTime).Msgf("Executing endpoint %s", info.FullMethod)
		}
		return resp, err
	}
}

func (s *GrpcServer) Run() error {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))

	if err != nil {
		return err
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()

	s.server = grpc.NewServer(grpc.ChainUnaryInterceptor(
		logUnaryInterceptor(&log),
		grpc_recovery.UnaryServerInterceptor(),
	))

	api.RegisterPurchaseServiceServer(s.server, s.service)

	if err := s.server.Serve(listen); err != nil {
		log.Error().Msgf("Error start %v", err)
		return err
	}
	s.wg.Done()
	return nil
}

func (s *GrpcServer) Disposal() {
	log.Debug().Msg("Shutdown GRPC server ...")
	s.server.GracefulStop()
}
