# fuel-station-exporter
A [prometheus](https://prometheus.io) exporter for getting some metrics of some fuel stations in Eindhoven, The Netherlands. Currently supports the Esso, Tango and Tinq locations in the city.

## Installation
If you have a working Go installation, getting the binary should be as simple as

```
go get github.com/wouter0100/fuel-station-exporter
```

## Usage
```plain
$ fuel-station-exporter
```

The following environment variables are required:
```
FUEL_STATION_LISTEN_ADDR=localhost:9684 
FUEL_STATION_TIMEOUT=10
```

## Metrics
The following metrics are currently returned:
```
# HELP current_fuel_price The current fuel price for different fuel types.
# TYPE current_fuel_price gauge
current_fuel_price{station="esso-express-eindhoven",type="Diesel"} 1.189
current_fuel_price{station="esso-express-eindhoven",type="Euro95 E5"} 1.499
current_fuel_price{station="tango-eindhoven-ruysdaelbaan",type="Diesel"} 1.209
current_fuel_price{station="tango-eindhoven-ruysdaelbaan",type="Euro95 E5"} 1.499
current_fuel_price{station="tinq-eindhoven-hurksestraat",type="Diesel"} 1.149
current_fuel_price{station="tinq-eindhoven-hurksestraat",type="Euro95 E10"} 1.479
current_fuel_price{station="tinq-eindhoven-hurksestraat",type="Euro95 E5"} 1.499
```
