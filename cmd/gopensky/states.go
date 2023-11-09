package main

import (
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var statesCommand = &cobra.Command{
	Use:     "states",
	Short:   "retrieve state vector information",
	Run:     runStatesCommand,
	PreRunE: preStateRun,
}

func preStateRun(cmd *cobra.Command, args []string) error {
	if cmdTime < -1 {
		return errInvalidTimeInput
	}

	return nil
}

func registerStatesCommand() {
	// states command
	statesCommand.Flags().StringSliceVarP(&cmdIcao24List, "icao24", "i", cmdIcao24List,
		"comma separates (,) list of unique ICAO 24-bit address of the transponder in hex string representation")

	statesCommand.Flags().Int64VarP(&cmdTime, "time", "t", cmdTime,
		"the time which the state vectors in this response are associated with (current time will be used if omitted)")

	statesCommand.Flags().BoolVarP(&cmdStatesExtended, "extended", "e", cmdStatesExtended,
		"request the category of aircraft ")

	statesCommand.Flags().Float64SliceVar(&cmdStatesBoundingBox, "box", nil,
		"query a certain area defined by a bounding box of WGS84 coordinates (lamin,lomin,lamax,lomax)")
}
