package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/navidys/gopensky"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func preRunFlightArrivalsDepartures(cmd *cobra.Command, args []string) error {
	if strings.TrimSpace(cmdAirport) == "" {
		return gopensky.ErrInvalidAirportName
	}

	if cmdBeginTime <= 0 || cmdEndTime <= 0 {
		return gopensky.ErrInvalidUnixTime
	}

	return nil
}

func preRunAircraft(cmd *cobra.Command, args []string) error {
	if strings.TrimSpace(cmdAircraft) == "" {
		return gopensky.ErrInvalidAircraftName
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

func preRunTracks(cmd *cobra.Command, args []string) error {
	if cmdTime < 0 {
		return gopensky.ErrInvalidUnixTime
	}

	if strings.TrimSpace(cmdAircraft) == "" {
		return gopensky.ErrInvalidAircraftName
	}

	return nil
}

func preRun(cmd *cobra.Command, args []string) error {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if cmdDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if cmdUsername != "" {
		fmt.Print("Enter Password: ") //nolint:forbidigo

		bytePassword, err := term.ReadPassword(0)
		if err != nil {
			return errPasswordEntry
		}

		cmdPassword = strings.TrimSpace(string(bytePassword))
	}

	return nil
}
