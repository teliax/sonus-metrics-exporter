package metrics

import (
	"encoding/xml"
	"fmt"
	"sonus-metrics-exporter/lib"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	sipArsName      = "SIP ARS"
	sipArsURLFormat = "%s/operational/addressContext/%s/zone/%s/sipArsStatus/"
)

var SipArsMetric = lib.SonusMetric{
	Name:       sipArsName,
	Processor:  processSipArs,
	URLGetter:  getSipArsUrl,
	APIMetrics: sipArsMetrics,
	Repetition: lib.RepeatPerAddressContextZone,
}

func getSipArsUrl(ctx lib.MetricContext) string {
	return fmt.Sprintf(sipArsURLFormat, ctx.APIBase, ctx.AddressContext, ctx.Zone)
}

var sipArsMetrics = map[string]*prometheus.Desc{
	"SIPARS_Endpoint_State": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "sipars", "endpoint_status"),
		"State of a sipArs monitored endpoint",
		[]string{"zone", "endpoint_address", "endpoint_port", "state_name"}, nil,
	),
}

func processSipArs(ctx lib.MetricContext, xmlBody *[]byte) {
	var (
		errors []*error
		sipArs = new(sipArsCollection)
	)

	if len(*xmlBody) == 0 {
		// Empty response is fine, no need to attempt to parse and error out
		ctx.ResultChannel <- lib.MetricResult{Name: sipArsName, Success: true}
		return
	}

	err := xml.Unmarshal(*xmlBody, &sipArs)

	if err != nil {
		log.Errorf("Failed to deserialize sipArsStatus XML: %v", err)
		errors = append(errors, &err)
		ctx.ResultChannel <- lib.MetricResult{Name: sipArsName, Success: false, Errors: errors}
		return
	}

	for _, status := range sipArs.SipArsStatus {
		var endpoint string
		if status.EndpointDomainName != "" {
			endpoint = status.EndpointDomainName
		} else {
			endpoint = status.EndpointIpAddress
		}

		ctx.MetricChannel <- prometheus.MustNewConstMetric(sipArsMetrics["SIPARS_Endpoint_State"], prometheus.GaugeValue, status.stateToFloat(), ctx.Zone, endpoint, status.EndpointIpPortNum, status.EndpointArsState)
	}

	log.Infof("SIP ARS Metrics for Address Context %q, zone %q collected", ctx.AddressContext, ctx.Zone)
	ctx.ResultChannel <- lib.MetricResult{Name: sipArsName, Success: true, Errors: errors}
}

/*
<collection xmlns:y="http://tail-f.com/ns/rest">
  <sipArsStatus xmlns="http://sonusnet.com/ns/mibs/SONUS-SIP-TRUNK-GROUP/1.0">
    <sigZoneId>1000</sigZoneId>
    <recordIndex>1</recordIndex>
    <sigPortNum>1000</sigPortNum>
    <endpointDomainName></endpointDomainName>
    <endpointIpAddress>2600:1f14:cea3:7d4:7904:5b5a:1284:b7a0</endpointIpAddress>
    <endpointIpPortNum>5060</endpointIpPortNum>
    <endpointArsState>blacklisted</endpointArsState>
    <endpointStateTransitionTime>2022-09-08T13:05:24.000006+00:00</endpointStateTransitionTime>
  </sipArsStatus>
...
</collection>
*/

type sipArsCollection struct {
	SipArsStatus []*sipArsStatus `xml:"sipArsStatus"`
}

type sipArsStatus struct {
	SigZoneId               float64 `xml:"sigZoneId"`
	RecordIndex             float64 `xml:"recordIndex"`
	SigPortNum              float64 `xml:"sigPortNum"`
	EndpointDomainName      string  `xml:"endpointDomainName"`
	EndpointIpAddress       string  `xml:"endpointIpAddress"`
	EndpointIpPortNum       string  `xml:"endpointIpPortNum"`
	EndpointArsState        string  `xml:"endpointArsState"`
	EndpointStateTransition string  `xml:"endpointStateTransition"`
}

func (s sipArsStatus) stateToFloat() float64 {
	if s.EndpointArsState == "blacklisted" {
		return 1
	} else {
		return 0
	}
}
