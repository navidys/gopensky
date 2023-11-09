package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/navidys/gopensky"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var arrivalsCommand = &cobra.Command{
	Use:     "arrivals",
	Short:   "Retrieve flights for a certain airport which arrived within a given time interval",
	Run:     runArrivals,
	PreRunE: preRunFlightArrivalsDepartures,
}

func runArrivals(cmd *cobra.Command, args []string) {
	conn, err := gopensky.NewConnection(context.Background(), cmdUsername, cmdPassword)
	if err != nil {
		log.Error().Msgf("%v", err)

		return
	}

	flights, err := gopensky.GetArrivalsByAirport(conn, cmdAirport, cmdBeginTime, cmdEndTime)
	if err != nil {
		log.Error().Msgf("%v", err)

		return
	}

	if cmdPrintJSON {
		jsonResult, err := json.MarshalIndent(flights, "", "    ")
		if err != nil {
			log.Error().Msgf("%v", err)

			return
		}

		fmt.Printf("%s\n", jsonResult) //nolint:forbidigo
	} else {
		printFlightTable(flights)
	}
}
