package pkg

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultKeepalivedConf = "/etc/keepalived/keepalived.conf"
)

// Exec to get state of keepalived
func (m *KeepalivedMetrics) GetState() (keepalivedVip map[string]int) {
	gaugeValue := updateKeepalivedVIP()
	hostName := accquireHostname()

	keepalivedVip = map[string]int{
		fmt.Sprintf("%v", hostName): gaugeValue,
	}

	return
}

// Collect func is for write gauge value to channel
func (m *KeepalivedMetrics) Collect(ch chan<- prometheus.Metric) {
	keepalived_vip := m.GetState()

	for host, values := range keepalived_vip {
		ch <- prometheus.MustNewConstMetric(
			m.Keepalived_vip,
			prometheus.GaugeValue,
			float64(values),
			host,
		)
	}
}

// Write describe to channel
func (m *KeepalivedMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.Keepalived_vip
}

func NewKeepalivedMetrics() *KeepalivedMetrics {
	return &KeepalivedMetrics{
		Keepalived_vip: prometheus.NewDesc(
			"keepalived_vip_ready",
			"vip on current node",
			[]string{"host"},
			nil,
		),
	}
}

// accquire current hostname
func accquireHostname() string {
	std, err := exec.Command("/bin/bash", "-c", "hostname", "-f").CombinedOutput()
	if err != nil {
		glog.Errorf("failed to get current node hostname")
		return ""
	}

	hostname := strings.TrimSuffix(string(std), "\n")
	return hostname
}

// read keepalived config and return context
func readKeepalivedFile() string {
	file, err := ioutil.ReadFile(defaultKeepalivedConf)
	if err != nil {
		glog.Errorf("fail to open keepalived config on path %v, err: %v", defaultKeepalivedConf, err)
		return ""
	}

	return string(file)
}

// parse keepalived config and return ip collection of keepalived VIP
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

// inspect current node whether vip is on current machine
func inspectKeepalivedVIP() bool {
	var vipCheckArray []bool
	keepalivedContent := readKeepalivedFile()
	currentHost := accquireHostname()

	if keepalivedContent == "" {
		glog.Fatalf("keepalived config parse failed, return with empty string")
		return false
	}

	IPCollection := parseKeepalivedVIP(keepalivedContent)

	for _, v := range IPCollection {
		cmd := fmt.Sprintf("ip addr | grep %v 2>/dev/null", v)
		std, err := exec.Command("bash", "-c", cmd).Output()

		if err != nil {
			glog.Errorf("exec command: [%v], error occurs: %v", cmd, err)
			return false
		}

		currentVIP := parseKeepalivedVIP(string(std))
		if len(currentVIP) == 0 {
			glog.Errorf("current host: \"%v\" lost vip", currentHost)
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

// Exec to inspect current node vip state
func updateKeepalivedVIP() int {
	if inspectKeepalivedVIP() == true {
		return 1
	}
	return 0
}
