package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"go-fuel-station-exporter/models"
	"log"
	"net/http"
	"time"

	_ "go-fuel-station-exporter/stations"
)

type fuelCollector struct {
	client             *http.Client
	userAgent          string
	currentPriceMetric *prometheus.GaugeVec
}

func (f fuelCollector) Describe(desc chan<- *prometheus.Desc) {
	f.currentPriceMetric.Describe(desc)
}

func (f fuelCollector) Collect(ch chan<- prometheus.Metric) {
	for _, station := range models.GetStations() {
		results, err := station.CollectFunc(f.client)
		if err != nil {
			log.Fatal(err)
		}

		for _, result := range results {
			f.currentPriceMetric.With(prometheus.Labels{
				"type": result.FuelType.String(),
				"station": station.Name,
			}).Set(result.Amount)
		}
	}

	f.currentPriceMetric.Collect(ch)
}

func newCollector(timeout time.Duration) *fuelCollector {
	return &fuelCollector{
		client: &http.Client{
			Timeout: timeout,
		},
		currentPriceMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "current_fuel_price",
			Help: "The current fuel price for different fuel types.",
		}, []string{"type", "station"}),
	}
}
