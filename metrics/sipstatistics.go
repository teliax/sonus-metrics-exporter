package metrics

import (
	"encoding/xml"
	"fmt"

	"sonus-metrics-exporter/lib"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	sipStatisticsName      = "SipStatistic"
	sipStatisticsURLFormat = "%s/operational/addressContext/%s/zone/%s/sipCurrentStatistics/"
)

var SipStatisticMetric = lib.SonusMetric{
	Name:       sipStatisticsName,
	Processor:  processSipStatistics,
	URLGetter:  getSipStatisticsUrl,
	APIMetrics: sipStatisticMetrics,
	Repetition: lib.RepeatPerAddressContextZone,
}

func getSipStatisticsUrl(ctx lib.MetricContext) string {
	return fmt.Sprintf(sipStatisticsURLFormat, ctx.APIBase, ctx.AddressContext, ctx.Zone)
}

var sipStatisticMetrics = map[string]*prometheus.Desc{
	"TG_SIP_Req_Sent": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "TG", "sip_req_sent"),
		"Number of SIP requests sent",
		[]string{"zone", "name", "method"}, nil,
	),
	"TG_SIP_Req_Received": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "TG", "sip_req_recv"),
		"Number of SIP requests received",
		[]string{"zone", "name", "method"}, nil,
	),
	"TG_SIP_Resp_Sent": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "TG", "sip_resp_sent"),
		"Number of SIP responses sent",
		[]string{"zone", "name", "code"}, nil,
	),
	"TG_SIP_Resp_Received": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "TG", "sip_resp_recv"),
		"Number of SIP responses received",
		[]string{"zone", "name", "code"}, nil,
	),
}

func processSipStatistics(ctx lib.MetricContext, xmlBody *[]byte) {
	var (
		errors   []*error
		sipStats = new(sipStatisticCollection)
	)

	if len(*xmlBody) == 0 {
		ctx.ResultChannel <- lib.MetricResult{Name: sipStatisticsName, Success: true}
		return
	}

	err := xml.Unmarshal(*xmlBody, &sipStats)

	if err != nil {
		log.Errorf("Failed to deserialize sipCurrentStatistics XML: %v", err)
		errors = append(errors, &err)
		ctx.ResultChannel <- lib.MetricResult{Name: sipStatisticsName, Success: false, Errors: errors}
		return
	}

	for _, sipStat := range sipStats.SipStatistics {
		var sipReqSent = map[string]float64{
			"INVITE":             sipStat.SndInvite,
			"PRACK":              sipStat.SndPrack,
			"INFO":               sipStat.SndInfo,
			"REFER":              sipStat.SndRefer,
			"BYE":                sipStat.SndBye,
			"CANCEL":             sipStat.SndCancel,
			"REGISTER":           sipStat.SndRegister,
			"UPDATE":             sipStat.SndUpdate,
			"SUBSCRIBE":          sipStat.SndSubscriber,
			"NOTIFY":             sipStat.SndNotify,
			"OPTIONS":            sipStat.SndOption,
			"MESSAGE":            sipStat.SndMessage,
			"PUBLISH":            sipStat.SndPublish,
			"INVITE (retrans)":   sipStat.InvReTransmit,
			"REGISTER (retrans)": sipStat.RegReTransmit,
			"BYE (retrans)":      sipStat.ByeReTransmit,
			"CANCEL (retrans)":   sipStat.CancelReTransmit,
			"Other (retrans)":    sipStat.OtherReTransmit,
		}
		for n, v := range sipReqSent {
			ctx.MetricChannel <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, v, ctx.Zone, sipStat.TrunkGroupName, n)
		}

		var sipReqReceived = map[string]float64{
			"INVITE":    sipStat.RcvInvite,
			"PRACK":     sipStat.RcvPrack,
			"INFO":      sipStat.RcvInfo,
			"REFER":     sipStat.RcvRefer,
			"BYE":       sipStat.RcvBye,
			"CANCEL":    sipStat.RcvCancel,
			"REGISTER":  sipStat.RcvRegister,
			"UPDATE":    sipStat.RcvUpdate,
			"SUBSCRIBE": sipStat.RcvSubscriber,
			"NOTIFY":    sipStat.RcvNotify,
			"OPTIONS":   sipStat.RcvOption,
			"MESSAGE":   sipStat.RcvMessage,
			"PUBLISH":   sipStat.RcvPublish,
			"Unknown":   sipStat.RcvUnknownMsg,
		}
		for n, v := range sipReqReceived {
			ctx.MetricChannel <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, v, ctx.Zone, sipStat.TrunkGroupName, n)
		}

		var sipRespSent = map[string]float64{
			"ACK":              sipStat.SndAck,
			"18x":              sipStat.Snd18x,
			"1xx":              sipStat.Snd1xx,
			"2xx":              sipStat.Snd2xx,
			"Non-INVITE 2xx":   sipStat.SndNonInv2xx,
			"3xx":              sipStat.Snd3xx,
			"4xx":              sipStat.Snd4xx,
			"5xx":              sipStat.Snd5xx,
			"6xx":              sipStat.Snd6xx,
			"Non-INVITE error": sipStat.SndNonInvErr,
		}
		for n, v := range sipRespSent {
			ctx.MetricChannel <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Sent"], prometheus.CounterValue, v, ctx.Zone, sipStat.TrunkGroupName, n)
		}

		var sipRespReceived = map[string]float64{
			"ACK":              sipStat.RcvAck,
			"18x":              sipStat.Rcv18x,
			"1xx":              sipStat.Rcv1xx,
			"2xx":              sipStat.Rcv2xx,
			"Non-INVITE 2xx":   sipStat.RcvNonInv2xx,
			"3xx":              sipStat.Rcv3xx,
			"4xx":              sipStat.Rcv4xx,
			"5xx":              sipStat.Rcv5xx,
			"6xx":              sipStat.Rcv6xx,
			"Non-INVITE error": sipStat.RcvNonInvErr,
		}
		for n, v := range sipRespReceived {
			ctx.MetricChannel <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Received"], prometheus.CounterValue, v, ctx.Zone, sipStat.TrunkGroupName, n)
		}
	}
	log.Infof("SIP Statistics Metrics for Address Context %q, zone %q collected", ctx.AddressContext, ctx.Zone)
	ctx.ResultChannel <- lib.MetricResult{Name: sipStatisticsName, Success: true}
}

