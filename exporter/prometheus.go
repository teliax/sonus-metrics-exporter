package exporter

import (
	"crypto/tls"
	"net/http"
	"strconv"
	"time"

	"sonus-metrics-exporter/config"
	"sonus-metrics-exporter/lib"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	// metricDisposition is a counter metric that tracks success per each metric type
	metricDisposition = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "sonus",
		Subsystem: "exporter",
		Name:      "metric_disposition",
		Help:      "Number of times each metric has succeeded or failed being collected",
	}, []string{"name", "successful"})

	//metricDuration is a summary that tracks how long metrics api requests and processing take
	metricDuration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  "sonus",
		Subsystem:  "exporter",
		Name:       "metric_duration",
		Help:       "How long metrics took to query and process",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"name", "stage"})
)

func init() {
	prometheus.MustRegister(metricDisposition)
	prometheus.MustRegister(metricDuration)
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
		addressContexts           []*addressContext
		apiBase                   string
		collectCount, resultCount uint
		httpTransport             = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		httpClient                = &http.Client{Transport: httpTransport, Timeout: e.APITimeout}
		results                   = make(chan lib.MetricResult)
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
			ac       = addressContext{Name: acName}
			err      error
			response *httpResponse
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
			ctx := lib.MetricContext{APIBase: apiBase, MetricChannel: ch, ResultChannel: results}
			if metric.Repetition == lib.RepeatNone {
				collectCount++
				doHTTPAndProcess(e, metric, ctx, httpClient)
			} else if metric.Repetition == lib.RepeatPerAddressContext {
				for _, ac := range addressContexts {
					ctx.AddressContext = ac.Name

					collectCount++
					doHTTPAndProcess(e, metric, ctx, httpClient)
				}
			} else if metric.Repetition == lib.RepeatPerAddressContextZone {
				for _, ac := range addressContexts {
					for _, zone := range ac.Zones {
						ctx.AddressContext = ac.Name
						ctx.Zone = zone.Name

						collectCount++
						doHTTPAndProcess(e, metric, ctx, httpClient)
					}
				}
			} else if metric.Repetition == lib.RepeatPerAddressContextIpInterfaceGroup {
				for _, ac := range addressContexts {
					for _, ipig := range ac.IPInterfaceGroups {
						ctx.AddressContext = ac.Name
						ctx.IPInterfaceGroup = ipig.Name

						collectCount++
						doHTTPAndProcess(e, metric, ctx, httpClient)
					}
				}
			}
		}
	}()

	for {
		select {
		case result := <-results:
			var successString = strconv.FormatBool(result.Success)
			metricDisposition.WithLabelValues(result.Name, successString).Inc()

			resultCount++
			if resultCount == collectCount {
				log.Info("Done collectin'")
				httpTransport.CloseIdleConnections()
				return
			}
		}
	}

}

func doHTTPAndProcess(e *Exporter, metric lib.SonusMetric, ctx lib.MetricContext, httpClient *http.Client) {
	url := metric.URLGetter(ctx)
	ht := time.Now()
	response, err := doHTTPRequest(httpClient, url, e.APIUser, e.APIPass)
	metricDuration.WithLabelValues(metric.Name, "http").Observe(time.Since(ht).Seconds())

	if err != nil {
		log.Errorf("Unable to perform HTTP request to %q. Error: %v", url, err)
		ctx.ResultChannel <- lib.MetricResult{Success: false, Errors: []*error{&err}}
	} else {
		go func(m lib.SonusMetric, c lib.MetricContext, r *httpResponse) {
			mt := time.Now()
			m.Processor(c, r.body)
			metricDuration.WithLabelValues(m.Name, "process").Observe(time.Since(mt).Seconds())
		}(metric, ctx, response)
	}
}
