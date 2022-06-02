package exporter

import (
	"encoding/xml"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func (a *AddressContext) GetZoneStatusURL(ctx MetricContext) string {
	return fmt.Sprintf("%s/operational/addressContext/%s/zoneStatus/", ctx.APIBase, a.Name)
}

func (a *AddressContext) GetIPInterfaceGroupURL(ctx MetricContext) string {
	return fmt.Sprintf("%s/operational/addressContext/%s/ipInterfaceGroup/", ctx.APIBase, a.Name)
}

var ZoneStatusMetrics = map[string]*prometheus.Desc{
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

type AddressContext struct {
	Name              string
	Zones             []*ZoneStatus       `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 zoneStatus"`
	IPInterfaceGroups []*IPInterfaceGroup `xml:"http://sonusnet.com/ns/mibs/SONUS-GEN2-IP-INTERFACE/1.0 ipInterfaceGroup"`
}

type ZoneStatus struct {
	Name                 string  `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 name"`
	TotalCallsAvailable  float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 totalCallsAvailable"`
	InboundCallsUsage    float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 inboundCallsUsage"`
	OutboundCallsUsage   float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 outboundCallsUsage"`
	TotalCallsConfigured float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 totalCallsConfigured"`
	ActiveSipRegCount    float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 activeSipRegCount"`
}

type IPInterfaceGroup struct {
	Name        string `xml:"name"`
	IpInterface struct {
		Name string `xml:"name"`
	} `xml:"http://sonusnet.com/ns/mibs/SONUS-GEN2-IP-INTERFACE/1.0 ipInterface"`
}

func ProcessZones(addressContext *AddressContext, xmlBody *[]byte, ch chan<- prometheus.Metric) error {
	err := xml.Unmarshal(*xmlBody, &addressContext)
	if err != nil {
		log.Errorf("Failed to deserialize zoneStatus XML: %v", err)
		return err
	}

	for _, zone := range addressContext.Zones {
		ch <- prometheus.MustNewConstMetric(ZoneStatusMetrics["Zone_Total_Calls_Configured"], prometheus.GaugeValue, zone.TotalCallsConfigured, addressContext.Name, zone.Name)
		ch <- prometheus.MustNewConstMetric(ZoneStatusMetrics["Zone_Usage_Total"], prometheus.GaugeValue, zone.InboundCallsUsage, "inbound", addressContext.Name, zone.Name)
		ch <- prometheus.MustNewConstMetric(ZoneStatusMetrics["Zone_Usage_Total"], prometheus.GaugeValue, zone.OutboundCallsUsage, "outbound", addressContext.Name, zone.Name)
	}
	log.Info("Zone Status and Metrics collected")
	return nil
}

func ProcessIPInterfaceGroups(addressContext *AddressContext, xmlBody *[]byte, ch chan<- prometheus.Metric) error {
	err := xml.Unmarshal(*xmlBody, &addressContext)
	if err != nil {
		log.Errorf("Failed to deserialize ipInterfaceGroup XML: %v", err)
		return err
	}
	return nil
}
