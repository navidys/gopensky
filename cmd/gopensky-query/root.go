package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

//nolint:gochecknoglobals
var (
	cmdUsername          string
	cmdPassword          string
	cmdIcao24List        string
	cmdTime              int64 = -1
	cmdDebug                   = false
	cmdPrintJSON               = false
	cmdStatesExtended          = false
	cmdStatesBoundingBox []float64

	rootCmd = &cobra.Command{
		Use:               "gopensky-query",
		Short:             "Query opensky network live API",
		Long:              "Query opensky network live API information (ADS-B and Mode S data)",
		PersistentPreRunE: preRun,
	}

	errInvalidTimeInput = errors.New("invalid time entry")
	errPasswordEntry    = errors.New("password entry error")
)

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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() { //nolint:gochecknoinits
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVarP(&cmdUsername, "username", "u", cmdUsername, "connection username")
	rootCmd.PersistentFlags().BoolVarP(&cmdDebug, "debug", "d", cmdDebug, "run in debug mode")
	rootCmd.PersistentFlags().BoolVarP(&cmdPrintJSON, "json", "j", cmdPrintJSON, "print json output")

	// states command
	statesCommand.Flags().StringVarP(&cmdIcao24List, "icao24", "i", "",
		"comma separates (,) list of unique ICAO 24-bit address of the transponder in hex string representation")

	statesCommand.Flags().Int64VarP(&cmdTime, "time", "t", cmdTime,
		"the time which the state vectors in this response are associated with (current time will be used if omitted)")

	statesCommand.Flags().BoolVarP(&cmdStatesExtended, "extended", "e", cmdStatesExtended,
		"request the category of aircraft ")

	statesCommand.Flags().Float64SliceVar(&cmdStatesBoundingBox, "box", nil,
		"query a certain area defined by a bounding box of WGS84 coordinates ([lamin,lomin,lamax,lomax])")

	rootCmd.AddCommand(statesCommand)
}
