package server

import (
	"context"
	"fmt"
	"github.com/ozonva/ova-purchase-api/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
)

type MetricsServer struct {
	wg     *sync.WaitGroup
	config config.PrometheusConfiguration
	server *http.Server
}

func NewMetricServer(wg *sync.WaitGroup, config config.PrometheusConfiguration) Server {
	return &MetricsServer{
		wg:     wg,
		config: config,
	}
}

func (s *MetricsServer) Run() error {
	mux := http.NewServeMux()
	mux.Handle(s.config.Path, promhttp.Handler())

	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler: mux,
	}
	err := s.server.ListenAndServe()

	if err == http.ErrServerClosed {
		return nil
	}
	s.wg.Done()
	return err
}

func (s *MetricsServer) Disposal() {
	log.Debug().Msg("Shutdown metrics server ...")
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msgf("Failed to shutdown metrics server")
	}
}
