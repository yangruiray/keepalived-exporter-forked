package main

import (
	"fmt"
	"flag"
	"regexp"
	"os/exec"
	"net/http"
	"io/ioutil"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultKeepalivedConf = "/etc/keepalived/keepalived.conf"
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

func readKeepalivedFile() string {
	file, err := ioutil.ReadFile(defaultKeepalivedConf)
	if err != nil {
		glog.Errorf("fail to open keepalived config on path %v, err: %v", defaultKeepalivedConf, err)
		return ""
	}

	return string(file)
}

func parseKeepalivedVIP(inputContent string) []string {
	var ipCollection []string
	vipMap := make(map[string]bool)

	vipPattern := `(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])`
	vipObject := regexp.MustCompile(vipPattern)
	matchVIP := vipObject.FindAllStringSubmatch(inputContent, -1)

	for _, value := range matchVIP {
		// TODO should we ignore localhost?
		if value[0] != "127.0.0.1" {
			if _, ok := vipMap[value[0]]; !ok {
				vipMap[value[0]] = true
				ipCollection = append(ipCollection, value[0])
			}
		}
	}

	if len(ipCollection) == 0 {
		// glog.Errorf("keepalived ip address is empty")
		return nil
	}

	return ipCollection
}

func inspectKeepalivedVIP() bool {
	var vipCheckArray []bool
	keepalivedContent := readKeepalivedFile()

	if keepalivedContent == "" {
		glog.Fatalf("keepalived config parse failed, return with empty string")
		return false
	}

	IPCollection := parseKeepalivedVIP(keepalivedContent)

	for _, v := range IPCollection {
		cmd := fmt.Sprintf("ip addr | grep %v &2>/dev/null", v)
		std, err := exec.Command("bash", "-c", cmd).CombinedOutput()

		if err != nil {
			glog.Errorf("exec command: [%v], error occurs: %v", cmd, err)
			keepalived_status.Set(0.0)
			return false
		}

		currentVIP := parseKeepalivedVIP(string(std))
		if len(currentVIP) == 0 {
			// glog.Errorf("current host lost vip")
			return false
		}

		if parseKeepalivedVIP(string(std))[0] == v {
			vipCheckArray = append(vipCheckArray, true)
		} else {
			vipCheckArray = append(vipCheckArray, false)
		}
	}

	for _, v := range vipCheckArray {
		if v == false {
			glog.Infof("current host lost vip")
			return false
		}
	}

	if len(vipCheckArray) == len(IPCollection) {
		return true
	}

	//glog.Infof("current host vip counts mismatch keepalived conf's counts")
	return false
}

func inspectKeepalivedStatus() bool {
	cmd := "systemctl is-active keepalived.service | grep -w active &2>/dev/null"
	std, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		glog.Errorf("fail to exec command to get keepalived status, err: %v", err)
		keepalived_status.Set(0.0)
		return false
	}

	if string(std) != "" {
		return true
	}

	return false
}

func updateKeepalivedMetrics() {
	if inspectKeepalivedStatus() == true {
		keepalived_status.Set(1.0)
	} else {
		keepalived_status.Set(0.0)
	}
	//time.Sleep(defaultUpdateInterval * time.Second)
}

func updateKeepalivedVIP() {
	if inspectKeepalivedVIP() == true {
		keepalived_vip.Set(1.0)
	} else {
		keepalived_vip.Set(0.0)
	}
	//time.Sleep(defaultUpdateInterval * time.Second)
}

func init() {
	prometheus.MustRegister(keepalived_status)
	prometheus.MustRegister(keepalived_vip)
}

func main() {

	fmt.Println("start a exporter for keepalived ...")
	flag.Parse()
	//updateKeepalivedMetrics()
	updateKeepalivedVIP()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*add, nil)
}
