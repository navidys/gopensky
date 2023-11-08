package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

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
		"LastSeen",
		"EstDepartureAirp",
		"EstArrivalAirp",
		"Callsign",
		"EstDepAirpHorizDist",
		"EstDepAirpVertDist",
		"EstArrAirpHorizDist",
		"EstArrAirpVertDist",
		"DepAirportCandCount",
		"ArrAirportCandCount",
	)

	fmt.Fprintln(writer, header)

	for _, flightData := range flights {
		firstSeen := time.Unix(flightData.FirstSeen, 0).Format("2023-10-09 09:52:38")
		lastSeen := time.Unix(flightData.LastSeen, 0).Format("2023-10-09 09:52:38")
		estDepartureAirport := ""
		estArrivalAirport := ""
		callsign := ""

		if flightData.EstDepartureAirport != nil {
			estDepartureAirport = *flightData.EstDepartureAirport
		}

		if flightData.EstArrivalAirport != nil {
			estDepartureAirport = *flightData.EstArrivalAirport
		}

		if flightData.Callsign != nil {
			estDepartureAirport = *flightData.Callsign
		}

		data := fmt.Sprintf("\n%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d",
			flightData.Icao24,
			firstSeen,
			lastSeen,
			estDepartureAirport,
			estArrivalAirport,
			callsign,
			flightData.EstDepartureAirportVertDistance,
			flightData.EstDepartureAirportVertDistance,
			flightData.EstArrivalAirportHorizDistance,
			flightData.EstArrivalAirportVertDistance,
			flightData.DepartureAirportCandidatesCount,
			flightData.ArrivalAirportCandidatesCount,
		)

		fmt.Fprintln(writer, data)
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
