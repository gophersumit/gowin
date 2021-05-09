/*
Copyright © 2021 Sumit Agrawal gophersumit@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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

	"github.com/fatih/color"
	"github.com/gophersumit/gowin/pkg/CowinPublicV2"
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "find available hospitals by pincode using cowin",
	Long:  `find will help you see hopsitals with vaccinate availabity for next 30 days.`,
	Run: func(cmd *cobra.Command, args []string) {

		pincode, err := cmd.Flags().GetInt("pincode")
		if err != nil {
			log.Fatalf("Error reading pincode")
		}
		city, err := cmd.Flags().GetString("city")
		if err != nil {
			log.Fatalf("Error reading pincode")
		}

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

		weeks := 12
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
		header := color.New(color.BgBlack).Add(color.FgWhite).Add(color.Bold)
		header.Printf("%30s|\t%10s|\t%10s|\t%10s|\t%s|\t%8s\n", "Center Name", "City", "Vaccine", "Date", "Capacity", "Pincode")
		wg.Add(len(dates))
		for _, date := range dates {
			go func(d string) {

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
					printCenters(centers)
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
					printCenters(centers)
				}
				wg.Done()

			}(date)
		}

		wg.Wait()
	},
}

func init() {
	var pincode int

	var city string
	findCmd.Flags().IntVarP(&pincode, "pincode", "p", 0, "Pincode to search for")
	findCmd.Flags().StringVarP(&city, "city", "c", "", "City to search for")

	rootCmd.AddCommand(findCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printCenters(centers CowinCenters) {
	green := color.New(color.FgGreen).Add(color.Underline).Add(color.BgBlack).Add(color.Bold)
	//	red := color.New(color.FgRed).Add(color.Underline).Add(color.BgBlack).Add(color.Bold)

	for _, v := range centers.Centers {
		for _, s := range v.Sessions {
			if s.AvailableCapacity > 0 {
				green.Printf("%30s|\t%10s|\t%10s|\t%10s|\t%8d|\t%8d|\n", v.Name, v.DistrictName, s.Vaccine, s.Date, s.AvailableCapacity, v.Pincode)
			} else {
				//		red.Printf("%30s|\t%10s|\t%10s|\t%10s|\t%8d|\n", v.Name, v.DistrictName, s.Vaccine, s.Date, s.AvailableCapacity)
			}
		}
	}
}
