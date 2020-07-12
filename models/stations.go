package models

import "net/http"

type Station struct {
	Name string
	CollectFunc func(*http.Client) ([]Price, error)
}

var stations = make([]Station, 0)

func RegisterStation(station Station) {
	stations = append(stations, station)
}

func GetStations() []Station {
	return stations
}
