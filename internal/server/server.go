package server

import (
	"context"
	"fmt"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jmoiron/sqlx"
	"github.com/ozonva/ova-purchase-api/internal/repo"
	api "github.com/ozonva/ova-purchase-api/pkg/ova-purchase-api"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

type Server interface {
	Run()
}

type server struct {
	db       *sqlx.DB
	grpcPort uint
	httpPort uint
}

func logUnaryInterceptor(log *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)
		if err != nil {
			log.Error().Time("time", startTime).Msgf("Executing endpoint %s", info.FullMethod)
		} else {
			log.Info().Time("time", startTime).Msgf("Executing endpoint %s", info.FullMethod)
		}
		return resp, err
	}
}

func (s *server) runHttp(wg *sync.WaitGroup) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := api.RegisterPurchaseServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", s.grpcPort), opts)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", s.httpPort), mux)
	if err != nil {
		panic(err)
	}
	wg.Done()
}

func (s *server) runGrpc(wg *sync.WaitGroup) {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", s.grpcPort))

	if err != nil {
		panic(err)
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logUnaryInterceptor(&log),
		grpc_recovery.UnaryServerInterceptor(),
	))

	api.RegisterPurchaseServiceServer(server, NewPurchaseServer(repo.NewRepo(s.db)))

	if err := server.Serve(listen); err != nil {
		panic(err)
	}
	wg.Done()
}

func (s *server) Run() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go s.runGrpc(&wg)
	go s.runHttp(&wg)
	wg.Wait()
}

func NewServer(db *sqlx.DB, grpcPort uint, httpPort uint) Server {
	return &server{
		db:       db,
		grpcPort: grpcPort,
		httpPort: httpPort,
	}
}
