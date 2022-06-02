package exporter

import (
	"crypto/tls"
	"net/http"
	"time"

	"sonus-metrics-exporter/config"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type Repetition uint8

const (
	repeatNone Repetition = iota
	repeatPerAddressContext
	repeatPerAddressContextIpInterfaceGroup
	repeatPerAddressContextZone
)

// SonusMetric describes a class of metric, and how to load and process its data
type SonusMetric struct {
	Processor  func(MetricContext, *[]byte, chan<- prometheus.Metric, chan<- bool)
	URLGetter  func(MetricContext) string
	APIMetrics map[string]*prometheus.Desc
	Repetition
}

// This is the list of metric classes to be collected
var metricList = []SonusMetric{
	{ProcessDSPUsage, GetDSPUrl, DSPMetrics, repeatNone},
	{ProcessFans, GetFanUrl, FanMetrics, repeatNone},
	{ProcessIPInterfaceStatus, GetIPInterfaceGroupUrl, IPInterfaceMetrics, repeatPerAddressContextIpInterfaceGroup},
	{ProcessPowerSupplies, GetPowerSupplyUrl, PowerSupplyMetrics, repeatNone},
	{ProcessSipStatistics, GetSipStatisticsUrl, SipStatisticMetrics, repeatPerAddressContextZone},
	{ProcessTGs, GetTGUrl, TGMetrics, repeatNone},
}

type MetricContext struct {
	APIBase          string
	AddressContext   string
	Zone             string
	IPInterfaceGroup string
}

// Exporter is used to store Metrics data and embeds the config struct.
// This is done so that the relevant functions have easy access to the
// user defined runtime configuration when the Collect method is called.
type Exporter struct {
	APIMetrics map[string]*prometheus.Desc
	config.Config
}

// AddMetrics - Adds all the metrics to a map of strings, returns the map.
func AddMetrics() map[string]*prometheus.Desc {

	APIMetrics := make(map[string]*prometheus.Desc)

	// Merge Zone metrics
	for k, v := range ZoneStatusMetrics {
		APIMetrics[k] = v
	}

	// Merge each class' metrics
	for _, metric := range metricList {
		for k, v := range metric.APIMetrics {
			APIMetrics[k] = v
		}
	}

	return APIMetrics
}

// Describe - loops through the API metrics and passes them to prometheus.Describe
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.APIMetrics {
		ch <- m
	}
}

// Collect function, called on by Prometheus Client library
// This function is called when a scrape is performed on the /metrics page
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	var (
		addressContexts           []*AddressContext
		results                   = make(chan bool)
		requestCount, resultCount uint
	)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 10,
	}

	// Create addressContext structs and identify zones
	for _, ac := range e.Config.APIAddressContexts {
		var (
			acs      AddressContext
			err      error
			response *Response
		)

		acs = AddressContext{ac, nil, nil}
		addressContexts = append(addressContexts, &acs)

		// TODO: Something should iterate over APIURLs to figure out which url is the active SBC
		zoneStatusUrl := acs.GetZoneStatusURL(MetricContext{e.APIURLs[0], ac, "", ""})
		response, err = doHTTPRequest(client, zoneStatusUrl, e.APIUser, e.APIPass)
		if err != nil {
			log.Errorf("Unable to perform HTTP request to %q. Error: %v", zoneStatusUrl, err)
			return
		}
		err = ProcessZones(&acs, response.body, ch)

		ipInterfaceGroupUrl := acs.GetIPInterfaceGroupURL(MetricContext{e.APIURLs[0], ac, "", ""})
		response, err = doHTTPRequest(client, ipInterfaceGroupUrl, e.APIUser, e.APIPass)
		if err != nil {
			log.Errorf("Unable to perform HTTP request to %q. Error: %v", zoneStatusUrl, err)
			return
		}
		err = ProcessIPInterfaceGroups(&acs, response.body, ch)
	}

	// Perform HTTP requests one at a time then delegate xml deserialization and metric processing to a goroutine
	go func() {
		for _, metric := range metricList {
			if metric.Repetition == repeatNone {
				var ctx = MetricContext{e.APIURLs[0], "", "", ""}
				var url = metric.URLGetter(ctx)
				requestCount++
				response, err := doHTTPRequest(client, url, e.APIUser, e.APIPass)

				if err != nil {
					log.Errorf("Unable to perform HTTP request to %q. Error: %v", url, err)
					results <- false
				} else {
					go metric.Processor(ctx, response.body, ch, results)
				}
			} else if metric.Repetition == repeatPerAddressContext {
				for _, ac := range addressContexts {
					var ctx = MetricContext{e.APIURLs[0], ac.Name, "", ""}

					var url = metric.URLGetter(ctx)
					requestCount++
					response, err := doHTTPRequest(client, url, e.APIUser, e.APIPass)

					if err != nil {
						log.Errorf("Unable to perform HTTP request to %q. Error: %v", url, err)
						results <- false
					} else {
						go metric.Processor(ctx, response.body, ch, results)
					}
				}
			} else if metric.Repetition == repeatPerAddressContextZone {
				for _, ac := range addressContexts {
					for _, zone := range ac.Zones {
						var ctx = MetricContext{e.APIURLs[0], ac.Name, zone.Name, ""}

						var url = metric.URLGetter(ctx)
						requestCount++
						response, err := doHTTPRequest(client, url, e.APIUser, e.APIPass)

						if err != nil {
							log.Errorf("Unable to perform HTTP request to %q. Error: %v", url, err)
							results <- false
						} else {
							go metric.Processor(ctx, response.body, ch, results)
						}
					}
				}
			} else if metric.Repetition == repeatPerAddressContextIpInterfaceGroup {
				for _, ac := range addressContexts {
					for _, ipig := range ac.IPInterfaceGroups {
						var ctx = MetricContext{e.APIURLs[0], ac.Name, "", ipig.Name}

						var url = metric.URLGetter(ctx)
						requestCount++
						response, err := doHTTPRequest(client, url, e.APIUser, e.APIPass)

						if err != nil {
							log.Errorf("Unable to perform HTTP request to %q. Error: %v", url, err)
							results <- false
						} else {
							go metric.Processor(ctx, response.body, ch, results)
						}
					}
				}
			}
		}
	}()

	for {
		select {
		case <-results:
			// TODO: Collect metrics about success/failures
			resultCount++
			if resultCount == requestCount {
				log.Info("Done collectin'")
				transport.CloseIdleConnections()
				return
			}
		}
	}

}
