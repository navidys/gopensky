package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/navidys/gopensky"
	"github.com/rs/zerolog/log"
)

func printFlightsTable(flights []gopensky.FlighData) { //nolint:funlen
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	info := []string{
		"Header aliases:",
		"  EDA   = estimated departure airport",
		"  EAA   = estimated arrival airport",
		"  EDAHD = estimated departure airport horizontal distance",
		"  EDAVD = estimated departure airport vertical distance",
		"  EAAHD = estimated arrival airport horizontal distance",
		"  EAAVD = estimated arrival airport vertical distance",
		"  DACC  = departure airport candidates count",
		"  AACC  = arrival airport candidates count",
		"", // just to print a new line
	}

	if _, err := fmt.Fprintln(os.Stdout, strings.Join(info, "\n")); err != nil {
		log.Error().Msgf("%v", err)
	}

	header := fmt.Sprintf("%s\t%s\t%s\t%8s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
		"ICAO24", "EDA", "EAA", "CALLSIGN", "EDAHD", "EDAVD",
		"EAAHD", "EAAVD", "DACC", "AACC", "FIRST SEEN", "LAST SEEN")

	if _, err := fmt.Fprintln(writer, header); err != nil {
		log.Error().Msgf("%v", err)
	}

	for _, flightData := range flights {
		firstSeen := time.Unix(flightData.FirstSeen, 0).UTC()
		lastSeen := time.Unix(flightData.LastSeen, 0).UTC()
		estDepartureAirport := ""
		estArrivalAirport := ""
		callsign := ""

		if flightData.EstDepartureAirport != nil {
			estDepartureAirport = *flightData.EstDepartureAirport
		}

		if flightData.EstArrivalAirport != nil {
			estArrivalAirport = *flightData.EstArrivalAirport
		}

		if flightData.Callsign != nil {
			estDepartureAirport = *flightData.Callsign
		}

		data := fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%s\t%s",
			flightData.Icao24,
			estDepartureAirport,
			estArrivalAirport,
			callsign,
			flightData.EstDepartureAirportVertDistance,
			flightData.EstDepartureAirportVertDistance,
			flightData.EstArrivalAirportHorizDistance,
			flightData.EstArrivalAirportVertDistance,
			flightData.DepartureAirportCandidatesCount,
			flightData.ArrivalAirportCandidatesCount,
			firstSeen,
			lastSeen,
		)

		if _, err := fmt.Fprintln(writer, data); err != nil {
			log.Error().Msgf("%v", err)
		}
	}

	if err := writer.Flush(); err != nil {
		log.Error().Msgf("failed to flush template: %v", err)
	}
}

func printStatesTable(states *gopensky.States) { //nolint:funlen
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

	if _, err := fmt.Fprintln(writer, header); err != nil {
		log.Error().Msgf("%v", err)
	}

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

		if _, err := fmt.Fprintln(writer, data); err != nil {
			log.Error().Msgf("%v", err)
		}
	}

	if err := writer.Flush(); err != nil {
		log.Error().Msgf("failed to flush template: %v", err)
	}
}
