package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	resolveErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "porkbun_resolve_errors_total",
			Help: "The total number of errors trying to resolve the current ip address",
		},
		[]string{"record_type"},
	)
	updateErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "porkbun_update_errors_total",
			Help: "The total number of errors found",
		},
		[]string{"host", "domain"},
	)
	credentialsErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "porkbun_credentials_errors_total",
		Help: "The total number of credentials errors",
	})
	connectionErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "porkbun_connection_errors_total",
		Help: "The total number of connection errors (while connecting to porkbun)",
	})
	updateSuccessTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "porkbun_update_success_total",
			Help: "The total number of successful updates",
		},
		[]string{"host", "domain", "record_type"},
	)
)

// initMetrics initializes the metrics with the records provided in the configuration
// We want to have 0 as a baseline for all possible labels instead of non-existing metrics
func initMetrics(records []Record) {
	for _, record := range records {
		updateSuccessTotal.WithLabelValues(record.Host, record.Domain, "A").Add(0)
		updateSuccessTotal.WithLabelValues(record.Host, record.Domain, "AAAA").Add(0)
		updateErrorsTotal.WithLabelValues(record.Host, record.Domain).Add(0)
		resolveErrorsTotal.WithLabelValues("A").Add(0)
		resolveErrorsTotal.WithLabelValues("AAAA").Add(0)
	}
}

func healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "OK")
	if err != nil {
		log.Errorf("error responding to request %v", err)
	}
}

func startMetricsServer(port int) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", healthcheckHandler)
	log.Printf("Starting metrics server on port %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
