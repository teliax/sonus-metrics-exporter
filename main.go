package main

import (
	"net/http"

	"sonus-metrics-exporter/config"
	"sonus-metrics-exporter/exporter"
	"sonus-metrics-exporter/lib"
	"sonus-metrics-exporter/metrics"

	"github.com/fatih/structs"
	"github.com/infinityworks/go-common/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	log            *logrus.Logger
	applicationCfg config.Config

	metricList = []lib.SonusMetric{
		metrics.DSPMetric,
		metrics.FanMetric,
		metrics.IPInterfaceMetric,
		metrics.PowerSupplyMetric,
		metrics.SipStatisticMetric,
		metrics.TGMetric,
	}
)

func init() {
	applicationCfg = config.Init()
	log = logger.Start(&applicationCfg)
}

func main() {

	log.WithFields(structs.Map(applicationCfg)).Info("Starting Exporter")

	ex := exporter.Exporter{
		Metrics: metricList,
		Config:  applicationCfg,
	}

	// Register Metrics from each of the endpoints
	// This invokes the Collect method through the prometheus client libraries.
	prometheus.MustRegister(&ex)

	// Setup HTTP handler
	http.Handle(applicationCfg.MetricsPath(), promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		                <head><title>Sonus Exporter</title></head>
		                <body>
		                   <h1>Sonus Metrics Exporter</h1>
						   <p>For more information, visit <a href=https://github.com/teliax/sonus-metrics-exporter>GitHub</a></p>
		                   <p><a href='` + applicationCfg.MetricsPath() + `'>Metrics</a></p>
		                   </body>
		                </html>
		              `))
	})
	log.Fatal(http.ListenAndServe(":"+applicationCfg.ListenPort(), nil))
}
