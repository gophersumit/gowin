package cmd

import "github.com/fatih/color"

func printCenters(cowinCenters []CowinCenters) {
	header := color.New(color.BgBlack).Add(color.FgWhite).Add(color.Bold)
	header.Printf("%30s|\t%10s|\t%10s|\t%10s|\t%s|\t%8s\n", "Center Name", "City", "Vaccine", "Date", "Capacity", "Pincode")

	green := color.New(color.FgGreen).Add(color.Underline).Add(color.BgBlack).Add(color.Bold)

	for _, center := range cowinCenters {
		for _, v := range center.Centers {
			for _, s := range v.Sessions {
				if s.AvailableCapacity > 0 {
					green.Printf("%30s|\t%10s|\t%10s|\t%10s|\t%8d|\t%8d|\n", v.Name, v.DistrictName, s.Vaccine, s.Date, s.AvailableCapacity, v.Pincode)
				} else {
					//		red.Printf("%30s|\t%10s|\t%10s|\t%10s|\t%8d|\n", v.Name, v.DistrictName, s.Vaccine, s.Date, s.AvailableCapacity)
				}
			}
		}
	}
}
