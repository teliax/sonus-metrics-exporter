package metrics

import (
	"encoding/xml"

	"sonus-metrics-exporter/lib"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var DSPMetric = lib.SonusMetric{
	Name:       "DSP",
	Processor:  processDSPUsage,
	URLGetter:  getDSPUrl,
	APIMetrics: dspMetrics,
	Repetition: lib.RepeatNone,
}

const dspUrlSuffix = "/operational/system/dspStatus/dspUsage/"

func getDSPUrl(ctx lib.MetricContext) string {
	return ctx.APIBase + dspUrlSuffix
}

var dspMetrics = map[string]*prometheus.Desc{
	"DSP_Resources_Used": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "dsp", "resources_used"),
		"Usage of DSP resources per slot",
		[]string{"system", "slot"}, nil,
	),
	"DSP_Resources_Total": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "dsp", "resources_total"),
		"Total compression resources",
		[]string{"system"}, nil,
	),
	"DSP_Compression_Utilization": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "dsp", "compression_utilization"),
		"Compression resource utilization, in percent",
		[]string{"system"}, nil,
	),
	"DSP_Codec_Utilization": prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "dsp", "codec_utilization"),
		"Codec utilization, in percent",
		[]string{"system", "codec"}, nil,
	),
}

func processDSPUsage(ctx lib.MetricContext, xmlBody *[]byte, ch chan<- prometheus.Metric, result chan<- lib.MetricResult) {
	dsp := new(dspUsageCollection)
	err := xml.Unmarshal(*xmlBody, &dsp)
	if err != nil {
		log.Errorf("Failed to deserialize dspStatus XML: %v", err)
		result <- lib.MetricResult{Success: false, Errors: []*error{&err}}
		return
	}

	var d = dsp.DSPUsage

	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Resources_Used"], prometheus.GaugeValue, d.Slot1ResourcesUtilized, d.SystemName, "1")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Resources_Used"], prometheus.GaugeValue, d.Slot2ResourcesUtilized, d.SystemName, "2")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Resources_Used"], prometheus.GaugeValue, d.Slot3ResourcesUtilized, d.SystemName, "3")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Resources_Used"], prometheus.GaugeValue, d.Slot4ResourcesUtilized, d.SystemName, "4")

	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Resources_Total"], prometheus.GaugeValue, d.CompressionTotal, d.SystemName)

	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Compression_Utilization"], prometheus.GaugeValue, d.CompressionUtilization, d.SystemName)

	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.G711Utilization, d.SystemName, "G.711")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.G711SsUtilization, d.SystemName, "G.711 Silence Suppression")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.G726Utilization, d.SystemName, "G.726")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.G7231Utilization, d.SystemName, "G.723.1")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.G722Utilization, d.SystemName, "G.722")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.G7221Utilization, d.SystemName, "G.722.1")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.G729AbUtilization, d.SystemName, "G.729")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.EcmUtilization, d.SystemName, "ECM")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.IlbcUtilization, d.SystemName, "iLBC")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.AmrNbUtilization, d.SystemName, "AMR-NB")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.AmrWbUtilization, d.SystemName, "AMR-WB")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.ToneUtilization, d.SystemName, "Tone")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.G711V8Utilization, d.SystemName, "G.711 V8")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.G711SsV8Utilization, d.SystemName, "G.711 Silence Suppression V8")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.G726V8Utilization, d.SystemName, "G.726 V8")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.G7231V8Utilization, d.SystemName, "G.723.1 V8")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.G722V8Utilization, d.SystemName, "G.722 V8")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.G7221V8Utilization, d.SystemName, "G.722.1 V8")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.G729AbV8Utilization, d.SystemName, "G.729 V8")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.EcmV34Utilization, d.SystemName, "ECM V.34")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.IlbcV8Utilization, d.SystemName, "iLBC V8")
	ch <- prometheus.MustNewConstMetric(dspMetrics["DSP_Codec_Utilization"], prometheus.GaugeValue, d.OpusUtilization, d.SystemName, "Opus")

	log.Info("DSP Metrics collected")
	result <- lib.MetricResult{Success: true}
}

/*
<collection xmlns:y="http://tail-f.com/ns/rest">
  <dspUsage xmlns="http://sonusnet.com/ns/mibs/SONUS-DRM-DSPSTATUS/1.0">
    <systemName>densbc01</systemName>
    <slot1ResourcesUtilized>68</slot1ResourcesUtilized>
    <slot2ResourcesUtilized>0</slot2ResourcesUtilized>
    ...
  </dspUsage>
</collection>
*/

type dspUsageCollection struct {
	DSPUsage *dspUsage `xml:"http://sonusnet.com/ns/mibs/SONUS-DRM-DSPSTATUS/1.0 dspUsage"`
}

