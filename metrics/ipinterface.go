package metrics

import (
	"encoding/xml"
	"fmt"

	"sonus-metrics-exporter/lib"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	ipInterfaceName           = "IPInterface"
	ipInterfaceGroupURLFormat = "%s/operational/addressContext/%s/ipInterfaceGroup/%s/ipInterfaceStatus/"
)

var IPInterfaceMetric = lib.SonusMetric{
	Name:       ipInterfaceName,
	Processor:  processIPInterfaceStatus,
	URLGetter:  getIPInterfaceGroupUrl,
	APIMetrics: ipInterfaceMetrics,
	Repetition: lib.RepeatPerAddressContextIpInterfaceGroup,
}

func getIPInterfaceGroupUrl(ctx lib.MetricContext) string {
	return fmt.Sprintf(ipInterfaceGroupURLFormat, ctx.APIBase, ctx.AddressContext, ctx.IPInterfaceGroup)
}

var ipInterfaceMetrics = map[string]*prometheus.Desc{
	"IPInterface_Oper_Status": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "ipinterface", "status"),
		"Current status of ipInterfaceGroup",
		[]string{"name", "status_text"}, nil,
	),
	"IPInterface_Packets_Received": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "ipinterface", "rxpackets"),
		"Number of packets received on ipInterfaceGroup",
		[]string{"name"}, nil,
	),
	"IPInterface_Packets_Transmitted": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "ipinterface", "txpackets"),
		"Number of packets transmitted on ipInterfaceGroup",
		[]string{"name"}, nil,
	),
	"IPInterface_Bandwidth_Receive": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "ipinterface", "rxbandwidth"),
		"Receive bandwidth in use on interface, in bytes per second",
		[]string{"name"}, nil,
	),
	"IPInterface_Bandwidth_Transmit": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "ipinterface", "txbandwidth"),
		"Transmit bandwidth in use on interface, in bytes per second",
		[]string{"name"}, nil,
	),
	"IPInterface_Media_Streams": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "ipinterface", "media_streams"),
		"Number of media streams currently on ipInterfaceGroup",
		[]string{"name"}, nil,
	),
}

func processIPInterfaceStatus(ctx lib.MetricContext, xmlBody *[]byte) {
	var (
		errors       []*error
		ipInterfaces = new(ipInterfaceStatusCollection)
	)

	if len(*xmlBody) == 0 {
		ctx.ResultChannel <- lib.MetricResult{Name: ipInterfaceName, Success: true}
		return
	}

	err := xml.Unmarshal(*xmlBody, &ipInterfaces)

	if err != nil {
		log.Errorf("Failed to deserialize ipInterfaceStatus XML: %v", err)
		errors = append(errors, &err)
		ctx.ResultChannel <- lib.MetricResult{Name: ipInterfaceName, Success: false, Errors: errors}
		return
	}

	for _, ipInterfaceGroup := range ipInterfaces.IPInterfaces {
		ctx.MetricChannel <- prometheus.MustNewConstMetric(ipInterfaceMetrics["IPInterface_Oper_Status"], prometheus.GaugeValue, ipInterfaceGroup.operStateToMetric(), ipInterfaceGroup.Name, ipInterfaceGroup.OperState)

		ctx.MetricChannel <- prometheus.MustNewConstMetric(ipInterfaceMetrics["IPInterface_Packets_Received"], prometheus.CounterValue, ipInterfaceGroup.RxPackets, ipInterfaceGroup.Name)
		ctx.MetricChannel <- prometheus.MustNewConstMetric(ipInterfaceMetrics["IPInterface_Packets_Transmitted"], prometheus.CounterValue, ipInterfaceGroup.TxPackets, ipInterfaceGroup.Name)

		ctx.MetricChannel <- prometheus.MustNewConstMetric(ipInterfaceMetrics["IPInterface_Bandwidth_Receive"], prometheus.GaugeValue, ipInterfaceGroup.RxActualBandwidth, ipInterfaceGroup.Name)
		ctx.MetricChannel <- prometheus.MustNewConstMetric(ipInterfaceMetrics["IPInterface_Bandwidth_Transmit"], prometheus.GaugeValue, ipInterfaceGroup.TxActualBandwidth, ipInterfaceGroup.Name)

		ctx.MetricChannel <- prometheus.MustNewConstMetric(ipInterfaceMetrics["IPInterface_Media_Streams"], prometheus.GaugeValue, ipInterfaceGroup.NumMediaStreams, ipInterfaceGroup.Name)
	}

	log.Infof("IP Interface Metrics for Address Context %q, ipInterfaceGroup %q collected", ctx.AddressContext, ctx.IPInterfaceGroup)
	ctx.ResultChannel <- lib.MetricResult{Name: ipInterfaceName, Success: true}
}

/*
<collection xmlns:y="http://tail-f.com/ns/rest">
  <ipInterfaceStatus xmlns="http://sonusnet.com/ns/mibs/SONUS-GEN2-IP-INTERFACE/1.0">
    <name>CORE_1024</name>
    <ifindex>14</ifindex>
    <operState>resAllocated</operState>
    <oosReason>notApplicable</oosReason>
    <rxPackets>17654956</rxPackets>
    <txPackets>20793856</txPackets>
    <allocatedBandwidth>0</allocatedBandwidth>
    <bwDeviation>0</bwDeviation>
    <numMediaStreams>0</numMediaStreams>
    <rxActualBandwidth>763</rxActualBandwidth>
    <txActualBandwidth>267</txActualBandwidth>
    <localIpType>static</localIpType>
    <fixedIpV4>0.0.0.0</fixedIpV4>
    <fixedIpPrefixV4>0</fixedIpPrefixV4>
    <fixedIpV6>::10</fixedIpV6>
    <fixedIpPrefixV6>112</fixedIpPrefixV6>
    <floatingIpV4>0.0.0.0</floatingIpV4>
    <floatingIpV6>::</floatingIpV6>
  </ipInterfaceStatus>
...
</collection>
*/

type ipInterfaceStatusCollection struct {
	IPInterfaces []*ipInterfaceStatus `xml:"http://sonusnet.com/ns/mibs/SONUS-GEN2-IP-INTERFACE/1.0 ipInterfaceStatus,omitempty"`
}

type ipInterfaceStatus struct {
	Name              string  `xml:"http://sonusnet.com/ns/mibs/SONUS-GEN2-IP-INTERFACE/1.0 name"`
	OperState         string  `xml:"http://sonusnet.com/ns/mibs/SONUS-GEN2-IP-INTERFACE/1.0 operState"`
	RxPackets         float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GEN2-IP-INTERFACE/1.0 rxPackets"`
	TxPackets         float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GEN2-IP-INTERFACE/1.0 txPackets"`
	RxActualBandwidth float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GEN2-IP-INTERFACE/1.0 rxActualBandwidth"`
	TxActualBandwidth float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GEN2-IP-INTERFACE/1.0 txActualBandwidth"`
	NumMediaStreams   float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GEN2-IP-INTERFACE/1.0 numMediaStreams"`
}

func (i ipInterfaceStatus) operStateToMetric() float64 {
	switch i.OperState {
	case "resAllocated":
		return 0
	default:
		return 1
	}
}
