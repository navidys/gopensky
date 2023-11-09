package main

import (
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var flightsCommand = &cobra.Command{
	Use:     "flights",
	Short:   "Retrieve flights for a certain time interval",
	Run:     runFlightsCommand,
	PreRunE: preRunFlights,
}

func registerFlightsCommand() {
	// flights command
	flightsCommand.Flags().Int64VarP(&cmdBeginTime, "being", "b", cmdBeginTime,
		"start of time interval to retrieve flights for as Unix time (seconds since epoch)")

	flightsCommand.Flags().Int64VarP(&cmdEndTime, "end", "e", cmdEndTime,
		"end of time interval to retrieve flights for as Unix time (seconds since epoch)")
}
