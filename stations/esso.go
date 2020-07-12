package stations

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"fuel-station-exporter/models"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	models.RegisterStation(models.Station{
		Name:        "esso-express-eindhoven",
		CollectFunc: collectEsso,
	})
}

func collectEsso(client *http.Client) ([]models.Price, error) {
	var results = make([]models.Price, 0)

	res, err := client.Get("https://www.tankstation.nl/tankstation/esso-express-eindhoven/")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var loopErr error

	doc.Find(".gasoline-prices span").EachWithBreak(func(i int, s *goquery.Selection) bool {
		var fuelType models.FuelType

		fuelTypeString := strings.TrimSpace(strings.ReplaceAll(s.Nodes[0].FirstChild.Data, ":", ""))
		price, err := strconv.ParseFloat(strings.ReplaceAll(s.Find("em").Text(), "€", ""), 64)
		if err != nil {
			loopErr = err
			return false
		}

		if fuelTypeString == "Euro 95" {
			fuelType = models.Euro95E5
		} else if fuelTypeString == "Diesel" {
			fuelType = models.Diesel
		} else {
			loopErr = fmt.Errorf("unable to convert `%+v` to a FuelType in Esso", fuelTypeString)
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
