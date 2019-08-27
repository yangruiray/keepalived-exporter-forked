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
)

func init() {
	prometheus.MustRegister(pkg.Keepalived_vip)
}

func main() {
	fmt.Println("start a exporter for keepalived ...")
	flag.Parse()

	pkg.UpdateKeepalivedVIP()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*add, nil)
}