package metrics

import (
	"encoding/xml"

	"sonus-metrics-exporter/lib"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	powerSupplyName      = "PowerSupply"
	powerSupplyUrlSuffix = "/operational/system/powerSupplyStatus/"
)

var PowerSupplyMetric = lib.SonusMetric{
	Name:       powerSupplyName,
	Processor:  processPowerSupplies,
	URLGetter:  getPowerSupplyUrl,
	APIMetrics: powerSupplyMetrics,
	Repetition: lib.RepeatNone,
}

func getPowerSupplyUrl(ctx lib.MetricContext) string {
	return ctx.APIBase + powerSupplyUrlSuffix
}

var powerSupplyMetrics = map[string]*prometheus.Desc{
	"PowerSupply_Power_Fault": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "powersupply", "powerfault"),
		"Is there a power fault, per supply",
		[]string{"server", "powerSupplyID"}, nil,
	),
	"PowerSupply_Voltage_Fault": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "powersupply", "voltagefault"),
		"Is there a voltage fault, per supply",
		[]string{"server", "powerSupplyID"}, nil,
	),
}

func processPowerSupplies(ctx lib.MetricContext, xmlBody *[]byte) {
	var (
		errors        []*error
		powerSupplies = new(powerSupplyCollection)
	)

	err := xml.Unmarshal(*xmlBody, &powerSupplies)

	if err != nil {
		log.Errorf("Failed to deserialize powerSupplyStatus XML: %v", err)
		errors = append(errors, &err)
		ctx.ResultChannel <- lib.MetricResult{Name: powerSupplyName, Success: false, Errors: errors}
		return
	}

	for _, psu := range powerSupplies.PowerSupplyStatus {
		ctx.MetricChannel <- prometheus.MustNewConstMetric(powerSupplyMetrics["PowerSupply_Power_Fault"], prometheus.GaugeValue, psu.powerFaultToMetric(), psu.ServerName, psu.PowerSupplyID)
		ctx.MetricChannel <- prometheus.MustNewConstMetric(powerSupplyMetrics["PowerSupply_Voltage_Fault"], prometheus.GaugeValue, psu.voltageFaultToMetric(), psu.ServerName, psu.PowerSupplyID)
	}

	log.Info("Power Supply Metrics collected")
	ctx.ResultChannel <- lib.MetricResult{Name: powerSupplyName, Success: true}
}

/*
<collection xmlns:y="http://tail-f.com/ns/rest">
  <powerSupplyStatus xmlns="http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0">
    <serverName>densbc01a</serverName>
    <powerSupplyId>PSA</powerSupplyId>
    <present>true</present>
    <productName>TECTROL  TC92S-1525R</productName>
    <serialNum>00000000</serialNum>
    <partNum>TC92S-1525R</partNum>
    <powerFault>false</powerFault>
    <voltageFault>false</voltageFault>
  </powerSupplyStatus>
...
</collection>
*/

type powerSupplyCollection struct {
	PowerSupplyStatus []*powerSupplyStatus `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 powerSupplyStatus,omitempty"`
}

type powerSupplyStatus struct {
	ServerName    string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 serverName"`
	PowerSupplyID string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 powerSupplyId"`
	Present       bool   `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 present"`
	PowerFault    bool   `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 powerFault"`
	VoltageFault  bool   `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 voltageFault"`
}

func (p powerSupplyStatus) powerFaultToMetric() float64 {
	if p.PowerFault {
		return 1
	} else {
		return 0
	}
}

func (p powerSupplyStatus) voltageFaultToMetric() float64 {
	if p.VoltageFault {
		return 1
	} else {
		return 0
	}
}
