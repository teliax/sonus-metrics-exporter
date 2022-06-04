package exporter

import (
	"crypto/tls"
	"net/http"
	"time"

	"sonus-metrics-exporter/config"
	"sonus-metrics-exporter/lib"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Exporter is used to store Metrics data and embeds the config struct.
// This is done so that the relevant functions have easy access to the
// user defined runtime configuration when the Collect method is called.
type Exporter struct {
	Metrics []lib.SonusMetric
	config.Config
}

// Describe - loops through the API metrics and passes them to prometheus.Describe
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, am := range zoneStatusMetrics {
		ch <- am
	}

	for _, m := range e.Metrics {
		for _, am := range m.APIMetrics {
			ch <- am
		}
	}
}

// Collect function, called on by Prometheus Client library
// This function is called when a scrape is performed on the /metrics page
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	var (
		addressContexts           []*AddressContext
		results                   = make(chan lib.MetricResult)
		requestCount, resultCount uint
	)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 10,
	}

	// Create addressContext structs, and identify zones and ipInterfaceGroups
	for _, acName := range e.Config.APIAddressContexts {
		var (
			ac       = AddressContext{Name: acName}
			err      error
			response *HTTPResponse
		)

		addressContexts = append(addressContexts, &ac)

		// TODO: Something should iterate over APIURLs to figure out which url is the active SBC
		zoneStatusUrl := ac.getZoneStatusURL(lib.MetricContext{APIBase: e.APIURLs[0], AddressContext: ac.Name})
		response, err = doHTTPRequest(client, zoneStatusUrl, e.APIUser, e.APIPass)
		if err != nil {
			log.Errorf("Unable to perform HTTP request to %q. Error: %v", zoneStatusUrl, err)
			return
		}
		err = processZones(&ac, response.body, ch)

		ipInterfaceGroupUrl := ac.getIPInterfaceGroupURL(lib.MetricContext{APIBase: e.APIURLs[0], AddressContext: ac.Name})
		response, err = doHTTPRequest(client, ipInterfaceGroupUrl, e.APIUser, e.APIPass)
		if err != nil {
			log.Errorf("Unable to perform HTTP request to %q. Error: %v", zoneStatusUrl, err)
			return
		}
		err = processIPInterfaceGroups(&ac, response.body)
	}

	// Perform HTTP requests one at a time then delegate xml deserialization and metric processing to a goroutine
	go func() {
		for _, metric := range e.Metrics {
			if metric.Repetition == lib.RepeatNone {
				var (
					ctx = lib.MetricContext{APIBase: e.APIURLs[0]}
					url = metric.URLGetter(ctx)
				)

				requestCount++
				response, err := doHTTPRequest(client, url, e.APIUser, e.APIPass)

				if err != nil {
					log.Errorf("Unable to perform HTTP request to %q. Error: %v", url, err)
					results <- lib.MetricResult{Success: false, Errors: []*error{&err}}
				} else {
					go metric.Processor(ctx, response.body, ch, results)
				}
			} else if metric.Repetition == lib.RepeatPerAddressContext {
				for _, ac := range addressContexts {
					var (
						ctx = lib.MetricContext{APIBase: e.APIURLs[0], AddressContext: ac.Name}
						url = metric.URLGetter(ctx)
					)

					requestCount++
					response, err := doHTTPRequest(client, url, e.APIUser, e.APIPass)

					if err != nil {
						log.Errorf("Unable to perform HTTP request to %q. Error: %v", url, err)
						results <- lib.MetricResult{Success: false, Errors: []*error{&err}}
					} else {
						go metric.Processor(ctx, response.body, ch, results)
					}
				}
			} else if metric.Repetition == lib.RepeatPerAddressContextZone {
				for _, ac := range addressContexts {
					for _, zone := range ac.Zones {
						var (
							ctx = lib.MetricContext{APIBase: e.APIURLs[0], AddressContext: ac.Name, Zone: zone.Name}
							url = metric.URLGetter(ctx)
						)

						requestCount++
						response, err := doHTTPRequest(client, url, e.APIUser, e.APIPass)

						if err != nil {
							log.Errorf("Unable to perform HTTP request to %q. Error: %v", url, err)
							results <- lib.MetricResult{Success: false, Errors: []*error{&err}}
						} else {
							go metric.Processor(ctx, response.body, ch, results)
						}
					}
				}
			} else if metric.Repetition == lib.RepeatPerAddressContextIpInterfaceGroup {
				for _, ac := range addressContexts {
					for _, ipig := range ac.IPInterfaceGroups {
						var (
							ctx = lib.MetricContext{APIBase: e.APIURLs[0], AddressContext: ac.Name, IPInterfaceGroup: ipig.Name}
							url = metric.URLGetter(ctx)
						)

						requestCount++
						response, err := doHTTPRequest(client, url, e.APIUser, e.APIPass)

						if err != nil {
							log.Errorf("Unable to perform HTTP request to %q. Error: %v", url, err)
							results <- lib.MetricResult{Success: false, Errors: []*error{&err}}
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
