package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"net/http"
)

type ClusterManager struct {
	Zone               string
	ProcessCounterDesc *prometheus.Desc
}

// Describe 指标描述
func (c *ClusterManager) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ProcessCounterDesc
}

// Collect 指标采集
func (c *ClusterManager) Collect(ch chan<- prometheus.Metric) {
		SystemProcessBycount := c.SystemState()
	for host,processCount := range SystemProcessBycount {
		ch<- prometheus.MustNewConstMetric(c.ProcessCounterDesc,
			prometheus.CounterValue,
			float64(processCount),
			host)
	}
}

// SystemState 采集方法
func (c *ClusterManager) SystemState() (processCountByHost map[string]int){
	processCountByHost = map[string]int{
		"192" : 111,
	}
	return
}

// NewClusterManager 创建集群采集管理结构体
func NewClusterManager(zone string) *ClusterManager {
	return &ClusterManager{
		Zone: "testServer",
		ProcessCounterDesc: prometheus.NewDesc(
			"clustermanager_process_total",
			"Number of restart process",
			[]string{"host"},
			prometheus.Labels{"zone": "zone"}),
	}
}

func main() {
	workeA := NewClusterManager("demo")
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(workeA)
	gatherers := prometheus.Gatherers{prometheus.DefaultGatherer,reg}
	h := promhttp.HandlerFor(gatherers,
		promhttp.HandlerOpts{
			ErrorLog: log.NewErrorLogger(),
			ErrorHandling: promhttp.ContinueOnError,
		})
	http.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		h.ServeHTTP(writer,request)
		http.Handle("/metrics",promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8080",nil))
	})
	select {

	}
}
