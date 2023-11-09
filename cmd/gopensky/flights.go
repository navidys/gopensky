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
var flightsCommand = &cobra.Command{
	Use:     "flights",
	Short:   "Retrieve flights for a certain time interval ",
	Run:     runFlights,
	PreRunE: preRunFlights,
}

func runFlights(cmd *cobra.Command, args []string) {
	conn, err := gopensky.NewConnection(context.Background(), cmdUsername, cmdPassword)
	if err != nil {
		log.Error().Msgf("%v", err)

		return
	}

	flights, err := gopensky.GetFlightsByInterval(conn, cmdBeginTime, cmdEndTime)
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
