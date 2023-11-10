package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/navidys/gopensky"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func runStatesCommand(cmd *cobra.Command, args []string) {
	var boundingBoxOpts *gopensky.BoundingBoxOptions

	if len(cmdStatesBoundingBox) == 4 { //nolint:gomnd
		boundingBoxOpts = &gopensky.BoundingBoxOptions{
			Lamin: cmdStatesBoundingBox[0],
			Lomin: cmdStatesBoundingBox[1],
			Lamax: cmdStatesBoundingBox[2],
			Lomax: cmdStatesBoundingBox[3],
		}
	}

	conn, err := gopensky.NewConnection(context.Background(), cmdUsername, cmdPassword)
	if err != nil {
		log.Error().Msgf("%v", err)

		return
	}

	states, err := gopensky.GetStates(conn, cmdTime, cmdIcao24List, boundingBoxOpts, cmdStatesExtended)
	if err != nil {
		log.Error().Msgf("%v", err)

		return
	}

	if cmdPrintJSON {
		jsonResult, err := json.MarshalIndent(states, "", "    ")
		if err != nil {
			log.Error().Msgf("%v", err)

			return
		}

		fmt.Printf("%s\n", jsonResult) //nolint:forbidigo
	} else {
		printStatesTable(states)
	}
}

func runFlightsCommand(cmd *cobra.Command, args []string) {
	var (
		flightsData []gopensky.FlighData
		err         error
		conn        context.Context
	)

	conn, err = gopensky.NewConnection(context.Background(), cmdUsername, cmdPassword)
	if err != nil {
		log.Error().Msgf("%v", err)

		return
	}

	switch cmd.Use {
	case "flights":
		flightsData, err = gopensky.GetFlightsByInterval(conn, cmdBeginTime, cmdEndTime)
	case "arrivals":
		flightsData, err = gopensky.GetArrivalsByAirport(conn, cmdAirport, cmdBeginTime, cmdEndTime)
	case "departures":
		flightsData, err = gopensky.GetDeparturesByAirport(conn, cmdAirport, cmdBeginTime, cmdEndTime)
	case "aircraft":
		flightsData, err = gopensky.GetFlightsByAircraft(conn, cmdAircraft, cmdBeginTime, cmdEndTime)
	default:
		log.Error().Msg("unsupported api /flights operation")
	}

	if err != nil {
		log.Error().Msgf("%v", err)

		return
	}

	if cmdPrintJSON {
		jsonResult, err := json.MarshalIndent(flightsData, "", "    ")
		if err != nil {
			log.Error().Msgf("%v", err)

			return
		}

		fmt.Printf("%s\n", jsonResult) //nolint:forbidigo
	} else {
		printFlightsTable(flightsData)
	}
}

func runTracksCommand(cmd *cobra.Command, args []string) {
}
