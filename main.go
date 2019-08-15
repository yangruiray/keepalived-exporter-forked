package main

import (
	"fmt"
	"log"
	"flag"
	"os/exec"
	"net/http"


	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
)

const (
	defaultKeepalivedConf = "/etc/keepalived/keepalived.conf"
	defaultUpdateInterval = 2
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
		Name: "keepalived_vip",
		Help: "The gauge is represents keepalive vip network status",
	})
)

func inspectKeepalivedStatus() string {
	cmd := "systemctl is-active keepalived.service | grep -w active &2>/dev/null"
	std, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatalf("fail to exec command to get keepalived status, err: %v", err)
		keepalived_status.Set(0.0)
	}

	return string(std)
}

func updateKeepalivedMetrics() {
	go func() {
		for {
			if inspectKeepalivedStatus() == "active" {
				keepalived_status.Set(1.0)
			} else {
				keepalived_status.Set(0.0)
			}
			time.Sleep(defaultUpdateInterval * time.Second)
		}
	}()
}

func init() {
	prometheus.MustRegister(keepalived_status)
	prometheus.MustRegister(keepalived_vip)
}

func main() {
	fmt.Println("start a exporter for keepalived ...")
	flag.Parse()
	updateKeepalivedMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*add, nil)
}
