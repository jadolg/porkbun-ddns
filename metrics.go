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
	updateErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "porkbun_update_errors_total",
			Help: "The total number of errors found",
		},
		[]string{"host", "domain"},
	)
	credentialsErrorTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "porkbun_credentials_error_total",
		Help: "The total number of credentials errors",
	})
)

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
