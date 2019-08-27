package pkg

import "github.com/prometheus/client_golang/prometheus"

var (
	Keepalived_status = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "keepalived_status",
		Help: "The gauge is represents keepalived process status",
	})

	Keepalived_vip = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "keepalived_vip",
		Help: "The gauge is represents keepalive vip network status",
	})
)