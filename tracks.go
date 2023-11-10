package gopensky

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
)

func GetTrackByAircraft(ctx context.Context, icao24 string, time int64) (FlightTrack, error) {
	var (
		flightTrack         FlightTrack
		flightTrackResponse FlightTrackResponse
	)

	if icao24 == "" {
		return flightTrack, ErrInvalidAircraftName
	}

	if time < 0 {
		return flightTrack, ErrInvalidAircraftName
	}

	conn, err := getClient(ctx)
	if err != nil {
		return flightTrack, fmt.Errorf("client: %w", err)
	}

	requestParams := getTracksRequestParams(time, icao24)

	response, err := conn.doGetRequest(ctx, "/tracks/all", requestParams)
	if err != nil {
		return flightTrack, fmt.Errorf("do request: %w", err)
	}

	defer response.Body.Close()

	if err := response.process(&flightTrackResponse); err != nil {
		return flightTrack, err
	}

	flightTrack, err = parseFlightTrackResponse(&flightTrackResponse)
	if err != nil {
		return flightTrack, fmt.Errorf("parse track: %w", err)
	}

	return flightTrack, nil
}

func parseFlightTrackResponse(response *FlightTrackResponse) (FlightTrack, error) {
	var flightTrack FlightTrack

	log.Info().Msgf("%v", response)

	flightTrack.Icao24 = response.Icao24
	flightTrack.Callsign = response.Callsign
	flightTrack.StartTime = time.Unix(int64(response.StartTime), 1000).Unix()
	flightTrack.EndTime = time.Unix(int64(response.EndTime), 1000).Unix()

	return flightTrack, nil
}

func getTracksRequestParams(time int64, icao24 string) url.Values {
	requestParams := make(url.Values)
	if time >= 0 {
		requestParams.Add("time", fmt.Sprintf("%d", time))
	}

	if icao24 != "" {
		requestParams.Add("icao24", icao24)
	}

	return requestParams
}
