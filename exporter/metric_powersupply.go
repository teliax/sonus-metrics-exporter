package exporter

import (
	"encoding/xml"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const powerSupplyUrlSuffix = "/operational/system/powerSupplyStatus/"

func GetPowerSupplyUrl(ctx MetricContext) string {
	return ctx.APIBase + powerSupplyUrlSuffix
}

var PowerSupplyMetrics = map[string]*prometheus.Desc{
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

func ProcessPowerSupplies(ctx MetricContext, xmlBody *[]byte, ch chan<- prometheus.Metric, result chan<- bool) {
	powerSupplies := new(powerSupplyCollection)
	err := xml.Unmarshal(*xmlBody, &powerSupplies)
	if err != nil {
		log.Errorf("Failed to deserialize powerSupplyStatus XML: %v", err)
		result <- false
		return
	}

	for _, psu := range powerSupplies.PowerSupplyStatus {
		ch <- prometheus.MustNewConstMetric(PowerSupplyMetrics["PowerSupply_Power_Fault"], prometheus.GaugeValue, psu.powerFaultToMetric(), psu.ServerName, psu.PowerSupplyID)
		ch <- prometheus.MustNewConstMetric(PowerSupplyMetrics["PowerSupply_Voltage_Fault"], prometheus.GaugeValue, psu.voltageFaultToMetric(), psu.ServerName, psu.PowerSupplyID)
	}
	log.Info("Power Supply Metrics collected")
	result <- true
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
