package main

import (
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var flightsCommand = &cobra.Command{
	Use:   "flights",
	Short: "Retrieve flights for a certain time interval (<= 2 hours)",
	Long: `Retrieve flights for a certain time interval.
The given time interval must not be larger than 30 days!`,
	Run:     runFlightsCommand,
	PreRunE: preRunFlights,
}

func registerFlightsCommand() {
	// flights command
	flightsCommand.Flags().Int64VarP(&cmdBeginTime, "being", "b", cmdBeginTime,
		"start of time interval to retrieve flights for as unix time (seconds since epoch)")

	flightsCommand.Flags().Int64VarP(&cmdEndTime, "end", "e", cmdEndTime,
		"end of time interval to retrieve flights for as unix time (seconds since epoch)")
}
