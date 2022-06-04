package metrics

import (
	"encoding/xml"
	"strconv"
	"strings"

	"sonus-metrics-exporter/lib"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var FanMetric = lib.SonusMetric{
	Name:       "Fan",
	Processor:  processFans,
	URLGetter:  getFanUrl,
	APIMetrics: fanMetrics,
	Repetition: lib.RepeatNone,
}

const fanUrlSuffix = "/operational/system/fanStatus/"

func getFanUrl(ctx lib.MetricContext) string {
	return ctx.APIBase + fanUrlSuffix
}

var fanMetrics = map[string]*prometheus.Desc{
	"Fan_Speed": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "fan", "speed"),
		"Current speed of fans, in RPM",
		[]string{"server", "fanID"}, nil,
	),
}

func processFans(ctx lib.MetricContext, xmlBody *[]byte, ch chan<- prometheus.Metric, result chan<- lib.MetricResult) {
	fans := new(fanCollection)
	err := xml.Unmarshal(*xmlBody, &fans)
	if err != nil {
		log.Errorf("Failed to deserialize fanStatus XML: %v", err)
		result <- lib.MetricResult{Success: false, Errors: []*error{&err}}
		return
	}

	for _, fan := range fans.FanStatus {
		var fanRpm, err = fanStatus.speedToRPM(*fan)
		if err != nil {
			log.Errorf("Failed to convert fan speed (%q) to rpm: %v", fan.Speed, err)
			break
		}
		ch <- prometheus.MustNewConstMetric(fanMetrics["Fan_Speed"], prometheus.GaugeValue, fanRpm, fan.ServerName, fan.FanID)
	}
	log.Info("Fan Metrics collected")
	result <- lib.MetricResult{Success: true}
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

func (f fanStatus) speedToRPM() (float64, error) {
	var rpm = strings.TrimSuffix(f.Speed, " RPM")
	return strconv.ParseFloat(rpm, 64)
}
