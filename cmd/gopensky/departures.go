package main

import (
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var departuresCommand = &cobra.Command{
	Use:     "departures",
	Short:   "Retrieve flights for a certain airport which departed within a given time interval",
	Run:     runFlightsCommand,
	PreRunE: preRunFlightArrivalsDepartures,
}

func registerDeparturesCommand() {
	// departures command
	departuresCommand.Flags().StringVarP(&cmdAirport, "airport", "a", cmdAirport,
		"ICAO identifier for the airport")

	departuresCommand.Flags().Int64VarP(&cmdBeginTime, "being", "b", cmdBeginTime,
		"start of time interval to retrieve flights for as Unix time (seconds since epoch)")

	departuresCommand.Flags().Int64VarP(&cmdEndTime, "end", "e", cmdEndTime,
		"end of time interval to retrieve flights for as Unix time (seconds since epoch)")
}
