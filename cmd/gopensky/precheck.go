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

func preStateRun(cmd *cobra.Command, args []string) error {
	if cmdTime < 0 {
		return gopensky.ErrInvalidUnixTime
	}

	return nil
}

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

	if cmdUsername != "" {
		fmt.Print("Enter password: ") //nolint:forbidigo

		bytePassword, err := term.ReadPassword(0)
		if err != nil {
			return errPasswordEntry
		}

		fmt.Println() //nolint:forbidigo

		cmdPassword = strings.TrimSpace(string(bytePassword))
	}

	return nil
}
