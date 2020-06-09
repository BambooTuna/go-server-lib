package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	namespace string
	counter   map[string]prometheus.Counter
	gauge     map[string]prometheus.Gauge
	histogram map[string]prometheus.Histogram
	summary   map[string]prometheus.Summary
}

func CreateMetrics(namespace string) Metrics {
	return Metrics{namespace: namespace, counter: make(map[string]prometheus.Counter), gauge: make(map[string]prometheus.Gauge)}
}

func (m *Metrics) Counter(name string, label map[string]string) prometheus.Counter {
	if value, ok := m.counter[name]; ok {
		return value
	}
	c := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   m.namespace,
		Name:        name + "_Counter",
		Help:        name + " help",
		ConstLabels: label,
	})
	m.counter[name] = c
	return c
}

func (m Metrics) Gauge(name string, label map[string]string) prometheus.Gauge {
	if value, ok := m.gauge[name]; ok {
		return value
	}
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   m.namespace,
		Name:        name + "_Gauge",
		Help:        name + " help",
		ConstLabels: label,
	})
	m.gauge[name] = g
	return g
}

func (m Metrics) Histogram(name string, label map[string]string) prometheus.Histogram {
	if value, ok := m.histogram[name]; ok {
		return value
	}
	h := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace:   m.namespace,
		Name:        name + "_Histogram",
		Help:        name + " help",
		ConstLabels: label,
	})
	m.histogram[name] = h
	return h
}

func (m Metrics) Summary(name string, label map[string]string) prometheus.Summary {
	if value, ok := m.summary[name]; ok {
		return value
	}
	s := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:   m.namespace,
		Name:        name + "_Summary",
		Help:        name + " help",
		ConstLabels: label,
	})
	m.summary[name] = s
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
