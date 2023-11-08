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
	cmdIcao24List        []string
	cmdTime              int64 = -1
	cmdDebug                   = false
	cmdPrintJSON               = false
	cmdStatesExtended          = false
	cmdStatesBoundingBox []float64
	cmdAirport           string
	cmdBeginTime         int64
	cmdEndTime           int64
	cmdPrintVersion      bool

	buildVersion  string
	buildRevision string

	rootCmd = &cobra.Command{
		Use:               "gopensky",
		Short:             "Query opensky network live API",
		Long:              "Query opensky network live API information (ADS-B and Mode S data)",
		PersistentPreRunE: preRun,
		Run:               run,
	}

	errInvalidTimeInput = errors.New("invalid time entry")
	errPasswordEntry    = errors.New("password entry error")
)

func run(cmd *cobra.Command, args []string) {
	if cmdPrintVersion {
		fmt.Printf("%s version %s-%s\n", cmd.Use, buildVersion, buildRevision) //nolint:forbidigo

		os.Exit(0)
	}
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

	// states command
	statesCommand.Flags().StringSliceVarP(&cmdIcao24List, "icao24", "i", cmdIcao24List,
		"comma separates (,) list of unique ICAO 24-bit address of the transponder in hex string representation")

	statesCommand.Flags().Int64VarP(&cmdTime, "time", "t", cmdTime,
		"the time which the state vectors in this response are associated with (current time will be used if omitted)")

	statesCommand.Flags().BoolVarP(&cmdStatesExtended, "extended", "e", cmdStatesExtended,
		"request the category of aircraft ")

	statesCommand.Flags().Float64SliceVar(&cmdStatesBoundingBox, "box", nil,
		"query a certain area defined by a bounding box of WGS84 coordinates (lamin,lomin,lamax,lomax)")

	// flights command
	arrivalsCommand.Flags().StringVarP(&cmdAirport, "airport", "a", cmdAirport,
		"ICAO identifier for the airport")

	arrivalsCommand.Flags().Int64VarP(&cmdBeginTime, "being", "b", cmdBeginTime,
		"start of time interval to retrieve flights for as Unix time (seconds since epoch)")

	arrivalsCommand.Flags().Int64VarP(&cmdEndTime, "end", "e", cmdEndTime,
		"end of time interval to retrieve flights for as Unix time (seconds since epoch)")

	rootCmd.AddCommand(statesCommand)
	rootCmd.AddCommand(arrivalsCommand)
}
