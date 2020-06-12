package metrics

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type Metrics struct {
	namespace string
	metrics   *sync.Map
}

func CreateMetrics(namespace string) Metrics {
	return Metrics{namespace: namespace}
}

func createKey(name string, label map[string]string) string {
	bytes, err := json.Marshal(label)
	if err != nil {
		return name
	}
	return name + string(bytes)
}

func (m *Metrics) Counter(name string, label map[string]string) prometheus.Counter {
	key := createKey(name, label)
	if value, ok := m.metrics.Load(key); ok {
		switch value := value.(type) {
		case prometheus.Counter:
			return value
		}
	}
	c := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   m.namespace,
		Name:        name + "_Counter",
		Help:        name + " help",
		ConstLabels: label,
	})
	m.metrics.Store(key, &c)
	return c
}

func (m *Metrics) Gauge(name string, label map[string]string) prometheus.Gauge {
	key := createKey(name, label)
	if value, ok := m.metrics.Load(key); ok {
		switch value := value.(type) {
		case prometheus.Gauge:
			return value
		}
	}
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   m.namespace,
		Name:        name + "_Gauge",
		Help:        name + " help",
		ConstLabels: label,
	})
	m.metrics.Store(key, &g)
	return g
}

func (m *Metrics) Histogram(name string, label map[string]string) prometheus.Histogram {
	key := createKey(name, label)
	if value, ok := m.metrics.Load(key); ok {
		switch value := value.(type) {
		case prometheus.Histogram:
			return value
		}
	}
	h := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace:   m.namespace,
		Name:        name + "_Histogram",
		Help:        name + " help",
		ConstLabels: label,
	})
	m.metrics.Store(key, &h)
	return h
}

func (m *Metrics) Summary(name string, label map[string]string) prometheus.Summary {
	key := createKey(name, label)
	if value, ok := m.metrics.Load(key); ok {
		switch value := value.(type) {
		case prometheus.Summary:
			return value
		}
	}
	s := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:   m.namespace,
		Name:        name + "_Summary",
		Help:        name + " help",
		ConstLabels: label,
	})
	m.metrics.Store(key, &s)
	return s
}

func (m *Metrics) Describe(ch chan<- *prometheus.Desc) {
	m.metrics.Range(func(_, value interface{}) bool {
		switch value := value.(type) {
		case prometheus.Counter:
			ch <- value.Desc()
		case prometheus.Gauge:
			ch <- value.Desc()
		case prometheus.Histogram:
			ch <- value.Desc()
		case prometheus.Summary:
			ch <- value.Desc()
		}
		return true
	})
}

func (m *Metrics) Collect(ch chan<- prometheus.Metric) {
	m.metrics.Range(func(_, value interface{}) bool {
		switch value := value.(type) {
		case prometheus.Counter:
			value.Collect(ch)
		case prometheus.Gauge:
			value.Collect(ch)
		case prometheus.Histogram:
			value.Collect(ch)
		case prometheus.Summary:
			value.Collect(ch)
		}
		return true
	})
}
