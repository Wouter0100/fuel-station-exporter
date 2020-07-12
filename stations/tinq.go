package stations

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-fuel-station-exporter/models"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	models.RegisterStation(models.Station{
		Name:        "tinq-eindhoven-hurksestraat",
		CollectFunc: collectTinq,
	})
}

func collectTinq(client *http.Client) ([]models.Price, error) {
	var results = make([]models.Price, 0)

	res, err := client.Get("https://www.tinq.nl/tankstations/eindhoven-hurksestraat")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var loopErr error

	doc.Find(".node--type-price").EachWithBreak(func(i int, s *goquery.Selection) bool {
		fuelTypeStringSelection := s.Find(".field--name-taxonomy-term-title").First()
		if fuelTypeStringSelection.Nodes == nil {
			loopErr = fmt.Errorf("no fuel type found in Tinq")
			return false
		}

		fuelTypeString := strings.TrimSpace(fuelTypeStringSelection.Text())

		var fuelType models.FuelType

		if fuelTypeString == "Euro95 E10" {
			fuelType = models.Euro95E10
		} else if fuelTypeString == "Euro95 E5" {
			fuelType = models.Euro95E5
		} else if fuelTypeString == "Diesel" {
			fuelType = models.Diesel
		} else {
			loopErr = fmt.Errorf("unable to convert `%+v` to a FuelType in Tinq", fuelTypeString)
			return false
		}

		priceSelection := s.Find(".field--name-field-prices-price-pump").First()
		if priceSelection.Nodes == nil {
			loopErr = fmt.Errorf("no pricings found for `%+v` for Tinq", fuelTypeString)
			return false
		}

		priceString, exists := priceSelection.Attr("content")
		if !exists {
			loopErr = fmt.Errorf("unable to get content for `%+v` price for Tinq", fuelTypeString)
			return false
		}

		price, err := strconv.ParseFloat(priceString, 64)
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
