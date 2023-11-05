package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/navidys/gopensky"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var statesCommand = &cobra.Command{
	Use:     "states",
	Short:   "retrieve state vector information",
	Run:     runStates,
	PreRunE: preStateRun,
}

func preStateRun(cmd *cobra.Command, args []string) error {
	if cmdTime < -1 {
		return errInvalidTimeInput
	}

	return nil
}

func runStates(cmd *cobra.Command, args []string) {
	var boundingBoxOpts *gopensky.BoundingBoxOptions

	if len(cmdStatesBoundingBox) == 4 { //nolint:gomnd
		boundingBoxOpts = &gopensky.BoundingBoxOptions{
			Lamin: cmdStatesBoundingBox[0],
			Lomin: cmdStatesBoundingBox[1],
			Lamax: cmdStatesBoundingBox[2],
			Lomax: cmdStatesBoundingBox[3],
		}
	}

	conn, err := gopensky.NewConnection(context.Background(), cmdUsername, cmdPassword)
	if err != nil {
		log.Error().Msgf("%v", err)

		return
	}

	states, err := gopensky.GetStates(conn, cmdTime, cmdIcao24List, boundingBoxOpts, cmdStatesExtended)
	if err != nil {
		log.Error().Msgf("%v", err)

		return
	}

	if cmdPrintJSON {
		jsonResult, err := json.MarshalIndent(states, "", "    ")
		if err != nil {
			log.Error().Msgf("%v", err)

			return
		}

		fmt.Printf("%s\n", jsonResult) //nolint:forbidigo
	} else {
		printStatesTemplate(states)
	}
}

func printStatesTemplate(states *gopensky.States) { //nolint:funlen
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	header := fmt.Sprintf("\n%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
		"Icao24",
		"Callsign",
		"OriginCountry",
		"TimePosition",
		"LastContact",
		"Longitude",
		"Latitude",
		"BaroAltitude",
		"OnGround",
		"Velocity",
		"TrueTrack",
		"VerticalRate",
		"Sensors",
		"GeoAltitude",
		"Squawk",
		"Spi",
		"PositionSource",
		"Category",
	)

	fmt.Fprintln(writer, header)

	floatToString := func(data *float64) string {
		if data != nil {
			return fmt.Sprintf("%.4f", *data)
		}

		return "<nil>"
	}

	for _, state := range states.States {
		var (
			callSign     = "<nil>"
			timePosition = "<nil>"
			sensors      = "<nil>"
			squawk       = "<nil>"
		)

		if state.Squawk != nil {
			squawk = *state.Squawk
		}

		geoAltitude := floatToString(state.GeoAltitude)
		verticalRate := floatToString(state.VerticalRate)
		trueTrack := floatToString(state.TrueTrack)
		velocity := floatToString(state.Velocity)
		baroAltitude := floatToString(state.BaroAltitude)
		latitude := floatToString(state.Latitude)
		longitude := floatToString(state.Longitude)

		if state.Sensors != nil {
			sensors = fmt.Sprintf("%v", state.Sensors)
		}

		if state.TimePosition != nil {
			timePosition = fmt.Sprintf("%d", *state.TimePosition)
		}

		if state.Callsign != nil {
			callSign = *state.Callsign
		}

		data := fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%s\t%s\t%s\t%v\t%s\t%s\t%s\t%s\t%s\t%s\t%v\t%d\t%d",
			state.Icao24,
			callSign,
			state.OriginCountry,
			timePosition,
			state.LastContact,
			longitude,
			latitude,
			baroAltitude,
			state.OnGround,
			velocity,
			trueTrack,
			verticalRate,
			sensors,
			geoAltitude,
			squawk,
			state.Spi,
			state.PositionSource,
			state.Category,
		)

		fmt.Fprintln(writer, data)
	}

	if err := writer.Flush(); err != nil {
		log.Error().Msgf("failed to flush template: %v", err)
	}
}
