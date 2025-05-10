package system

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	promClient "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/khusainnov/iam/app/config"
)

const (
	shutdownTimeout = 5 * time.Second
)

type System struct {
	Registry *promClient.Registry

	meterProvider *sdkmetric.MeterProvider
	prometheus    http.Handler

	Mux *http.ServeMux
	srv *http.Server
}

func New(log *zap.Logger, cfg *config.Config) (*System, error) {
	if cfg.Server.Addr == "" {
		cfg.Server.Addr = ":53000"
	}

	registry := promClient.NewPedanticRegistry()
	registry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(
			collectors.WithGoCollections(collectors.GoRuntimeMemStatsCollection|collectors.GoRuntimeMetricsCollection),
		),
		collectors.NewBuildInfoCollector(),
	)

	promExporter, err := prometheus.New(
		prometheus.WithRegisterer(registry),
	)
	if err != nil {
		return nil, fmt.Errorf("prometheus: %w", err)
	}

	metricProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(promExporter),
	)

	mux := http.NewServeMux()

	s := &System{
		Registry: registry,

		meterProvider: metricProvider,
		prometheus:    promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
		Mux:           mux,
		srv: &http.Server{
			Handler: mux,
			Addr:    cfg.Server.Addr,
		},
	}

	otel.SetMeterProvider(s.MeterProvider())
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		),
	)

	s.registerRoot()
	s.registerProfiler()
	s.registerPrometheus()

	log.Info("Metrics initialized",
		zap.String("http.addr", s.srv.Addr),
	)

	return s, nil
}

func (s *System) registerProfiler() {
	// Routes for pprof
	s.Mux.HandleFunc("/debug/pprof/", pprof.Index)
	s.Mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.Mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.Mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.Mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Manually add support for paths linked to by index page at /debug/pprof/
	s.Mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	s.Mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	s.Mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	s.Mux.Handle("/debug/pprof/block", pprof.Handler("block"))
}

func (s *System) registerPrometheus() {
	s.Mux.Handle("/metrics", s.prometheus)
}

func (s *System) MeterProvider() metric.MeterProvider {
	return s.meterProvider
}

func (s *System) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *System) registerRoot() {
	s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var b strings.Builder
		b.WriteString("Service is up and running.\n\n")
		b.WriteString("\nAvailable debug endpoints:\n")
		for _, s := range []struct {
			Name        string
			Description string
		}{
			{"/metrics", "prometheus metrics"},
			{"/debug/pprof/", "exported pprof"},
		} {
			b.WriteString(fmt.Sprintf("%-20s - %s\n", s.Name, s.Description))
		}

		_, _ = fmt.Fprintln(w, b.String())
	})
}

func (s *System) Run(ctx context.Context) error {
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	wg.Go(func() error {
		// Wait until g ctx canceled, then try to shut down server
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return s.Shutdown(ctx)
	})

	return wg.Wait()
}