type dspUsage struct {
	SystemName                         string  `xml:"systemName"`
	Slot1ResourcesUtilized             float64 `xml:"slot1ResourcesUtilized"`
	Slot2ResourcesUtilized             float64 `xml:"slot2ResourcesUtilized"`
	Slot3ResourcesUtilized             float64 `xml:"slot3ResourcesUtilized"`
	Slot4ResourcesUtilized             float64 `xml:"slot4ResourcesUtilized"`
	CompressionTotal                   float64 `xml:"compressionTotal"`
	CompressionAvailable               float64 `xml:"compressionAvailable"`
	CompressionUtilization             float64 `xml:"compressionUtilization"`
	CompressionHighPriorityUtilization float64 `xml:"compressionHighPriorityUtilization"`
	CompressionAllocFailures           float64 `xml:"compressionAllocFailures"`
	G711Total                          float64 `xml:"g711Total"`
	G711Utilization                    float64 `xml:"g711Utilization"`
	G711SsTotal                        float64 `xml:"g711SsTotal"`
	G711SsUtilization                  float64 `xml:"g711SsUtilization"`
	G726Total                          float64 `xml:"g726Total"`
	G726Utilization                    float64 `xml:"g726Utilization"`
	G7231Total                         float64 `xml:"g7231Total"`
	G7231Utilization                   float64 `xml:"g7231Utilization"`
	G722Total                          float64 `xml:"g722Total"`
	G722Utilization                    float64 `xml:"g722Utilization"`
	G7221Total                         float64 `xml:"g7221Total"`
	G7221Utilization                   float64 `xml:"g7221Utilization"`
	G729AbTotal                        float64 `xml:"g729AbTotal"`
	G729AbUtilization                  float64 `xml:"g729AbUtilization"`
	EcmTotal                           float64 `xml:"ecmTotal"`
	EcmUtilization                     float64 `xml:"ecmUtilization"`
	IlbcTotal                          float64 `xml:"ilbcTotal"`
	IlbcUtilization                    float64 `xml:"ilbcUtilization"`
	AmrNbTotal                         float64 `xml:"amrNbTotal"`
	AmrNbUtilization                   float64 `xml:"amrNbUtilization"`
	AmrNbT140Total                     float64 `xml:"amrNbT140Total"`
	AmrNbT140Utilization               float64 `xml:"amrNbT140Utilization"`
	AmrWbTotal                         float64 `xml:"amrWbTotal"`
	AmrWbUtilization                   float64 `xml:"amrWbUtilization"`
	AmrWbT140Total                     float64 `xml:"amrWbT140Total"`
	AmrWbT140Utilization               float64 `xml:"amrWbT140Utilization"`
	Evrcb0Total                        float64 `xml:"evrcb0Total"`
	Evrcb0Utilization                  float64 `xml:"evrcb0Utilization"`
	Evrc0Total                         float64 `xml:"evrc0Total"`
	Evrc0Utilization                   float64 `xml:"evrc0Utilization"`
	ToneTotal                          float64 `xml:"toneTotal"`
	ToneAvailable                      float64 `xml:"toneAvailable"`
	ToneUtilization                    float64 `xml:"toneUtilization"`
	ToneHighPriorityUtilization        float64 `xml:"toneHighPriorityUtilization"`
	ToneAllocFailures                  float64 `xml:"toneAllocFailures"`
	EfrTotal                           float64 `xml:"efrTotal"`
	EfrUtilization                     float64 `xml:"efrUtilization"`
	G711V8Total                        float64 `xml:"g711V8Total"`
	G711V8Utilization                  float64 `xml:"g711V8Utilization"`
	G711SsV8Total                      float64 `xml:"g711SsV8Total"`
	G711SsV8Utilization                float64 `xml:"g711SsV8Utilization"`
	G726V8Total                        float64 `xml:"g726V8Total"`
	G726V8Utilization                  float64 `xml:"g726V8Utilization"`
	G7231V8Total                       float64 `xml:"g7231V8Total"`
	G7231V8Utilization                 float64 `xml:"g7231V8Utilization"`
	G722V8Total                        float64 `xml:"g722V8Total"`
	G722V8Utilization                  float64 `xml:"g722V8Utilization"`
	G7221V8Total                       float64 `xml:"g7221V8Total"`
	G7221V8Utilization                 float64 `xml:"g7221V8Utilization"`
	G729AbV8Total                      float64 `xml:"g729AbV8Total"`
	G729AbV8Utilization                float64 `xml:"g729AbV8Utilization"`
	EcmV34Total                        float64 `xml:"ecmV34Total"`
	EcmV34Utilization                  float64 `xml:"ecmV34Utilization"`
	IlbcV8Total                        float64 `xml:"ilbcV8Total"`
	IlbcV8Utilization                  float64 `xml:"ilbcV8Utilization"`
	OpusTotal                          float64 `xml:"opusTotal"`
	OpusUtilization                    float64 `xml:"opusUtilization"`
	EvsTotal                           float64 `xml:"evsTotal"`
	EvsUtilization                     float64 `xml:"evsUtilization"`
	Silk8Total                         float64 `xml:"silk8Total"`
	Silk8Utilization                   float64 `xml:"silk8Utilization"`
	Silk16Total                        float64 `xml:"silk16Total"`
	Silk16Utilization                  float64 `xml:"silk16Utilization"`
}
