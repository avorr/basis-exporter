package main

import (
	"basis-exporter/client"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	NA = "NA"
)

var (
	Api                *client.Api
	computeReqDuration time.Duration
	deskReqDuration    time.Duration
	totalReqDuration   time.Duration
	port               string
	region             string
	req2xxCount        float64
	req5xxCount        float64
)

type Collector struct {
	fooMetric               *prometheus.Desc
	barMetric               *prometheus.Desc
	cpuMetric               *prometheus.GaugeVec
	ramMetric               *prometheus.GaugeVec
	diskMetric              *prometheus.GaugeVec
	reqCount                *prometheus.GaugeVec
	portalApiLatency        *prometheus.GaugeVec
	portalApiLatencyCompute *prometheus.GaugeVec
	portalApiLatencyDisks   *prometheus.GaugeVec
	computeMetrics          ComputeMetrics
}

type Disks struct {
	ID int `json:"id"`
	//PciSlot int `json:"pciSlot"`
}

type Interfaces struct {
	//ConnID      int    `json:"connId"`
	//ConnType    string `json:"connType"`
	//DefGw       string `json:"defGw"`
	//Enabled     bool   `json:"enabled"`
	//FlipgroupID int    `json:"flipgroupId"`
	//GUID        string `json:"guid"`
	IPAddress string `json:"ipAddress"`
	//ListenSSH   bool   `json:"listenSsh"`
	//Mac         string `json:"mac"`
	Name  string `json:"name"`
	NetID int    `json:"netId"`
	//NetMask     int    `json:"netMask"`
	NetType string `json:"netType"`
	//PciSlot     int    `json:"pciSlot"`
	//Target      string `json:"target"`
	//Type        string `json:"type"`
	//Vnfs        []any  `json:"vnfs"`
}

type Data struct {
	AccountID   int    `json:"accountId"`
	AccountName string `json:"accountName"`
	//ACL               []any        `json:"acl"`
	//AffinityLabel     string       `json:"affinityLabel"`
	//AffinityRules     []any        `json:"affinityRules"`
	//AffinityWeight    int          `json:"affinityWeight"`
	//AntiAffinityRules []any        `json:"antiAffinityRules"`
	//Arch              string       `json:"arch"`
	//BootOrder         []string     `json:"bootOrder"`
	//BootdiskSize      int          `json:"bootdiskSize"`
	//CloneReference    int          `json:"cloneReference"`
	//Clones            []any        `json:"clones"`
	//ComputeciID       int          `json:"computeciId"`
	Cpus int `json:"cpus"`
	//CreatedBy         string       `json:"createdBy"`
	//CreatedTime       int          `json:"createdTime"`
	//CustomFields      struct{}     `json:"customFields"`
	//DeletedBy         string       `json:"deletedBy"`
	//DeletedTime       int          `json:"deletedTime"`
	//Desc              string       `json:"desc"`
	//Devices           struct{}     `json:"devices"`
	Disks []Disks `json:"disks"`
	//Driver            string       `json:"driver"`
	//Gid               int          `json:"gid"`
	//GUID              int          `json:"guid"`
	ID int `json:"id"`
	//ImageID           int          `json:"imageId"`
	Interfaces []Interfaces `json:"interfaces"`
	//LockStatus        string       `json:"lockStatus"`
	//ManagerID         int          `json:"managerId"`
	//ManagerType       string       `json:"managerType"`
	//Migrationjob      int          `json:"migrationjob"`
	//Milestones        int          `json:"milestones"`
	Name string `json:"name"`
	//NeedReboot        bool         `json:"needReboot"`
	//Pinned            bool         `json:"pinned"`
	RAM int `json:"ram"`
	//ReferenceID       string       `json:"referenceId"`
	//Registered        bool         `json:"registered"`
	//ResName           string       `json:"resName"`
	RgID   int    `json:"rgId"`
	RgName string `json:"rgName"`
	//SnapSets          []any        `json:"snapSets"`
	//StatelessSepID    int          `json:"statelessSepId"`
	//StatelessSepType  string       `json:"statelessSepType"`
	//Status            string       `json:"status"`
	//Tags              Tags         `json:"tags"`
	//Tags           interface{} `json:"tags"`
	Tags map[string]string `json:"tags"`
	//TechStatus     string            `json:"techStatus"`
	//TotalDisksSize int               `json:"totalDisksSize"`
	//UpdatedBy      string            `json:"updatedBy"`
	//UpdatedTime    int               `json:"updatedTime"`
	//UserManaged    bool              `json:"userManaged"`
	//Vgpus          []any             `json:"vgpus"`
	//VinsConnected  int               `json:"vinsConnected"`
	//VirtualImageID int               `json:"virtualImageId"`
}