/*
<collection xmlns:y="http://tail-f.com/ns/rest">
  <sipCurrentStatistics xmlns="http://sonusnet.com/ns/mibs/SONUS-SIP-PEER-PERF-STATS/1.0">
    <name>TEST_Logan</name>
    <rcvInvite>0</rcvInvite>
    <sndInvite>0</sndInvite>
    <rcvAck>0</rcvAck>
    <sndAck>0</sndAck>
    <rcvPrack>0</rcvPrack>
...
</collection>
*/

type sipStatisticCollection struct {
	SipStatistics []*sipStatistics `xml:"http://sonusnet.com/ns/mibs/SONUS-SIP-PEER-PERF-STATS/1.0 sipCurrentStatistics,omitempty"`
}

type sipStatistics struct {
	TrunkGroupName                string  `xml:"name"`
	RcvInvite                     float64 `xml:"rcvInvite"`
	SndInvite                     float64 `xml:"sndInvite"`
	RcvAck                        float64 `xml:"rcvAck"`
	SndAck                        float64 `xml:"sndAck"`
	RcvPrack                      float64 `xml:"rcvPrack"`
	SndPrack                      float64 `xml:"sndPrack"`
	RcvInfo                       float64 `xml:"rcvInfo"`
	SndInfo                       float64 `xml:"sndInfo"`
	RcvRefer                      float64 `xml:"rcvRefer"`
	SndRefer                      float64 `xml:"sndRefer"`
	RcvBye                        float64 `xml:"rcvBye"`
	SndBye                        float64 `xml:"sndBye"`
	RcvCancel                     float64 `xml:"rcvCancel"`
	SndCancel                     float64 `xml:"sndCancel"`
	RcvRegister                   float64 `xml:"rcvRegister"`
	SndRegister                   float64 `xml:"sndRegister"`
	RcvUpdate                     float64 `xml:"rcvUpdate"`
	SndUpdate                     float64 `xml:"sndUpdate"`
	Rcv18x                        float64 `xml:"rcv18x"`
	Snd18x                        float64 `xml:"snd18x"`
	Rcv1xx                        float64 `xml:"rcv1xx"`
	Snd1xx                        float64 `xml:"snd1xx"`
	Rcv2xx                        float64 `xml:"rcv2xx"`
	Snd2xx                        float64 `xml:"snd2xx"`
	RcvNonInv2xx                  float64 `xml:"rcvNonInv2xx"`
	SndNonInv2xx                  float64 `xml:"sndNonInv2xx"`
	Rcv3xx                        float64 `xml:"rcv3xx"`
	Snd3xx                        float64 `xml:"snd3xx"`
	Rcv4xx                        float64 `xml:"rcv4xx"`
	Snd4xx                        float64 `xml:"snd4xx"`
	Rcv5xx                        float64 `xml:"rcv5xx"`
	Snd5xx                        float64 `xml:"snd5xx"`
	Rcv6xx                        float64 `xml:"rcv6xx"`
	Snd6xx                        float64 `xml:"snd6xx"`
	RcvNonInvErr                  float64 `xml:"rcvNonInvErr"`
	SndNonInvErr                  float64 `xml:"sndNonInvErr"`
	RcvUnknownMsg                 float64 `xml:"rcvUnknownMsg"`
	RcvSubscriber                 float64 `xml:"rcvSubscriber"`
	SndSubscriber                 float64 `xml:"sndSubscriber"`
	RcvNotify                     float64 `xml:"rcvNotify"`
	SndNotify                     float64 `xml:"sndNotify"`
	RcvOption                     float64 `xml:"rcvOption"`
	SndOption                     float64 `xml:"sndOption"`
	InvReTransmit                 float64 `xml:"invReTransmit"`
	RegReTransmit                 float64 `xml:"regReTransmit"`
	ByeReTransmit                 float64 `xml:"byeRetransmit"`
	CancelReTransmit              float64 `xml:"cancelReTransmit"`
	OtherReTransmit               float64 `xml:"otherReTransmit"`
	RcvMessage                    float64 `xml:"rcvMessage"`
	SndMessage                    float64 `xml:"sndMessage"`
	RcvPublish                    float64 `xml:"rcvPublish"`
	SndPublish                    float64 `xml:"sndPublish"`
	EmergencyAccept               float64 `xml:"emergencyAccept"`
	EmergencyRejectBWCall         float64 `xml:"emergencyRejectBWCall"`
	EmergencyRejectPolicer        float64 `xml:"emergencyRejectPolicer"`
	HpcAccept                     float64 `xml:"hpcAccept"`
	NumberOfCallsSendingAARs      float64 `xml:"numberOfCallsSendingAARs"`
	NumberOfReceivedAAAFailures   float64 `xml:"numberOfReceivedAAAFailures"`
	NumberOfTotalAARSent          float64 `xml:"numberOfTotalAARSent"`
	NumberOfTimeoutOrErrorAAR     float64 `xml:"numberOfTimeoutOrErrorAAR"`
	EmergencyRegAccept            float64 `xml:"emergencyRegAccept"`
	EmergencyRegRejectLimit       float64 `xml:"emergencyRegRejectLimit"`
	EmergencyRegRejectPolicer     float64 `xml:"emergencyRegRejectPolicer"`
	NumberOfReceivedAAASuccesses  float64 `xml:"numberOfReceivedAAASuccesses"`
	NumberOfReceivedRARs          float64 `xml:"numberOfReceivedRARs"`
	NumberOfReceivedASRs          float64 `xml:"numberOfReceivedASRs"`
	NumberOfSentSTRs              float64 `xml:"numberOfSentSTRs"`
	EmergencyOODAccept            float64 `xml:"emergencyOODAccept"`
	EmergencyOODRejectPolicer     float64 `xml:"emergencyOODRejectPolicer"`
	EmergencySubsAccept           float64 `xml:"emergencySubsAccept"`
	EmergencySubsRejectLimit      float64 `xml:"emergencySubsRejectLimit"`
	EmergencySubsRejectPolicer    float64 `xml:"emergencySubsRejectPolicer"`
	ParseError                    float64 `xml:"parseError"`
	NumberOfTotalUDRSent          float64 `xml:"numberOfTotalUDRSent"`
	NumberOfTimeoutOrErrorUDR     float64 `xml:"numberOfTimeoutOrErrorUDR"`
	NumberOfReceivedUDASuccesses  float64 `xml:"numberOfReceivedUDASuccesses"`
	NumberOfReceivedUDAFailures   float64 `xml:"numberOfReceivedUDAFailures"`
	TotNumOfS8hrOutbndReg         float64 `xml:"totNumOfS8hrOutbndReg"`
	NumOfS8hrOutbndRegSuc         float64 `xml:"numOfS8hrOutbndRegSuc"`
	NumOfS8hrOutbndRegFail        float64 `xml:"numOfS8hrOutbndRegFail"`
	TotNumOfS8hrOutbndNormalCall  float64 `xml:"totNumOfS8hrOutbndNormalCall"`
	NumOfS8hrOutbndNormalCallSuc  float64 `xml:"numOfS8hrOutbndNormalCallSuc"`
	NumOfS8hrOutbndNormalCallFail float64 `xml:"numOfS8hrOutbndNormalCallFail"`
	NumOfS8hrOutbndEmgCallRej     float64 `xml:"numOfS8hrOutbndEmgCallRej"`
	NumOfS8hrInboundRegSuc        float64 `xml:"numOfS8hrInboundRegSuc"`
	NumOfS8hrInboundRegFail       float64 `xml:"numOfS8hrInboundRegFail"`
	NumOfS8hrInboundEmgCallSuc    float64 `xml:"numOfS8hrInboundEmgCallSuc"`
	NumOfS8hrInboundEmgCallFail   float64 `xml:"numOfS8hrInboundEmgCallFail"`
	InHpcAccept                   float64 `xml:"inHpcAccept"`
	OutHpcAccept                  float64 `xml:"outHpcAccept"`
	Hpc403Out                     float64 `xml:"hpc403Out"`
	HpcOverloadExempt             float64 `xml:"hpcOverloadExempt"`
}
