package matrixid

import "github.com/prometheus/client_golang/prometheus"

const HandlerLabel = "handler"

type MetricApi struct {
	ErrorCount *prometheus.CounterVec
	enabled    bool
}

func NewMetric() MetricApi {
	var m MetricApi
	m.ErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "matrix",
			Subsystem: "api",
			Name:      "errors",
			Help:      "counts number of errors occuring in a handler",
		},
		[]string{HandlerLabel},
	)
	prometheus.MustRegister(m.ErrorCount)
	return m
}

func (m MetricApi) CountError(handler string) prometheus.Counter {
	return m.ErrorCount.With(prometheus.Labels{
		HandlerLabel: handler,
	})
}

type Metric interface {
	// CountError returns collector for error counts in an handler. By error counts
	// it means an occupance of if err!=nil expression in the handler body. This is
	// to allow high quality code as the API is already final and cover small
	// surface area.
	//
	// We hope to achieve 0 runtime errors for this service deployments.
	CountError(handler string) prometheus.Counter
}
