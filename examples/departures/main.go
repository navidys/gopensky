package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/navidys/gopensky"
	"github.com/rs/zerolog/log"
)

func main() {
	conn, err := gopensky.NewConnection(context.Background(), "", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// retrieve departed flights of:
	// airpor: LFPG (Charles de Gaulle)
	// being time: 1696755342 (Sunday October 08, 2023 08:55:42 UTC)
	// end time: 1696928142 (Tuesday October 10, 2023 08:55:42 UTC)

	flightsData, err := gopensky.GetDeparturesByAirport(conn, "LFPG", 1696755342, 1696928142)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	jsonResult, err := json.MarshalIndent(flightsData, "", "    ")
	if err != nil {
		log.Error().Msgf("%v", err)

		return
	}

	fmt.Printf("%s\n", jsonResult)
}
