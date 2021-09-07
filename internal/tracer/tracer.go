package tracer

import (
	"github.com/opentracing/opentracing-go"
	"github.com/ozonva/ova-purchase-api/internal/config"
	"github.com/ozonva/ova-purchase-api/internal/disposal"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"io"
)

type tracer struct {
	opentracer opentracing.Tracer
	closer     io.Closer
}

func NewTracer(config config.Configuration) (disposal.Disposal, error) {

	cfg := jaegercfg.Configuration{
		ServiceName: config.Application.Name,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: config.Jaeger.String(),
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	opentracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jLogger), jaegercfg.Metrics(jMetricsFactory))
	if err != nil {
		log.Error().Err(err).Msg("Could not initialize tracer")
		return nil, err
	}

	opentracing.SetGlobalTracer(opentracer)

	return &tracer{
		opentracer: opentracer,
		closer:     closer,
	}, nil
}

func (s *tracer) Disposal() {
	log.Debug().Msg("Closing tracer...")
	if err := s.closer.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close tracer")
	}
}
