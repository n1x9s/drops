package metrics

import "github.com/prometheus/client_golang/prometheus"

var HTTPRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "secondbrain_http_requests_total",
		Help: "Total HTTP requests handled by the Second Brain API.",
	},
	[]string{"method", "path", "status"},
)

func Register() {
	prometheus.MustRegister(HTTPRequests)
}
