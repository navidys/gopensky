package main

import (
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var aircraftCommand = &cobra.Command{
	Use:     "aircraft",
	Short:   "Retrieve flights for a particular aircraft within a certain time interval",
	Run:     runFlightsCommand,
	PreRunE: preRunAircraft,
}

func registerAircraftCommand() {
	// aircraft command
	aircraftCommand.Flags().StringVarP(&cmdAircraft, "aircraft", "a", cmdAircraft,
		"unique ICAO 24-bit address of the transponder in hex string representation.")

	aircraftCommand.Flags().Int64VarP(&cmdBeginTime, "being", "b", cmdBeginTime,
		"start of time interval to retrieve flights for as Unix time (seconds since epoch)")

	aircraftCommand.Flags().Int64VarP(&cmdEndTime, "end", "e", cmdEndTime,
		"end of time interval to retrieve flights for as Unix time (seconds since epoch)")
}
