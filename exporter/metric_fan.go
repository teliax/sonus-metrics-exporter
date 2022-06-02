package exporter

import (
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const fanUrlSuffix = "/operational/system/fanStatus/"

func GetFanUrl(ctx MetricContext) string {
	return ctx.APIBase + fanUrlSuffix
}

var FanMetrics = map[string]*prometheus.Desc{
	"Fan_Speed": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "fan", "speed"),
		"Current speed of fans, in RPM",
		[]string{"server", "fanID"}, nil,
	),
}

func ProcessFans(ctx MetricContext, xmlBody *[]byte, ch chan<- prometheus.Metric, result chan<- bool) {
	fans := new(fanCollection)
	err := xml.Unmarshal(*xmlBody, &fans)
	if err != nil {
		log.Errorf("Failed to deserialize fanStatus XML: %v", err)
		result <- false
		return
	}

	for _, fan := range fans.FanStatus {
		var fanRpm, err = fanStatus.SpeedToRPM(*fan)
		if err != nil {
			log.Errorf("Failed to convert fan speed (%q) to rpm: %v", fan.Speed, err)
			break
		}
		ch <- prometheus.MustNewConstMetric(FanMetrics["Fan_Speed"], prometheus.GaugeValue, fanRpm, fan.ServerName, fan.FanID)
	}
	log.Info("Fan Metrics collected")
	result <- true
}

/*
<collection xmlns:y="http://tail-f.com/ns/rest">
  <fanStatus xmlns="http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0">
    <serverName>densbc01a</serverName>
    <fanId>FAN1/BOT</fanId>
    <speed>5632 RPM</speed>
  </fanStatus>
...
</collection>
*/

type fanCollection struct {
	FanStatus []*fanStatus `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 fanStatus,omitempty"`
}

type fanStatus struct {
	ServerName string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 serverName"`
	FanID      string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 fanId"`
	Speed      string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 speed"`
}

func (f fanStatus) SpeedToRPM() (float64, error) {
	var rpm = strings.TrimSuffix(f.Speed, " RPM")
	var rpmFloat, err = strconv.ParseFloat(rpm, 64)
	return rpmFloat, err
}
