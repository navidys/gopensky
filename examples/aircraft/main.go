package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/navidys/gopensky"
)

func main() {
	conn, err := gopensky.NewConnection(context.Background(), "", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Retrieve Elon Musk's primary private jet flight data, he has many and one of
	// them has the ICAO24 transponder address a835af
	// begin time: 1693523464 (Thu Aug 31 2023 23:11:04)
	// end time: 1696029064 (Fri Sep 29 2023 23:11:04)

	flightsData, err := gopensky.GetFlightsByAircraft(conn, "a835af", 1693523464, 1696029064)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	for _, flightData := range flightsData {
		var (
			departedAirport string
			arrivalAriport  string
		)

		if flightData.EstDepartureAirport != nil {
			departedAirport = *flightData.EstDepartureAirport
		}

		if flightData.EstArrivalAirport != nil {
			arrivalAriport = *flightData.EstArrivalAirport
		}

		fmt.Printf("ICAO24: %s, Departed: %4s, Arrival: %4s, LastSeen: %s\n",
			flightData.Icao24,
			departedAirport,
			arrivalAriport,
			time.Unix(flightData.LastSeen, 0),
		)
	}
}
