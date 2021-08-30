package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
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
	grpcPort uint
	httpPort uint
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

	server := grpc.NewServer()
	api.RegisterPurchaseServiceServer(server, NewPurchaseServer(&log))

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

func NewServer(grpcPort uint, httpPort uint) Server {
	return &server{
		grpcPort: grpcPort,
		httpPort: httpPort,
	}
}
