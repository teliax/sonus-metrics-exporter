package lib

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	// Repetition is used to indicate to the exporter how a metric should be repeated
	Repetition uint8

	// SonusMetric describes a class of metric, and how to load and process its data
	SonusMetric struct {
		Name       string
		Processor  func(MetricContext, *[]byte)
		URLGetter  func(MetricContext) string
		APIMetrics map[string]*prometheus.Desc
		Repetition
	}

	// MetricContext contains information and channels used by metric collectors to do their thing
	MetricContext struct {
		APIBase          string
		AddressContext   string
		Zone             string
		IPInterfaceGroup string
		MetricChannel    chan<- prometheus.Metric
		ResultChannel    chan<- MetricResult
	}

	// MetricResult indicates success or failure+errors from a metric collector
	MetricResult struct {
		Name    string
		Success bool
		Errors  []*error
	}
)

const (
	RepeatNone Repetition = iota
	RepeatPerAddressContext
	RepeatPerAddressContextIpInterfaceGroup
	RepeatPerAddressContextZone
)
