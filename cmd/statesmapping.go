package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gophersumit/gowin/pkg/CowinPublicV2"
)

type StatesResponse struct {
	States []struct {
		StateID   int    `json:"state_id"`
		StateName string `json:"state_name"`
	} `json:"states"`
	TTL int `json:"ttl"`
}

type DistrictResponse struct {
	Districts []struct {
		DistrictID   int    `json:"district_id"`
		DistrictName string `json:"district_name"`
	} `json:"districts"`
	TTL int `json:"ttl"`
}

type District struct {
	DistrictID   int    `json:"district_id"`
	DistrictName string `json:"district_name"`
}
type StateDistrict struct {
	StateName string     `json:"state_name"`
	Districts []District `json:"districts"`
}
type StateDistrictMaster map[int]*StateDistrict

var masterData StateDistrictMaster

func init() {
	fmt.Println("Loading states and cities master data...")

	acceptLanguage := "en-US,en;q=0.9"
	params := CowinPublicV2.StatesParams{
		AcceptLanguage: &acceptLanguage,
	}

	cowinClient := CowinPublicV2.Client{
		Server:         "https://cdn-api.co-vin.in/api/",
		Client:         &http.Client{},
		RequestEditors: []CowinPublicV2.RequestEditorFn{},
	}

	response, err := cowinClient.States(context.Background(), &params, cowinClient.RequestEditors...)

	if err != nil {
		log.Fatalf("Error")
	}
	defer response.Body.Close()
	statesResponse := StatesResponse{}
	err = json.NewDecoder(response.Body).Decode(&statesResponse)

	if err != nil {
		log.Fatalln("Unable to get states")
	}

	count := len(statesResponse.States)

	wg := sync.WaitGroup{}
	wg.Add(count)
	masterData = make(StateDistrictMaster, count)

	for _, s := range statesResponse.States {
		sd := &StateDistrict{
			StateName: s.StateName,
			Districts: nil,
		}
		masterData[s.StateID] = sd
	}

	for _, s := range statesResponse.States {
		dParams := CowinPublicV2.DistrictsParams{
			AcceptLanguage: &acceptLanguage,
		}

		go func(id int) {
			response, err := cowinClient.Districts(context.Background(), strconv.Itoa(id), &dParams, cowinClient.RequestEditors...)
			if err != nil {
				log.Fatalln("Error getting district data")
			}

			defer response.Body.Close()
			dResponse := DistrictResponse{}
			json.NewDecoder(response.Body).Decode(&dResponse)

			for _, d := range dResponse.Districts {
				dist := District{
					DistrictID:   d.DistrictID,
					DistrictName: d.DistrictName,
				}
				masterData[id].Districts = append(masterData[id].Districts, dist)
			}
			wg.Done()
		}(s.StateID)
	}

	wg.Wait()

	fmt.Println("loaded...")
}
