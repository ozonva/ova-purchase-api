package app

import (
	"fmt"
	"github.com/ozonva/ova-purchase-api/internal/config"
	db2 "github.com/ozonva/ova-purchase-api/internal/db"
	"github.com/ozonva/ova-purchase-api/internal/disposal"
	"github.com/ozonva/ova-purchase-api/internal/kafka"
	metrics "github.com/ozonva/ova-purchase-api/internal/metric"
	"github.com/ozonva/ova-purchase-api/internal/repo"
	server2 "github.com/ozonva/ova-purchase-api/internal/server"
	"github.com/ozonva/ova-purchase-api/internal/tracer"
	"github.com/rs/zerolog/log"
	"sync"
)

type App interface {
	Start()
	Stop()
}

type application struct {
	configuration *config.Configuration
	disposals     []disposal.Disposal
}

func NewApp(config *config.Configuration) App {
	return &application{
		configuration: config,
		disposals:     make([]disposal.Disposal, 0),
	}
}

func (s *application) Start() {
	log.Info().Msg("Starting application")

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		s.configuration.Db.Username,
		s.configuration.Db.Password,
		s.configuration.Db.Host,
		s.configuration.Db.Port,
		s.configuration.Db.Name,
	)
	wg := &sync.WaitGroup{}
	wg.Add(3)

	db, err := db2.NewDB(url)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create connection to db")
	}

	tracer, err := tracer.NewTracer(*s.configuration)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to tracer")
	}

	metrics := metrics.NewMetrics("purchase", "counts")
	producer, err := kafka.NewProducer(*s.configuration.Kafka)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to kafka")
	}
	purchaseServer := server2.NewPurchaseServer(repo.NewRepo(db.Db), metrics, producer, s.configuration.Batch.Size)

	grpc := server2.NewGrpcServer(wg, *s.configuration.Grpc, purchaseServer)
	gateway := server2.NewGatewayServer(wg, *s.configuration.Gateway, *s.configuration.Grpc)
	metric := server2.NewMetricServer(wg, *s.configuration.Prometheus)

	s.disposals = append(s.disposals, grpc, gateway, metric, db, tracer)

	go func() {
		if err := gateway.Run(); err != nil {
			log.Error().Err(err).Msg("Failed run gateway server")
		}
	}()
	go func() {
		if err := grpc.Run(); err != nil {
			log.Error().Err(err).Msg("Failed run GRPC server")
		}
	}()
	go func() {
		if err := metric.Run(); err != nil {
			log.Error().Err(err).Msg("Failed run metrics server")
		}
	}()
}

func (s *application) Stop() {
	log.Info().Msg("Stopping application")
	for _, d := range s.disposals {
		d.Disposal()
	}
}
