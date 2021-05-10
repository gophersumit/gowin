package cmd

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gophersumit/gowin/pkg/CowinPublicV2"
)

func getResults(city string, pincode int) []CowinCenters {

	weeks := 12
	cowinCenters := make([]CowinCenters, weeks)

	cityId := 0
	for _, s := range masterData {
		for _, d := range s.Districts {
			if strings.EqualFold(city, d.DistrictName) {
				cityId = d.DistrictID
				break
			}
		}
	}

	if cityId == 0 && city != "" {
		log.Fatalln("Invalid City")
	}

	dates := make([]string, 12)
	for i := 0; i < weeks; i++ {
		dates[i] = time.Now().AddDate(0, 0, 7*i).Format("02-01-2006")
	}

	strPin := strconv.Itoa(pincode)
	cowinClient := CowinPublicV2.Client{
		Server:         "https://cdn-api.co-vin.in/api/",
		Client:         &http.Client{},
		RequestEditors: []CowinPublicV2.RequestEditorFn{},
	}
	wg := sync.WaitGroup{}

	wg.Add(len(dates))
	for i, date := range dates {
		go func(d string, i int) {

			if cityId > 0 {
				param := &CowinPublicV2.CalendarByDistrictParams{
					DistrictId:     strconv.Itoa(cityId),
					Date:           d,
					AcceptLanguage: nil,
				}

				response, err := cowinClient.CalendarByDistrict(context.Background(), param, cowinClient.RequestEditors...)

				if err != nil {
					log.Fatalf("Error")
				}
				defer response.Body.Close()
				centers := CowinCenters{}
				json.NewDecoder(response.Body).Decode(&centers)
				cowinCenters[i] = centers
			} else {
				param := &CowinPublicV2.CalendarByPinParams{
					Pincode:        strPin,
					Date:           d,
					AcceptLanguage: nil,
				}
				response, err := cowinClient.CalendarByPin(context.Background(), param, cowinClient.RequestEditors...)

				if err != nil {
					log.Fatalf("Error")
				}
				defer response.Body.Close()
				centers := CowinCenters{}
				json.NewDecoder(response.Body).Decode(&centers)
				cowinCenters[i] = centers
			}
			wg.Done()

		}(date, i)
	}

	wg.Wait()

	return cowinCenters
}
