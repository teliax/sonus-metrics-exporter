package exporter

import (
	"encoding/xml"
	"fmt"

	"sonus-metrics-exporter/lib"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func getServerStatusURL(apiBase string) string {
	return fmt.Sprintf("%s/operational/system/serverStatus/", apiBase)
}

func (a *addressContext) getZoneStatusURL(ctx lib.MetricContext) string {
	return fmt.Sprintf("%s/operational/addressContext/%s/zoneStatus/", ctx.APIBase, a.Name)
}

func (a *addressContext) getIPInterfaceGroupURL(ctx lib.MetricContext) string {
	return fmt.Sprintf("%s/operational/addressContext/%s/ipInterfaceGroup/", ctx.APIBase, a.Name)
}

var zoneStatusMetrics = map[string]*prometheus.Desc{
	"Zone_Total_Calls_Configured": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "zone", "total_calls_configured"),
		"Total call limit per zone",
		[]string{"addresscontext", "zone"}, nil,
	),
	"Zone_Usage_Total": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "zone", "usage_total"),
		"Total call limit per zone",
		[]string{"direction", "addresscontext", "zone"}, nil,
	),
}

type addressContext struct {
	Name              string
	Zones             []*zoneStatus       `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 zoneStatus"`
	IPInterfaceGroups []*ipInterfaceGroup `xml:"http://sonusnet.com/ns/mibs/SONUS-GEN2-IP-INTERFACE/1.0 ipInterfaceGroup"`
}

type zoneStatus struct {
	Name                 string  `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 name"`
	TotalCallsAvailable  float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 totalCallsAvailable"`
	InboundCallsUsage    float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 inboundCallsUsage"`
	OutboundCallsUsage   float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 outboundCallsUsage"`
	TotalCallsConfigured float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 totalCallsConfigured"`
	ActiveSipRegCount    float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 activeSipRegCount"`
}

type ipInterfaceGroup struct {
	Name         string `xml:"name"`
	IpInterfaces struct {
		Name string `xml:"name"`
	} `xml:"http://sonusnet.com/ns/mibs/SONUS-GEN2-IP-INTERFACE/1.0 ipInterface"`
}

func processZones(addressContext *addressContext, xmlBody *[]byte, ch chan<- prometheus.Metric) error {
	err := xml.Unmarshal(*xmlBody, &addressContext)
	if err != nil {
		log.Errorf("Failed to deserialize zoneStatus XML: %v", err)
		return err
	}

	for _, zone := range addressContext.Zones {
		ch <- prometheus.MustNewConstMetric(zoneStatusMetrics["Zone_Total_Calls_Configured"], prometheus.GaugeValue, zone.TotalCallsConfigured, addressContext.Name, zone.Name)
		ch <- prometheus.MustNewConstMetric(zoneStatusMetrics["Zone_Usage_Total"], prometheus.GaugeValue, zone.InboundCallsUsage, "inbound", addressContext.Name, zone.Name)
		ch <- prometheus.MustNewConstMetric(zoneStatusMetrics["Zone_Usage_Total"], prometheus.GaugeValue, zone.OutboundCallsUsage, "outbound", addressContext.Name, zone.Name)
	}
	log.Info("Zone Status and Metrics collected")
	return nil
}

func processIPInterfaceGroups(addressContext *addressContext, xmlBody *[]byte) error {
	err := xml.Unmarshal(*xmlBody, &addressContext)
	if err != nil {
		log.Errorf("Failed to deserialize ipInterfaceGroup XML: %v", err)
		return err
	}
	return nil
}
