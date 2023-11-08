package main

import (
	"context"
	"fmt"
	"os"

	"github.com/navidys/gopensky"
)

func main() {
	conn, err := gopensky.NewConnection(context.Background(), "", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// retrieve all states information
	statesData, err := gopensky.GetStates(conn, 0, nil, nil, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	for _, state := range statesData.States {
		longitude := "<nil>"
		latitude := "<nil>"

		if state.Longitude != nil {
			longitude = fmt.Sprintf("%f", *state.Longitude)
		}

		if state.Latitude != nil {
			latitude = fmt.Sprintf("%f", *state.Latitude)
		}

		fmt.Printf("ICAO24: %s, Longitude: %s, Latitude: %s, Origin Country: %s \n",
			state.Icao24,
			longitude,
			latitude,
			state.OriginCountry,
		)
	}
}
