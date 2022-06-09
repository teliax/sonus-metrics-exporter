package exporter

import (
	"crypto/tls"
	"net/http"
	"strconv"

	"sonus-metrics-exporter/config"
	"sonus-metrics-exporter/lib"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// MetricDisposition is a counter metric that tracks success per each metric type
var MetricDisposition = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "sonus",
	Subsystem: "exporter",
	Name:      "metric_disposition",
	Help:      "Number of times each metric has succeeded or failed being collected",
}, []string{"name", "successful"})

func init() {
	prometheus.MustRegister(MetricDisposition)
}

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
		apiBase                   string
		results                   = make(chan lib.MetricResult)
		requestCount, resultCount uint
		httpTransport             = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		httpClient                = &http.Client{
			Transport: httpTransport,
			Timeout:   e.APITimeout,
		}
	)

	for i, url := range e.APIURLs {
		serverStatusUrl := getServerStatusURL(url)
		response, err := doHTTPRequest(httpClient, serverStatusUrl, e.APIUser, e.APIPass)

		if err != nil {
			log.Errorf("Error encountered attemping to validate API_URL %q, index %d. Error: %v", url, i, err)
			continue
		}
		if response.response.StatusCode == 200 {
			apiBase = url
			log.Infof("Using API_URL %q.", apiBase)
			break
		} else {
			log.Errorf("Non-200 HTTP reponse (%d) validating API_URL %q.", response.response.StatusCode, url)
		}

	}

	if len(apiBase) == 0 {
		log.Error("Unable to find an active SBC in API_URLS")
		return
	}

	// Create addressContext structs, and identify zones and ipInterfaceGroups
	for _, acName := range e.Config.APIAddressContexts {
		var (
			ac       = AddressContext{Name: acName}
			err      error
			response *HTTPResponse
		)

		addressContexts = append(addressContexts, &ac)

		zoneStatusUrl := ac.getZoneStatusURL(lib.MetricContext{APIBase: apiBase, AddressContext: ac.Name})
		response, err = doHTTPRequest(httpClient, zoneStatusUrl, e.APIUser, e.APIPass)
		if err != nil {
			log.Errorf("Unable to perform HTTP request to %q. Error: %v", zoneStatusUrl, err)
			return
		}
		if response.response.StatusCode != 200 {
			log.Errorf("Non-200 HTTP reponse (%d) to %q.", response.response.StatusCode, zoneStatusUrl)
			return
		}
		err = processZones(&ac, response.body, ch)

		ipInterfaceGroupUrl := ac.getIPInterfaceGroupURL(lib.MetricContext{APIBase: apiBase, AddressContext: ac.Name})
		response, err = doHTTPRequest(httpClient, ipInterfaceGroupUrl, e.APIUser, e.APIPass)
		if err != nil {
			log.Errorf("Unable to perform HTTP request to %q. Error: %v", ipInterfaceGroupUrl, err)
			return
		}
		if response.response.StatusCode != 200 {
			log.Errorf("Non-200 HTTP reponse (%d) to %q.", response.response.StatusCode, ipInterfaceGroupUrl)
			return
		}
		err = processIPInterfaceGroups(&ac, response.body)
	}

	// Perform HTTP requests one at a time then delegate xml deserialization and metric processing to a goroutine
	go func() {
		for _, metric := range e.Metrics {
			if metric.Repetition == lib.RepeatNone {
				var (
					ctx = lib.MetricContext{APIBase: apiBase}
					url = metric.URLGetter(ctx)
				)

				requestCount++
				response, err := doHTTPRequest(httpClient, url, e.APIUser, e.APIPass)

				if err != nil {
					log.Errorf("Unable to perform HTTP request to %q. Error: %v", url, err)
					results <- lib.MetricResult{Success: false, Errors: []*error{&err}}
				} else {
					go metric.Processor(ctx, response.body, ch, results)
				}
			} else if metric.Repetition == lib.RepeatPerAddressContext {
				for _, ac := range addressContexts {
					var (
						ctx = lib.MetricContext{APIBase: apiBase, AddressContext: ac.Name}
						url = metric.URLGetter(ctx)
					)

					requestCount++
					response, err := doHTTPRequest(httpClient, url, e.APIUser, e.APIPass)

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
							ctx = lib.MetricContext{APIBase: apiBase, AddressContext: ac.Name, Zone: zone.Name}
							url = metric.URLGetter(ctx)
						)

						requestCount++
						response, err := doHTTPRequest(httpClient, url, e.APIUser, e.APIPass)

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
							ctx = lib.MetricContext{APIBase: apiBase, AddressContext: ac.Name, IPInterfaceGroup: ipig.Name}
							url = metric.URLGetter(ctx)
						)

						requestCount++
						response, err := doHTTPRequest(httpClient, url, e.APIUser, e.APIPass)

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
		case result := <-results:
			var successString = strconv.FormatBool(result.Success)
			MetricDisposition.WithLabelValues(result.Name, successString).Inc()

			resultCount++
			if resultCount == requestCount {
				log.Info("Done collectin'")
				httpTransport.CloseIdleConnections()
				return
			}
		}
	}

}