type Computes struct {
	Data       []Data `json:"data"`
	EntryCount int    `json:"entryCount"`
}

type DisksInfo struct {
	Data       []DiskData `json:"data"`
	EntryCount int        `json:"entryCount"`
}

type DiskData struct {
	//AccountID   int    `json:"accountId"`
	//AccountName string `json:"accountName"`
	//ACL         struct{} `json:"acl"`
	//Computes struct {
	//	Num4847 string `json:"4847"`
	//} `json:"computes"`
	//CreatedTime     int    `json:"createdTime"`
	//DeletedTime     int    `json:"deletedTime"`
	//Desc string `json:"desc"`
	//DestructionTime int    `json:"destructionTime"`
	//Devicename      string `json:"devicename"`
	//Gid             int    `json:"gid"`
	ID int `json:"id"`
	//ImageID         int    `json:"imageId"`
	//Images          []any  `json:"images"`
	//Iotune          struct {
	//	TotalIopsSec int `json:"total_iops_sec"`
	//} `json:"iotune"`
	//MachineID   any    `json:"machineId"`
	//MachineName any    `json:"machineName"`
	Name string `json:"name"`
	//Order       int    `json:"order"`
	//Params      string `json:"params"`
	//ParentID    int    `json:"parentId"`
	//PciSlot     int    `json:"pciSlot"`
	//Pool string `json:"pool"`+
	//PresentTo   []int  `json:"presentTo"`
	//PurgeTime   int    `json:"purgeTime"`
	//ResID   string `json:"resId"`
	//ResName string `json:"resName"`
	//Role        string `json:"role"`
	//SepID   int    `json:"sepId"`+
	//SepType string `json:"sepType"`+
	//Shareable   bool   `json:"shareable"`
	SizeMax int `json:"sizeMax"`
	//SizeUsed float64 `json:"sizeUsed"`+
	//Snapshots  []any  `json:"snapshots"`
	//Status     string `json:"status"`
	//TechStatus string `json:"techStatus"`
	Type string `json:"type"`
	//Vmid int    `json:"vmid"`
}

type ComputeMetrics interface {
	GetComputes() ([]Data, []DiskData, error)
}

func NewFooCollector(metrics ComputeMetrics, reg *prometheus.Registry) *Collector {
	labels := []string{
		"name",
		"location",
		"owner",
		"project",
		"phase",
		"contract_service",
		"accountId",
		"accountName",
		"rgId",
		"rgName",
		"net_id",
		"net_name",
		"id",
		"ip",
	}
	diskLabels := []string{
		"name",
		"location",
		"owner",
		"project",
		"phase",
		"contract_service",
		"accountId",
		"accountName",
		"rgId",
		"rgName",
		"id",
		"ip",
		"disk_type",
		"disk_name",
	}
	counterLabels := []string{
		"code_status",
		"hostname",
		"instance",
		"job",
		"monitor",
	}

	requestLabels := []string{
		"api_endpoint",
		"hostname",
		"instance",
		"job",
		"le",
		"monitor",
	}

	f := &Collector{
		cpuMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "portal_server_cpu",
			Help: "Compute cpu",
		}, labels),
		////portal_server_cpu {accountId="1", accountName="name", resource_group_id="2", resource_name="rg_name", net_id="3", net_name="extnet", id="12", compute="test_vm", ip="172.22.0.12"} 4
		//			portal_server_ram {accountId="1", accountName="name", resource_group_id="2", resource_name="rg_name", net_id="3", net_name="extnet", id="12", compute="test_vm", ip="172.22.0.12"} 8
		//fooMetric: prometheus.NewDesc("foo_metric", "Shows whether a foo has occurred in our cluster", nil, labels),
		ramMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "portal_server_ram",
			Help: "Compute ram",
		}, labels),
		diskMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "portal_server_disk",
			Help: "Compute ram",
		}, diskLabels),
		reqCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "portal_api_requests_total",
			Help: "Status code 200 counter",
		}, counterLabels),
		portalApiLatency: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "portal_api_latency_bucket",
			Help: "Api response time",
		}, requestLabels),
		portalApiLatencyCompute: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "portal_api_latency_compute",
			Help: "Api response compute time",
		}, requestLabels),
		portalApiLatencyDisks: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "portal_api_latency_disks",
			Help: "Api response disks time",
		}, requestLabels),
		computeMetrics: metrics,
	}

	//var requestsCounter uint64 = 0
	// register counter in Prometheus collector
	//prometheus.MustRegister(prometheus.NewCounterFunc(
	//func() float64 {
	//	return float64(atomic.LoadUint64(&requestsCounter))
	//}))
	// somewhere in your code
	//atomic.AddUint64(&requestsCounter, 1)
	//var requestsCounter uint64 = 0

	// register counter in Prometheus collector
	//prometheus.MustRegister(prometheus.NewCounterFunc(
	//prometheus.CounterOpts{
	//	Name: "requests_total",
	//	Help: "Counts number of requests",
	//},

	//),
	//)

	// somewhere in your code
	//atomic.LoadUint64(&requestsCounter)
	//atomic.AddUint64(&requestsCounter, 1)

	reg.MustRegister(f.cpuMetric)
	reg.MustRegister(f.ramMetric)
	reg.MustRegister(f.diskMetric)
	reg.MustRegister(f.reqCount)
	reg.MustRegister(f.portalApiLatency)
	reg.MustRegister(f.portalApiLatencyCompute)
	reg.MustRegister(f.portalApiLatencyDisks)
	return f
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	//ch <- collector.fooMetric
	//ch <- collector.barMetric
}

