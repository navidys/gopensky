package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/navidys/gopensky"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var arrivalsCommand = &cobra.Command{
	Use:     "arrivals",
	Short:   "Retrieve flights for a certain airport which arrived within a given time interval",
	Run:     runArrivals,
	PreRunE: preArrivalsRun,
}

func runArrivals(cmd *cobra.Command, args []string) {
	conn, err := gopensky.NewConnection(context.Background(), cmdUsername, cmdPassword)
	if err != nil {
		log.Error().Msgf("%v", err)

		return
	}

	flights, err := gopensky.GetArrivalsByAirport(conn, "ymml", cmdBeginTime, cmdEndTime)
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
		printArrivalsTemplate(flights)
	}
}

func printArrivalsTemplate(flights []gopensky.FlighData) {
}

func preArrivalsRun(cmd *cobra.Command, args []string) error {
	if strings.TrimSpace(cmdAirport) == "" {
		return gopensky.ErrEmptyAirportName
	}

	if cmdBeginTime <= 0 || cmdEndTime <= 0 {
		return gopensky.ErrInvalidUnixTime
	}

	return nil
}
