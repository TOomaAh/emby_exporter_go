package series

import "github.com/prometheus/client_golang/prometheus"

type SeriesCollector struct {
	series   SeriesInterface
	history  *prometheus.Desc
	schedule *prometheus.Desc
}

func NewSeriesCollector(series SeriesInterface) *SeriesCollector {
	return &SeriesCollector{
		series:   series,
		history:  prometheus.NewDesc("series_history", "", []string{"name", "season", "status", "airdate"}, nil),
		schedule: prometheus.NewDesc("series_schedule", "", []string{"name", "season", "airdate"}, nil),
	}
}

func (s *SeriesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.history
	ch <- s.schedule
}

func (s *SeriesCollector) Collect(ch chan<- prometheus.Metric) {
	history := s.series.GetHistory()
	schedule := s.series.GetTodayEpisodes()

	for i, h := range *history {
		ch <- prometheus.MustNewConstMetric(s.history, prometheus.GaugeValue, float64(i), h.Name, h.Season, h.Status, h.AirDate)
	}

	for i, sc := range *schedule {
		ch <- prometheus.MustNewConstMetric(s.schedule, prometheus.GaugeValue, float64(i), sc.Name, sc.Season, sc.AirDate)
	}
}
