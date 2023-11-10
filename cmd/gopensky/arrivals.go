package main

import (
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var arrivalsCommand = &cobra.Command{
	Use:   "arrivals",
	Short: "Retrieve flights for a certain airport which arrived within a given time interval (<= 7 days)",
	Long: `Retrieve flights for a certain airport which arrived within a given time interval.
The given time interval must not be larger than seven days!`,
	Run:     runFlightsCommand,
	PreRunE: preRunFlightArrivalsDepartures,
}

func registerArrivalsCommand() {
	// arrivals command
	arrivalsCommand.Flags().StringVarP(&cmdAirport, "airport", "a", cmdAirport,
		"ICAO identifier for the airport")

	arrivalsCommand.Flags().Int64VarP(&cmdBeginTime, "being", "b", cmdBeginTime,
		"start of time interval to retrieve flights for as unix time (seconds since epoch)")

	arrivalsCommand.Flags().Int64VarP(&cmdEndTime, "end", "e", cmdEndTime,
		"end of time interval to retrieve flights for as unix time (seconds since epoch)")
}
