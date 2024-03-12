package metrics

import "github.com/prometheus/client_golang/prometheus"

type MetricsApp struct {
	CounterReq  *prometheus.CounterVec
	DurationReq *prometheus.HistogramVec
}

func AddMetrics() *MetricsApp {
	countReq := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total_masuk",
			Help: "menghitung total request yang masuk",
		}, []string{"path", "method"})

	durationReq := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_endpoint_duration",
		Help: "durasi setiap enpoint diprocess",
	}, []string{"path", "method"})

	return &MetricsApp{
		CounterReq:  countReq,
		DurationReq: durationReq,
	}
}
