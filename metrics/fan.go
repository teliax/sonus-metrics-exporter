package metrics

import (
	"encoding/xml"
	"strconv"
	"strings"

	"sonus-metrics-exporter/lib"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	fanName      = "Fan"
	fanUrlSuffix = "/operational/system/fanStatus/"
)

var FanMetric = lib.SonusMetric{
	Name:       fanName,
	Processor:  processFans,
	URLGetter:  getFanUrl,
	APIMetrics: fanMetrics,
	Repetition: lib.RepeatNone,
}

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

func processFans(ctx lib.MetricContext, xmlBody *[]byte) {
	var (
		errors []*error
		fans   = new(fanCollection)
	)

	err := xml.Unmarshal(*xmlBody, &fans)

	if err != nil {
		log.Errorf("Failed to deserialize fanStatus XML: %v", err)
		errors = append(errors, &err)
		ctx.ResultChannel <- lib.MetricResult{Name: fanName, Success: false, Errors: errors}
		return
	}

	for _, fan := range fans.FanStatus {
		var fanRpm, err = fanStatus.speedToRPM(*fan)
		if err != nil {
			log.Errorf("Failed to convert fan speed (%q) to rpm: %v", fan.Speed, err)
			errors = append(errors, &err)
			break
		}
		ctx.MetricChannel <- prometheus.MustNewConstMetric(fanMetrics["Fan_Speed"], prometheus.GaugeValue, fanRpm, fan.ServerName, fan.FanID)
	}

	log.Info("Fan Metrics collected")
	ctx.ResultChannel <- lib.MetricResult{Name: fanName, Success: true, Errors: errors}
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
