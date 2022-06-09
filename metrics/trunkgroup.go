package metrics

import (
	"encoding/xml"

	"sonus-metrics-exporter/lib"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	trunkGroupName      = "TrunkGroup"
	trunkGroupUrlSuffix = "/operational/global/globalTrunkGroupStatus/"
)

var TGMetric = lib.SonusMetric{
	Name:       trunkGroupName,
	Processor:  processTGs,
	URLGetter:  getTGUrl,
	APIMetrics: tgMetrics,
	Repetition: lib.RepeatNone,
}

func getTGUrl(ctx lib.MetricContext) string {
	return ctx.APIBase + trunkGroupUrlSuffix
}

var tgMetrics = map[string]*prometheus.Desc{
	"TG_Bandwidth": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "TG", "bytes"),
		"Bandwidth in use by current calls",
		[]string{"zone", "name", "direction"}, nil,
	),
	"TG_OBState": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "TG", "outbound_state"),
		"State of outbound calls on the trunkgroup",
		[]string{"zone", "name"}, nil,
	),
	"TG_State": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "TG", "state"),
		"State of the trunkgroup",
		[]string{"zone", "name"}, nil,
	),
	"TG_TotalChans": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "TG", "total_channels"),
		"Number of configured channels",
		[]string{"zone", "name"}, nil,
	),
	"TG_Usage": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "TG", "usage_total"),
		"Number of active calls",
		[]string{"zone", "name", "direction"}, nil,
	),
}

func processTGs(ctx lib.MetricContext, xmlBody *[]byte, ch chan<- prometheus.Metric, result chan<- lib.MetricResult) {
	var (
		errors []*error
		tgs    = new(trunkGroupCollection)
	)

	err := xml.Unmarshal(*xmlBody, &tgs)

	if err != nil {
		log.Errorf("Failed to deserialize globalTrunkGroupStatus XML: %v", err)
		errors = append(errors, &err)
		result <- lib.MetricResult{Name: trunkGroupName, Success: false, Errors: errors}
		return
	}

	for _, tg := range tgs.TrunkGroupStatus {
		ch <- prometheus.MustNewConstMetric(tgMetrics["TG_Usage"], prometheus.GaugeValue, tg.InboundCallsUsage, tg.Zone, tg.Name, "inbound")
		ch <- prometheus.MustNewConstMetric(tgMetrics["TG_Usage"], prometheus.GaugeValue, tg.OutboundCallsUsage, tg.Zone, tg.Name, "outbound")
		ch <- prometheus.MustNewConstMetric(tgMetrics["TG_Bandwidth"], prometheus.GaugeValue, tg.BandwidthInboundUsage, tg.Zone, tg.Name, "inbound")
		ch <- prometheus.MustNewConstMetric(tgMetrics["TG_Bandwidth"], prometheus.GaugeValue, tg.BandwidthOutboundUsage, tg.Zone, tg.Name, "outbound")
		ch <- prometheus.MustNewConstMetric(tgMetrics["TG_TotalChans"], prometheus.GaugeValue, tg.TotalCallsConfigured, tg.Zone, tg.Name)
		ch <- prometheus.MustNewConstMetric(tgMetrics["TG_State"], prometheus.GaugeValue, trunkGroupStatus.stateToMetric(*tg), tg.Zone, tg.Name)
		ch <- prometheus.MustNewConstMetric(tgMetrics["TG_OBState"], prometheus.GaugeValue, trunkGroupStatus.outStateToMetric(*tg), tg.Zone, tg.Name)
	}

	log.Info("Trunk Group Metrics collected")
	result <- lib.MetricResult{Name: trunkGroupName, Success: true}
}

/*
<collection xmlns:y="http://tail-f.com/ns/rest">
  <globalTrunkGroupStatus xmlns="http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0">
    <name>TEST</name>
    <state>inService</state>
    <totalCallsAvailable>100</totalCallsAvailable>
    <totalCallsInboundReserved>0</totalCallsInboundReserved>
    <inboundCallsUsage>0</inboundCallsUsage>
    <outboundCallsUsage>0</outboundCallsUsage>
    <totalCallsConfigured>100</totalCallsConfigured>
    <priorityCallUsage>0</priorityCallUsage>
    <totalOutboundCallsReserved>0</totalOutboundCallsReserved>
    <bwCurrentLimit>-1</bwCurrentLimit>
    <bwAvailable>-1</bwAvailable>
    <bwInboundUsage>0</bwInboundUsage>
    <bwOutboundUsage>0</bwOutboundUsage>
    <packetOutDetectState>normal</packetOutDetectState>
    <addressContext>default</addressContext>
    <zone>zone_23</zone>
    <priorityBwUsage>0</priorityBwUsage>
  </globalTrunkGroupStatus>
  ...
</collection>
*/

type trunkGroupCollection struct {
	TrunkGroupStatus []*trunkGroupStatus `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 globalTrunkGroupStatus,omitempty"`
}

type trunkGroupStatus struct {
	AddressContext             string  `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 addressContext"`
	BandwidthAvailable         float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 bwAvailable"`
	BandwidthCurrentLimit      float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 bwCurrentLimit"`
	BandwidthInboundUsage      float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 bwInboundUsage"`
	BandwidthOutboundUsage     float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 bwOutboundUsage"`
	InboundCallsUsage          float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 inboundCallsUsage"`
	Name                       string  `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 name"`
	OutboundCallsUsage         float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 outboundCallsUsage"`
	PacketOutDetectState       string  `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 packetOutDetectState"`
	PriorityBwUsage            float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 priorityBwUsage"`
	PriorityCallUsage          float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 priorityCallUsage"`
	State                      string  `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 state"`
	TotalCallsAvailable        float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 totalCallsAvailable"`
	TotalCallsConfigured       float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 totalCallsConfigured"`
	TotalCallsInboundReserved  float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 totalCallsInboundReserved"`
	TotalOutboundCallsReserved float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 totalOutboundCallsReserved"`
	Zone                       string  `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 zone"`
}

func (t trunkGroupStatus) stateToMetric() float64 {
	switch t.State {
	case "inService":
		return 1
	default:
		return 0
	}
}

func (t trunkGroupStatus) outStateToMetric() float64 {
	switch t.PacketOutDetectState {
	case "normal":
		return 1
	default:
		return 0
	}
}
