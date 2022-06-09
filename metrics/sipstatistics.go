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

func processSipStatistics(ctx lib.MetricContext, xmlBody *[]byte, ch chan<- prometheus.Metric, result chan<- lib.MetricResult) {
	var (
		errors   []*error
		sipStats = new(sipStatisticCollection)
	)

	if len(*xmlBody) == 0 {
		result <- lib.MetricResult{Name: sipStatisticsName, Success: true}
		return
	}

	err := xml.Unmarshal(*xmlBody, &sipStats)

	if err != nil {
		log.Errorf("Failed to deserialize sipCurrentStatistics XML: %v", err)
		errors = append(errors, &err)
		result <- lib.MetricResult{Name: sipStatisticsName, Success: false, Errors: errors}
		return
	}

	for _, sipStat := range sipStats.SipStatistics {
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.SndInvite, ctx.Zone, sipStat.Name, "INVITE")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.SndPrack, ctx.Zone, sipStat.Name, "PRACK")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.SndInfo, ctx.Zone, sipStat.Name, "INFO")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.SndRefer, ctx.Zone, sipStat.Name, "REFER")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.SndBye, ctx.Zone, sipStat.Name, "BYE")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.SndCancel, ctx.Zone, sipStat.Name, "CANCEL")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.SndRegister, ctx.Zone, sipStat.Name, "REGISTER")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.SndUpdate, ctx.Zone, sipStat.Name, "UPDATE")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.SndSubscriber, ctx.Zone, sipStat.Name, "SUBSCRIBE") // Is this correct? "subscriber"?
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.SndNotify, ctx.Zone, sipStat.Name, "NOTIFY")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.SndOption, ctx.Zone, sipStat.Name, "OPTIONS")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.SndMessage, ctx.Zone, sipStat.Name, "MESSAGE")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.SndPublish, ctx.Zone, sipStat.Name, "PUBLISH")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.InvReTransmit, ctx.Zone, sipStat.Name, "INVITE (retrans)")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.RegReTransmit, ctx.Zone, sipStat.Name, "REGISTER (retrans)")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.ByeReTransmit, ctx.Zone, sipStat.Name, "BYE (retrans)")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.CancelReTransmit, ctx.Zone, sipStat.Name, "CANCEL (retrans)")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Sent"], prometheus.CounterValue, sipStat.OtherReTransmit, ctx.Zone, sipStat.Name, "Other (retrans)")

		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, sipStat.RcvInvite, ctx.Zone, sipStat.Name, "INVITE")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, sipStat.RcvPrack, ctx.Zone, sipStat.Name, "PRACK")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, sipStat.RcvInfo, ctx.Zone, sipStat.Name, "INFO")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, sipStat.RcvRefer, ctx.Zone, sipStat.Name, "REFER")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, sipStat.RcvBye, ctx.Zone, sipStat.Name, "BYE")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, sipStat.RcvCancel, ctx.Zone, sipStat.Name, "CANCEL")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, sipStat.RcvRegister, ctx.Zone, sipStat.Name, "REGISTER")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, sipStat.RcvUpdate, ctx.Zone, sipStat.Name, "UPDATE")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, sipStat.RcvSubscriber, ctx.Zone, sipStat.Name, "SUBSCRIBE") // Is this correct? "subscriber"?
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, sipStat.RcvNotify, ctx.Zone, sipStat.Name, "NOTIFY")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, sipStat.RcvOption, ctx.Zone, sipStat.Name, "OPTIONS")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, sipStat.RcvMessage, ctx.Zone, sipStat.Name, "MESSAGE")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, sipStat.RcvPublish, ctx.Zone, sipStat.Name, "PUBLISH")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Req_Received"], prometheus.CounterValue, sipStat.RcvUnknownMsg, ctx.Zone, sipStat.Name, "Unknown")

		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Sent"], prometheus.CounterValue, sipStat.SndAck, ctx.Zone, sipStat.Name, "ACK")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Sent"], prometheus.CounterValue, sipStat.Snd18x, ctx.Zone, sipStat.Name, "18x")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Sent"], prometheus.CounterValue, sipStat.Snd1xx, ctx.Zone, sipStat.Name, "1xx")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Sent"], prometheus.CounterValue, sipStat.Snd2xx, ctx.Zone, sipStat.Name, "2xx")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Sent"], prometheus.CounterValue, sipStat.SndNonInv2xx, ctx.Zone, sipStat.Name, "Non-INVITE 2xx")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Sent"], prometheus.CounterValue, sipStat.Snd3xx, ctx.Zone, sipStat.Name, "3xx")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Sent"], prometheus.CounterValue, sipStat.Snd4xx, ctx.Zone, sipStat.Name, "4xx")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Sent"], prometheus.CounterValue, sipStat.Snd5xx, ctx.Zone, sipStat.Name, "5xx")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Sent"], prometheus.CounterValue, sipStat.Snd6xx, ctx.Zone, sipStat.Name, "6xx")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Sent"], prometheus.CounterValue, sipStat.SndNonInvErr, ctx.Zone, sipStat.Name, "Non-INVITE error")

		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Received"], prometheus.CounterValue, sipStat.RcvAck, ctx.Zone, sipStat.Name, "ACK")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Received"], prometheus.CounterValue, sipStat.Rcv18x, ctx.Zone, sipStat.Name, "18x")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Received"], prometheus.CounterValue, sipStat.Rcv1xx, ctx.Zone, sipStat.Name, "1xx")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Received"], prometheus.CounterValue, sipStat.Rcv2xx, ctx.Zone, sipStat.Name, "2xx")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Received"], prometheus.CounterValue, sipStat.RcvNonInv2xx, ctx.Zone, sipStat.Name, "Non-INVITE 2xx")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Received"], prometheus.CounterValue, sipStat.Rcv3xx, ctx.Zone, sipStat.Name, "3xx")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Received"], prometheus.CounterValue, sipStat.Rcv4xx, ctx.Zone, sipStat.Name, "4xx")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Received"], prometheus.CounterValue, sipStat.Rcv5xx, ctx.Zone, sipStat.Name, "5xx")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Received"], prometheus.CounterValue, sipStat.Rcv6xx, ctx.Zone, sipStat.Name, "6xx")
		ch <- prometheus.MustNewConstMetric(sipStatisticMetrics["TG_SIP_Resp_Received"], prometheus.CounterValue, sipStat.RcvNonInvErr, ctx.Zone, sipStat.Name, "Non-INVITE error")
	}
	log.Infof("SIP Statistics Metrics for Address Context %q, zone %q collected", ctx.AddressContext, ctx.Zone)
	result <- lib.MetricResult{Name: sipStatisticsName, Success: true}
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
	Name                          string  `xml:"name"`
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
