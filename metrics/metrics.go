package metrics

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	namespace string
	counter   map[string]prometheus.Counter
	gauge     map[string]prometheus.Gauge
	histogram map[string]prometheus.Histogram
	summary   map[string]prometheus.Summary
}

func CreateMetrics(namespace string) Metrics {
	info := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   namespace,
		Name:        "Metrics_Info",
		Help:        "Metrics_Info help",
		ConstLabels: map[string]string{"version": "1.0.0"},
	})
	return Metrics{namespace: namespace, counter: map[string]prometheus.Counter{"Metrics_Info": info}, gauge: make(map[string]prometheus.Gauge), histogram: make(map[string]prometheus.Histogram), summary: make(map[string]prometheus.Summary)}
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
	if value, ok := m.counter[key]; ok {
		return value
	}
	c := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   m.namespace,
		Name:        name + "_Counter",
		Help:        name + " help",
		ConstLabels: label,
	})
	m.counter[key] = c
	return c
}

func (m Metrics) Gauge(name string, label map[string]string) prometheus.Gauge {
	key := createKey(name, label)
	if value, ok := m.gauge[key]; ok {
		return value
	}
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   m.namespace,
		Name:        name + "_Gauge",
		Help:        name + " help",
		ConstLabels: label,
	})
	m.gauge[key] = g
	return g
}

func (m Metrics) Histogram(name string, label map[string]string) prometheus.Histogram {
	key := createKey(name, label)
	if value, ok := m.histogram[key]; ok {
		return value
	}
	h := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace:   m.namespace,
		Name:        name + "_Histogram",
		Help:        name + " help",
		ConstLabels: label,
	})
	m.histogram[key] = h
	return h
}

func (m Metrics) Summary(name string, label map[string]string) prometheus.Summary {
	key := createKey(name, label)
	if value, ok := m.summary[key]; ok {
		return value
	}
	s := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:   m.namespace,
		Name:        name + "_Summary",
		Help:        name + " help",
		ConstLabels: label,
	})
	m.summary[key] = s
	return s
}

func (m Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, value := range m.counter {
		ch <- value.Desc()
	}
	for _, value := range m.gauge {
		ch <- value.Desc()
	}
	for _, value := range m.histogram {
		ch <- value.Desc()
	}
	for _, value := range m.summary {
		ch <- value.Desc()
	}
}

func (m Metrics) Collect(ch chan<- prometheus.Metric) {
	for _, value := range m.counter {
		value.Collect(ch)
	}
	for _, value := range m.gauge {
		value.Collect(ch)
	}
	for _, value := range m.histogram {
		value.Collect(ch)
	}
	for _, value := range m.summary {
		value.Collect(ch)
	}
}
