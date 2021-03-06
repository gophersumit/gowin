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
	"log"

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
		result := getResults(city, pincode)
		printCenters(result)

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
