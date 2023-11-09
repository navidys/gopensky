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

	// retrieve flights
	// being time: 1696832368 (Monday, October 9, 2023 6:19:28)
	// end time: 1696835968 (Monday, October 9, 2023 7:19:28)

	flightsData, err := gopensky.GetFlightsByInterval(conn, 1696832368, 1696835968)
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
