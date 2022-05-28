package exporter

import "github.com/prometheus/client_golang/prometheus"

// AddMetrics - Adds all of the metrics to a map of strings, returns the map.
func AddMetrics() map[string]*prometheus.Desc {

	APIMetrics := make(map[string]*prometheus.Desc)

	APIMetrics["TG_Usage"] = prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "TG", "usage_total"),
		"Number of active calls",
		[]string{"zone", "name", "direction"}, nil,
	)
	APIMetrics["TG_Bandwidth"] = prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "TG", "bytes"),
		"Bandwidth in use by current calls",
		[]string{"zone", "name", "direction"}, nil,
	)
	APIMetrics["TG_State"] = prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "TG", "state"),
		"State of the trunkgroup",
		[]string{"zone", "name"}, nil,
	)
	APIMetrics["TG_OBState"] = prometheus.NewDesc(
		prometheus.BuildFQName("sonus", "TG", "outbound_state"),
		"State of outbound calls on the trunkgroup",
		[]string{"zone", "name"}, nil,
	)

	return APIMetrics
}

// processMetrics - processes the response data and sets the metrics using it as a source
func (e *Exporter) processMetrics(data []*TGcollection, ch chan<- prometheus.Metric) error {

	// APIMetrics - range through the data slice
	for _, y := range data {
		for _, x := range y.TGglobalTrunkGroupStatus {

			state := x.TGstate.Value == "inService"
			outState := x.TGpacketOutDetectState.Value == "normal"
			var stateMetric float64
			var outStateMetric float64
			if state {
				stateMetric = 1
			}
			if outState {
				outStateMetric = 1
			}
			ch <- prometheus.MustNewConstMetric(e.APIMetrics["Usage"], prometheus.GaugeValue, x.TGinboundCallsUsage.Number, x.TGzone.Value, x.TGname.Value, "inbound")
			ch <- prometheus.MustNewConstMetric(e.APIMetrics["Usage"], prometheus.GaugeValue, x.TGoutboundCallsUsage.Number, x.TGzone.Value, x.TGname.Value, "outbound")
			ch <- prometheus.MustNewConstMetric(e.APIMetrics["Bandwidth"], prometheus.GaugeValue, x.TGbwInboundUsage.Number, x.TGzone.Value, x.TGname.Value, "inbound")
			ch <- prometheus.MustNewConstMetric(e.APIMetrics["Bandwidth"], prometheus.GaugeValue, x.TGbwOutboundUsage.Number, x.TGzone.Value, x.TGname.Value, "outbound")
			ch <- prometheus.MustNewConstMetric(e.APIMetrics["State"], prometheus.GaugeValue, stateMetric, x.TGzone.Value, x.TGname.Value)
			ch <- prometheus.MustNewConstMetric(e.APIMetrics["OBState"], prometheus.GaugeValue, outStateMetric, x.TGzone.Value, x.TGname.Value)
		}
	}

	return nil
}
