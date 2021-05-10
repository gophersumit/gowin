/*
Copyright © 2021 Sumit Agrawal <gophersumit@gmail.com>

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
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gophersumit/gowin/pkg/notify"
	"github.com/spf13/cobra"
)

// notifyCmd represents the notify command
var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		pincode, err := cmd.Flags().GetInt("pincode")
		if err != nil {
			log.Fatalf("Error reading pincode")
		}
		city, err := cmd.Flags().GetString("city")
		if err != nil {
			log.Fatalf("Error reading pincode")
		}
		fmt.Printf("Starting to monitor api\n")
		monitor(city, pincode, time.Now())
		doEvery(2*time.Minute, city, pincode, monitor)

	},
}

func init() {
	var pincode int
	var city string
	notifyCmd.Flags().IntVarP(&pincode, "pincode", "p", 0, "Pincode to monitor")
	notifyCmd.Flags().StringVarP(&city, "city", "c", "", "City to search monitor")
	rootCmd.AddCommand(notifyCmd)

}

func doEvery(d time.Duration, city string, pincode int, f func(city string, pincode int, x time.Time)) {
	for x := range time.Tick(d) {
		f(city, pincode, x)
	}
}

func monitor(city string, pincode int, x time.Time) {
	fmt.Printf("Started Monitoring at %s\n", x.Format("Jan 2006 15:04:05"))
	result := getResults(city, pincode)
	messages := []string{}
	for _, cowinCenter := range result {
		for _, c := range cowinCenter.Centers {
			for _, session := range c.Sessions {
				if session.AvailableCapacity > 0 {
					messages = append(messages, fmt.Sprintf("%d %s available at %s in %s on %s\n", session.AvailableCapacity, session.Vaccine, c.Name, c.DistrictName, session.Date))
				}
			}
		}
	}
	if len(messages) > 0 {
		messagesToshow := 3
		for i := 0; i < len(messages) && i < messagesToshow; i++ {
			notify.Alert("Vaccine Available!", messages[i])
		}
		fmt.Printf("Found %d centers.Stopping monitoring.To see complete list of centers try\n ", len(messages))
		green := color.New(color.FgGreen).Add(color.Underline).Add(color.BgBlack).Add(color.Bold).Add(color.Italic)
		if city != "" {
			green.Printf("gowin find -c=\"%s\"\n", city)
		}

		if pincode != 0 {
			green.Printf("gowin find -p=%d\n", pincode)
		}
		os.Exit(1)
	} else {
		fmt.Printf("No Vaccine available at %s\n", x.Format("Jan 2006 15:04:05"))
	}
}
