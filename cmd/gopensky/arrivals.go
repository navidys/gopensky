package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

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
		printArrivalsTemplate(flights)
	}
}

func printArrivalsTemplate(flights []gopensky.FlighData) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	header := fmt.Sprintf("\n%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
		"Icao24",
		"FirstSeen",
		"EstDepartureAirport",
		"LastSeen",
		"EstArrivalAirport",
		"Callsign",
		"EstDepAirportHorizDist",
		"EstDepAirportVertDist",
		"EstArrAirportHorizDist",
		"EstArrAirportVertDist",
		"DepAirportCandCount",
		"ArrAirportCandCount",
	)

	fmt.Fprintln(writer, header)

	if err := writer.Flush(); err != nil {
		log.Error().Msgf("failed to flush template: %v", err)
	}
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
