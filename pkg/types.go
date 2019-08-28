package pkg

import "github.com/prometheus/client_golang/prometheus"

type KeepalivedMetrics struct {
	Keepalived_vip *prometheus.Desc
}