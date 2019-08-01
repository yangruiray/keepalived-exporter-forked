package main

import (
	"fmt"
	"net/http"
	"flag"
	"time"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	add = flag.String("listen-address", ":9999", "The address to listen on for HTTP requests.")
)

var (
	keepalived_status = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "keepalived_status",
		Help: "The gauge is represents keepalived process status",
	})
	keepalived_vip = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "keepalived_vip"
		Help: "The gauge is represents keepalive vip network status"
	})
)

func init() {
	prometheus.MustRegister(keepalived_status)
	prometheus.MustRegister(keepalived_vip)
}

func main() {
	fmt.Println("start a exporter for keepalived ...")

}
