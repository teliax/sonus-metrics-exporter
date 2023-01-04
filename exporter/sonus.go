package exporter

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"time"

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

var serverStatusMetrics = map[string]*prometheus.Desc{
	"System_Redundancy_Role": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "system", "redundancy_role"),
		"Current role of server. 1 = active",
		[]string{"server", "role_name"}, nil,
	),
	"System_Sync_Status": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "system", "sync_status"),
		"Current synchronization status. 1 = syncCompleted",
		[]string{"server", "status_name"}, nil,
	),
	"System_Uptime": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "system", "uptime"),
		"Current uptime of server, in seconds",
		[]string{"server", "type"}, nil,
	),
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

type (
	serverStatusCollection struct {
		ServerStatus []*serverStatus `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 serverStatus"`
	}

	serverStatus struct {
		Name                     string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 name"`
		SerialNum                string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 serialNum"`
		ManagementRedundancyRole string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 mgmtRedundancyRole"`
		Uptime                   string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 upTime"`
		ApplicationUptime        string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 applicationUpTime"`
		SyncStatus               string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 syncStatus"`
	}

	serverUptimeType uint8
)

const (
	serverOSUptime serverUptimeType = iota
	serverAppUptime
)

var uptimeRegex = regexp.MustCompile(`^(\d{1,5}) Days (\d{2}):(\d{2}):(\d{2})$`)

func (s serverStatus) parseUptime(upType serverUptimeType) float64 {
	var (
		uptime    time.Duration
		uptimeStr string
	)

	switch upType {
	case serverOSUptime:
		uptimeStr = s.Uptime
	case serverAppUptime:
		uptimeStr = s.ApplicationUptime
	}

	uptimeFields := uptimeRegex.FindStringSubmatch(uptimeStr)

	if len(uptimeFields) == 5 {
		daysInt, _ := strconv.ParseInt(uptimeFields[1], 10, 32)
		hoursInt, _ := strconv.ParseInt(uptimeFields[2], 10, 32)
		minutesInt, _ := strconv.ParseInt(uptimeFields[3], 10, 32)
		secondsInt, _ := strconv.ParseInt(uptimeFields[4], 10, 32)

		uptime += time.Duration(daysInt) * time.Hour * 24
		uptime += time.Duration(hoursInt) * time.Hour
		uptime += time.Duration(minutesInt) * time.Minute
		uptime += time.Duration(secondsInt) * time.Second

		return uptime.Seconds()
	} else {
		log.Errorf("Unable to match uptime %q with regex.", uptimeStr)
		return 0
	}
}

func (s serverStatus) mgmtRedunRoleToFloat() float64 {
	switch s.ManagementRedundancyRole {
	case "active":
		return 1
	default:
		return 0
	}
}

func (s serverStatus) syncStatusToFloat() float64 {
	switch s.SyncStatus {
	case "syncCompleted":
		return 1
	default:
		return 0
	}
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

func processServerStatus(xmlBody *[]byte, ch chan<- prometheus.Metric) error {
	var serverStatuses = new(serverStatusCollection)

	err := xml.Unmarshal(*xmlBody, &serverStatuses)
	if err != nil {
		log.Errorf("Failed to deserialize serverStatus XML: %v", err)
		return err
	}

	for _, server := range serverStatuses.ServerStatus {
		ch <- prometheus.MustNewConstMetric(serverStatusMetrics["System_Redundancy_Role"], prometheus.GaugeValue, server.mgmtRedunRoleToFloat(), server.Name, server.ManagementRedundancyRole)
		ch <- prometheus.MustNewConstMetric(serverStatusMetrics["System_Sync_Status"], prometheus.GaugeValue, server.syncStatusToFloat(), server.Name, server.SyncStatus)
		ch <- prometheus.MustNewConstMetric(serverStatusMetrics["System_Uptime"], prometheus.CounterValue, server.parseUptime(serverOSUptime), server.Name, "os")
		ch <- prometheus.MustNewConstMetric(serverStatusMetrics["System_Uptime"], prometheus.CounterValue, server.parseUptime(serverAppUptime), server.Name, "application")
	}
	log.Info("Server Status and Metrics collected")
	return nil
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
