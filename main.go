package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	listenAddr := os.Getenv("FUEL_STATION_LISTEN_ADDR")
	timeout, _ := time.ParseDuration(os.Getenv("FUEL_STATION_TIMEOUT"))

	collector := newCollector(timeout)
	if err := prometheus.Register(collector); err != nil {
		log.Fatalf("Failed to register collector: %s", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", http.RedirectHandler("/metrics", http.StatusFound))

	log.Printf("Listen on %s...", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
