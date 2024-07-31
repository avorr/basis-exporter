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
	//FOO                string
	appVersion  = "dev"
	Api         *client.Api
	port        string
	region      string
	req4xxCount float64
	req5xxCount float64
)

type Collector struct {
	fooMetric  *prometheus.Desc
	barMetric  *prometheus.Desc
	cpuMetric  *prometheus.GaugeVec
	ramMetric  *prometheus.GaugeVec
	diskMetric *prometheus.GaugeVec
	reqCount   *prometheus.GaugeVec
	duration   *prometheus.HistogramVec
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
	Pool string `json:"pool"`
	//PresentTo   []int  `json:"presentTo"`
	//PurgeTime   int    `json:"purgeTime"`
	ResID string `json:"resId"`
	//ResName string `json:"resName"`
	//Role        string `json:"role"`
	SepID   int    `json:"sepId"`
	SepType string `json:"sepType"`
	//Shareable   bool   `json:"shareable"`
	SizeMax  int     `json:"sizeMax"`
	SizeUsed float64 `json:"sizeUsed"`
	//Snapshots  []any  `json:"snapshots"`
	//Status     string `json:"status"`
	//TechStatus string `json:"techStatus"`
	Type string `json:"type"`
	//Vmid int    `json:"vmid"`
}

func NewFooCollector(reg *prometheus.Registry) *Collector {
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
		"resId",
		"pool",
		"sepId",
		"sepType",
		"sizeUsed",
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

	f := &Collector{
		cpuMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "portal",
			Name:      "server_cpu",
			Help:      "Compute cpu",
		}, labels),
		ramMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "portal",
			Name:      "server_ram",
			Help:      "Compute ram",
		}, labels),
		diskMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "portal",
			Name:      "server_disk",
			Help:      "Compute ram",
		}, diskLabels),
		reqCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "portal",
			Name:      "api_requests_total",
			Help:      "Status code 2xx and 5xx counter",
		}, counterLabels),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "portal",
			Name:      "api_request",
			Help:      "Duration of the request.",
			// Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5),
			// Buckets: prometheus.LinearBuckets(0.1, 5, 5),
			//Buckets: []float64{0.1, 0.15, 0.2, 0.25, 0.3},
			Buckets: []float64{30, 50, 70, 90, 110, 130, 150, 170, 190},
		}, []string{"api_endpoint", "status", "method"}),
	}

	reg.MustRegister(f.cpuMetric, f.ramMetric, f.diskMetric, f.reqCount, f.duration)
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
	computes, disks, err := collector.GetComputes()
	if err != nil {
		log.Fatal(err)
	}
	if computes == nil || disks == nil {
		return
	}

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
			diskPool := volumeMap[disk.ID].Pool
			if diskPool == "" {
				diskPool = NA
			}
			resId := volumeMap[disk.ID].ResID
			if resId == "" {
				resId = NA
			}
			sepId := volumeMap[disk.ID].SepID
			sepType := volumeMap[disk.ID].SepType
			if sepType == "" {
				sepType = NA
			}
			sizeUsed := volumeMap[disk.ID].SizeUsed

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
				"rgName":   compute.RgName,
				"resId":    resId,
				"pool":     diskPool,
				"sepId":    strconv.Itoa(sepId),
				"sepType":  sepType,
				"sizeUsed": strconv.FormatFloat(sizeUsed, 'f', 6, 64),
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
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}), collectors.NewGoCollector())
	collector := NewFooCollector(reg)
	if err := reg.Register(collector); err != nil {
		panic(err)
	}

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func (collector *Collector) GetComputeInfo() ([]byte, []byte) {
	Api = client.NewApi()
	if err := Api.Auth(); err != nil {
		log.Fatalf(err.Error())
	}

	out1 := make(chan []byte)
	out2 := make(chan []byte)
	done := make(chan bool)

	//var wg sync.WaitGroup
	//var data = strings.NewReader("page=1&size=1")

	go func() {
		//computeReqStart := time.Now()
		now := time.Now()
		resp, statusCode, err := Api.NewRequestPost("compute/list", strings.NewReader(""))
		collector.duration.With(prometheus.Labels{
			"api_endpoint": fmt.Sprintf("%s/%s", client.ApiUrl, "compute/list"),
			"method":       "GET",
			"status":       strconv.Itoa(statusCode),
		}).Observe(time.Since(now).Seconds())

		//resp, err := ReadFromFile("data.txt")
		if statusCode >= 400 && statusCode <= 499 {
			req4xxCount++
			done <- false
			return
		}
		if statusCode >= 500 && statusCode <= 526 {
			req5xxCount++
			done <- false
			return
		}
		if err != nil {
			done <- false
			log.Println(err)
		}
		//WriteToFile(resp, "computes.txt")
		done <- true
		out1 <- resp
	}()

	go func() {
		now := time.Now()
		diskResp, statusCode, err := Api.NewRequestPost("disks/list", strings.NewReader(""))

		collector.duration.With(prometheus.Labels{
			"api_endpoint": fmt.Sprintf("%s/%s", client.ApiUrl, "disks/list"),
			"method":       "GET",
			"status":       strconv.Itoa(statusCode),
		}).Observe(time.Since(now).Seconds())

		//diskResp, err := ReadFromFile("disks.txt")
		if statusCode >= 400 && statusCode <= 499 {
			req4xxCount++
			done <- false
			return
		}
		if statusCode >= 500 && statusCode <= 526 {
			req5xxCount++
			done <- false
			return
		}
		if err != nil {

			done <- false
			log.Println(err)
		}
		//WriteToFile(diskResp, "disks.txt")
		done <- true
		out2 <- diskResp
	}()

	var requestsResult bool
	for j := 0; j <= 1; j++ {
		requestsResult = <-done
		if !requestsResult {
			break
		}
	}

	if !requestsResult {
		return nil, nil
	}

	return <-out1, <-out2
}

func (collector *Collector) GetComputes() ([]Data, []DiskData, error) {
	respBytes, disksResp := collector.GetComputeInfo()
	if respBytes == nil || disksResp == nil {
		return nil, nil, nil
	}

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
