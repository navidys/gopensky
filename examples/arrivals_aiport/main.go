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

	// retrieve arrivals flights of:
	// airpor: LFPG (Charles de Gaulle)
	// being time: 1696755342 (Sunday October 08, 2023 08:55:42 UTC)
	// end time: 1696928142 (Tuesday October 10, 2023 08:55:42 UTC)

	flightsData, err := gopensky.GetArrivalsByAirport(conn, "LFPG", 1696755342, 1696928142)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	for _, flightData := range flightsData {
		var (
			depAirport string
			callSign   string
		)

		if flightData.EstDepartureAirport != nil {
			depAirport = *flightData.EstDepartureAirport
		}

		if flightData.Callsign != nil {
			callSign = *flightData.Callsign
		}

		fmt.Printf("ICAO24: %s, callSign: %8s, Departed Airport: %4s, LastSeen: %s\n",
			flightData.Icao24,
			callSign,
			depAirport,
			time.Unix(flightData.LastSeen, 0),
		)
	}
}
