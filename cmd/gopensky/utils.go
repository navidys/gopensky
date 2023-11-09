package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/navidys/gopensky"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func preRunFlightArrivalsDepartures(cmd *cobra.Command, args []string) error {
	if strings.TrimSpace(cmdAirport) == "" {
		return gopensky.ErrEmptyAirportName
	}

	if cmdBeginTime <= 0 || cmdEndTime <= 0 {
		return gopensky.ErrInvalidUnixTime
	}

	return nil
}

func preRunFlights(cmd *cobra.Command, args []string) error {
	if cmdBeginTime <= 0 || cmdEndTime <= 0 {
		return gopensky.ErrInvalidUnixTime
	}

	return nil
}

func printFlightTable(flights []gopensky.FlighData) { //nolint:funlen
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	info := []string{
		"", // just to print a new line
		"EDA   = EstDepartureAirport",
		"EAA   = EstArrivalAirport",
		"EDAHD = EstDepartureAirportHorizDistance",
		"EDAVD = EstDepartureAirportVertDistance",
		"EAAHD = EstArrivalAirportHorizDistance",
		"EAAVD = EstArrivalAirportVertDistance",
		"DACC  = DepartureAirportCandidatesCount",
		"AACC  = ArrivalAirportCandidatesCount",
		"", // just to print a new line
	}

	if _, err := fmt.Fprintln(os.Stdout, strings.Join(info, "\n")); err != nil {
		log.Error().Msgf("%v", err)
	}

	header := fmt.Sprintf("%s\t%s\t%s\t%8s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
		"Icao24", "EDA", "EAA", "Callsign", "EDAHD", "EDAVD",
		"EAAHD", "EAAVD", "DACC", "AACC", "FirstSeen", "LastSeen")

	if _, err := fmt.Fprintln(writer, header); err != nil {
		log.Error().Msgf("%v", err)
	}

	for _, flightData := range flights {
		firstSeen := time.Unix(flightData.FirstSeen, 0).UTC()
		lastSeen := time.Unix(flightData.LastSeen, 0).UTC()
		estDepartureAirport := ""
		estArrivalAirport := ""
		callsign := ""

		if flightData.EstDepartureAirport != nil {
			estDepartureAirport = *flightData.EstDepartureAirport
		}

		if flightData.EstArrivalAirport != nil {
			estArrivalAirport = *flightData.EstArrivalAirport
		}

		if flightData.Callsign != nil {
			estDepartureAirport = *flightData.Callsign
		}

		data := fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%s\t%s",
			flightData.Icao24,
			estDepartureAirport,
			estArrivalAirport,
			callsign,
			flightData.EstDepartureAirportVertDistance,
			flightData.EstDepartureAirportVertDistance,
			flightData.EstArrivalAirportHorizDistance,
			flightData.EstArrivalAirportVertDistance,
			flightData.DepartureAirportCandidatesCount,
			flightData.ArrivalAirportCandidatesCount,
			firstSeen,
			lastSeen,
		)

		if _, err := fmt.Fprintln(writer, data); err != nil {
			log.Error().Msgf("%v", err)
		}
	}

	if err := writer.Flush(); err != nil {
		log.Error().Msgf("failed to flush template: %v", err)
	}
}
