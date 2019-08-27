package main

import (
	"fmt"
	"flag"
	"net/http"

	"github.com/keepalived-exporter/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	add = flag.String("listen-address", ":9999", "The address to listen on for HTTP requests.")

	keepalived_status = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "keepalived_status",
		Help: "The gauge is represents keepalived process status",
	})

	keepalived_vip = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "keepalived_vip",
		Help: "The gauge is represents keepalive vip network status",
	})
)

func init() {
	// prometheus.MustRegister(keepalived_status)
	prometheus.MustRegister(pkg.Keepalived_vip)
}

func main() {

	fmt.Println("start a exporter for keepalived ...")
	flag.Parse()
	//updateKeepalivedMetrics()
	pkg.UpdateKeepalivedVIP()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*add, nil)
}