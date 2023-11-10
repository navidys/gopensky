package main

import "github.com/spf13/cobra"

//nolint:gochecknoglobals
var tracksCommand = &cobra.Command{
	Use:   "tracks",
	Short: "Retrieve the trajectory for a certain aircraft at a given time (<= 30 days)",
	Long: `Retrieve the trajectory for a certain aircraft at a given time.
The trajectory is a list of waypoints containing position, barometric altitude, true track and an on-ground flag.
The given time interval must not be larger than 30 days!`,
	Run:     runTracksCommand,
	PreRunE: preRunTracks,
}

func registerTracksCommand() {
	// flights command
	tracksCommand.Flags().Int64VarP(&cmdTime, "time", "t", cmdTime,
		`unix time in seconds since epoch. It can be any time between start and end of a known flight
		 If time = 0, get the live track if there is any flight ongoing for the given aircraft.
		`)

	tracksCommand.Flags().StringVarP(&cmdAircraft, "aircraft", "a", cmdAircraft,
		"unique ICAO 24-bit address of the transponder in hex string representation.")
}
