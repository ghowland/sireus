package exporter

import (
	"fmt"
	"github.com/ghowland/sireus/code/data"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// Runs the Prometheus Exporter listener in the background
func RunExporterListener() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", data.SireusData.AppConfig.PrometheusExportPort), nil)
}
