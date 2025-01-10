package metrics

import (
	"TOomaAh/emby_exporter_go/pkg/emby"
	"TOomaAh/emby_exporter_go/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

type RegisterFunc func(server *emby.Server, logger logger.Interface) prometheus.Collector

type Register struct {
	registry   *prometheus.Registry
	server     *emby.Server
	logger     logger.Interface
	collectors []prometheus.Collector
}

func NewRegister(server *emby.Server, logger logger.Interface) *Register {
	return &Register{
		server:   server,
		logger:   logger,
		registry: prometheus.NewRegistry(),
	}
}

func (r *Register) RegisterCollector(method ...RegisterFunc) {
	for _, m := range method {
		collector := m(r.server, r.logger)
		r.collectors = append(r.collectors, collector)
	}
}

func (r *Register) Run() {
	for _, c := range r.collectors {
		r.registry.MustRegister(c)
	}
}

func (r *Register) Gather() ([]*io_prometheus_client.MetricFamily, error) {
	for _, c := range r.collectors {
		if dto, err := c.(prometheus.Gatherer).Gather(); err == nil {
			return dto, nil
		} else {
			return nil, err
		}
	}
	return nil, nil
}
