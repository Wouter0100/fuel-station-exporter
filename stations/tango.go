package stations

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-fuel-station-exporter/models"
	"net/http"
	"strconv"
)

func init() {
	models.RegisterStation(models.Station{
		Name:        "tango-eindhoven-ruysdaelbaan",
		CollectFunc: collectTango,
	})
}

func collectTango(client *http.Client) ([]models.Price, error) {
	var results = make([]models.Price, 0)

	res, err := client.Get("https://www.tango.nl/stations/tango-eindhoven-ruysdaelbaan/")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var loopErr error

	doc.Find(".pricing").EachWithBreak(func(i int, s *goquery.Selection) bool {
		fuelTypeString, exists := s.Attr("id")
		if !exists {
			loopErr = fmt.Errorf("unable to get FuelTypeString for Tango")
			return false
		}

		var fuelType models.FuelType

		if fuelTypeString == "euro95" {
			fuelType = models.Euro95E5
		} else if fuelTypeString == "diesel" {
			fuelType = models.Diesel
		} else {
			loopErr = fmt.Errorf("unable to convert `%+v` to a FuelType in Tango", fuelTypeString)
			return false
		}

		priceSelection := s.Find(".pump_price .price").First()

		if priceSelection.Nodes == nil {
			loopErr = fmt.Errorf("no pricing found for `%+v` in Tango", fuelTypeString)
			return false
		}

		price, err := strconv.ParseFloat(priceSelection.Nodes[0].FirstChild.Data, 64)
		if err != nil {
			loopErr = err
			return false
		}

		results = append(results, models.Price{
			Amount:   price,
			FuelType: fuelType,
		})

		return true
	})

	if loopErr != nil {
		return nil, loopErr
	}

	return results, nil
}