func getInterfaceInfo(interfaces []Interfaces) (map[string]string, error) {
	for _, netInt := range interfaces {
		if netInt.NetType == "EXTNET" {
			return map[string]string{
				"net_id":   strconv.Itoa(netInt.NetID),
				"net_name": netInt.Name,
				"ip":       netInt.IPAddress,
			}, nil
		}
	}

	return map[string]string{}, fmt.Errorf("NetType EXTNET not fount")
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	totalReqStart := time.Now()
	computes, disks, err := collector.computeMetrics.GetComputes()
	if err != nil {
		log.Fatal(err)
	}
	totalReqDuration = time.Since(totalReqStart)

	volumeMap := make(map[int]DiskData)

	for _, disk := range disks {
		volumeMap[disk.ID] = disk
	}

	for _, compute := range computes {
		serviceInfo := strings.Split(compute.Name, ".")
		location, owner, project, phase := NA, NA, NA, NA

		if len(serviceInfo) >= 6 {
			location = serviceInfo[3]
			owner = serviceInfo[4]
			project = serviceInfo[1]
			phase = serviceInfo[2]
		}

		netInfo, err := getInterfaceInfo(compute.Interfaces)
		contractService, ok := compute.Tags["contract_service"]
		if !ok {
			contractService = NA
		}
		if err != nil {
			netInfo["net_id"] = NA
			netInfo["net_name"] = NA
			netInfo["ip"] = NA
		}

		labels := map[string]string{
			"name":             compute.Name,
			"location":         location,
			"owner":            owner,
			"project":          project,
			"phase":            phase,
			"contract_service": contractService,
			"accountId":        strconv.Itoa(compute.AccountID),
			"accountName":      compute.AccountName,
			"rgId":             strconv.Itoa(compute.RgID),
			"rgName":           compute.RgName,
			"net_id":           netInfo["net_id"],
			"net_name":         netInfo["net_name"],
			"id":               strconv.Itoa(compute.ID),
			"ip":               netInfo["ip"],
		}

		collector.cpuMetric.With(labels).Set(float64(compute.Cpus))
		collector.ramMetric.With(labels).Set(float64(compute.RAM))

		for _, disk := range compute.Disks {
			diskType := volumeMap[disk.ID].Type
			if diskType == "" {
				diskType = NA
			}

			diskName := volumeMap[disk.ID].Name
			if diskName == "" {
				diskName = NA
			}

			diskLabels := map[string]string{
				"name":             compute.Name,
				"location":         location,
				"owner":            owner,
				"project":          project,
				"phase":            phase,
				"contract_service": contractService,
				"accountId":        strconv.Itoa(compute.AccountID),
				"accountName":      compute.AccountName,
				//"resource_group_id": strconv.Itoa(compute.RgID),
				"rgId": strconv.Itoa(compute.RgID),
				//"resource_name":     compute.ResName,
				"rgName": compute.RgName,
				//"net_id": netInfo["net_id"],
				//"net_name":   netInfo["net_name"],
				"id": strconv.Itoa(compute.ID),
				//"compute":    compute.Name,
				//"ip":        compute.Interfaces[0].IPAddress,
				"ip":        netInfo["ip"],
				"disk_type": diskType,
				"disk_name": diskName,
				//"disk_type": volumeMap[disk.ID].Type,
			}
			collector.diskMetric.With(diskLabels).Set(float64(volumeMap[disk.ID].SizeMax))
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = NA
		fmt.Println(err)
	}

	reqLabels := map[string]string{
		"code_status": "200",
		"hostname":    hostname,
		"instance":    fmt.Sprintf("%s:%s", hostname, port),
		"job":         "portal_exporter",
		"monitor":     strings.ToLower(region),
	}

	req2xxCount += 1
	collector.reqCount.With(reqLabels).Set(req2xxCount)

	reqLatency := map[string]string{
		"api_endpoint": client.ApiUrl,
		"hostname":     hostname,
		"instance":     fmt.Sprintf("%s:%s", hostname, port),
		"job":          "portal_exporter",
		"le":           "+inf",
		"monitor":      strings.ToLower(region),
	}

	collector.portalApiLatency.With(reqLatency).Set(totalReqDuration.Minutes())

	reqLatency["api_endpoint"] = fmt.Sprintf("%s/%s", client.ApiUrl, "compute/list")
	collector.portalApiLatencyCompute.With(reqLatency).Set(computeReqDuration.Minutes())

	reqLatency["api_endpoint"] = fmt.Sprintf("%s/%s", client.ApiUrl, "disks/list")
	collector.portalApiLatencyDisks.With(reqLatency).Set(deskReqDuration.Minutes())

	//computes := GetComputeInfo()
	//var result Computes
	//err := json.Unmarshal(computes, &result)
	//if err != nil {
	//log.Fatal(err)
	//}

	//for _, i2 := range result.Data {
	//	labels := map[string]string{"cluster_name": i2.Name, "vdcname": i2.Desc, "value": i2.RgName}
	//	collector.gmetric.With(labels).Set(metricValue)
	//}

	//m1 := prometheus.MustNewConstMetric(collector.fooMetric, prometheus.GaugeValue, metricValue)
	//m2 := prometheus.MustNewConstMetric(collector.barMetric, prometheus.GaugeValue, metricValue)
	//m1 = prometheus.NewMetricWithTimestamp(time.Now().Add(-time.Hour), m1)
	//m2 = prometheus.NewMetricWithTimestamp(time.Now(), m2)
	//ch <- m1
	//ch <- m2
}

func init() {
	ok := false
	port, ok = os.LookupEnv("BE_PORT")
	if !ok || port == "" {
		port = "2112"
	}
	region, ok = os.LookupEnv("BE_REGION")
	if !ok || region == "" {
		region = NA
	}
}

func main() {
	//var foo metrics
	//foo.GetComputes()
	//return
	//var m metrics
	//computes, _, err := m.GetComputes()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//volumeMap := make(map[int]DiskData)
	//for _, disk := range disks {
	//	if disk.SizeMax != 50 {
	//		volumeMap[disk.ID] = disk
	//	}
	//}

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}), collectors.NewGoCollector())
	collector := NewFooCollector(&metrics{}, reg)
	if err := reg.Register(collector); err != nil {
		panic(err)
	}

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func GetComputeInfo() ([]byte, []byte) {
	Api = client.NewApi()
	if err := Api.Auth(); err != nil {
		log.Fatalf(err.Error())
	}

	out1 := make(chan []byte)
	out2 := make(chan []byte)

	//var wg sync.WaitGroup
	//var data = strings.NewReader("page=1&size=1")

	go func() {
		computeReqStart := time.Now()
		resp, err := Api.NewRequestPost("compute/list", strings.NewReader(""))
		//resp, err := ReadFromFile("data.txt")
		if err != nil {
			log.Println(err)
		}
		computeReqDuration = time.Since(computeReqStart)
		//WriteToFile(resp, "computes.txt")
		out1 <- resp
	}()

	go func() {
		diskReqStart := time.Now()
		diskResp, err := Api.NewRequestPost("disks/list", strings.NewReader(""))
		//diskResp, err := ReadFromFile("disks.txt")
		if err != nil {
			log.Println(err)
		}
		deskReqDuration = time.Since(diskReqStart)
		//WriteToFile(diskResp, "disks.txt")
		out2 <- diskResp
	}()

	return <-out1, <-out2
	//return resp, diskResp
}

type metrics struct{}

func (d *metrics) GetComputes() ([]Data, []DiskData, error) {
	respBytes, disksResp := GetComputeInfo()

	var computes Computes
	if err := json.Unmarshal(respBytes, &computes); err != nil {
		return nil, nil, err
	}

	var disks DisksInfo
	if err := json.Unmarshal(disksResp, &disks); err != nil {
		return nil, nil, err
	}

	return computes.Data, disks.Data, nil
}

func WriteToFile(data []byte, file string) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}(f)

	_, err = f.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadFromFile(file string) ([]byte, error) {
	return os.ReadFile(file)
}
