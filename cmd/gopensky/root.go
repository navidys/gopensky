package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	cmdUsername          string
	cmdPassword          string
	cmdIcao24List        []string
	cmdTime              int64 = -1
	cmdDebug                   = false
	cmdPrintJSON               = false
	cmdStatesExtended          = false
	cmdStatesBoundingBox []float64
	cmdAirport           string
	cmdAircraft          string
	cmdBeginTime         int64
	cmdEndTime           int64
	cmdPrintVersion      bool

	buildVersion  string
	buildRevision string

	rootCmd = &cobra.Command{
		Use:               "gopensky",
		Short:             "Query opensky network live API",
		Long:              "Query opensky network live API information (ADS-B and Mode S data)",
		Version:           fmt.Sprintf("%s-%s", buildVersion, buildRevision),
		PersistentPreRunE: preRun,
	}

	errInvalidTimeInput = errors.New("invalid time entry")
	errPasswordEntry    = errors.New("password entry error")
)

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
	rootCmd.PersistentFlags().BoolVarP(&cmdPrintVersion, "version", "v", cmdPrintVersion, "print version and exit")

	registerStatesCommand()
	registerArrivalsCommand()
	registerDeparturesCommand()
	registerFlightsCommand()
	registerAircraftCommand()

	rootCmd.AddCommand(statesCommand)
	rootCmd.AddCommand(arrivalsCommand)
	rootCmd.AddCommand(departuresCommand)
	rootCmd.AddCommand(flightsCommand)
	rootCmd.AddCommand(aircraftCommand)
}
