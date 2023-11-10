package gopensky

import (
	"context"
	"fmt"
	"net/url"
	"time"
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

	flightTrack.Icao24 = response.Icao24
	flightTrack.Callsign = response.Callsign
	flightTrack.EndTime = time.Unix(int64(response.EndTime), 0).Unix()
	// the api is not returning proper start time value
	// temporary checking if its <= 0 then allocated 1
	if flightTrack.StartTime <= 0 {
		flightTrack.StartTime = 1
	}

	for _, data := range response.Path {
		if len(data) < trackOnGroundIndex {
			return flightTrack, errWaypointsDataCount
		}

		var waypoint WayPoint

		// BaroAltitude index
		if data[trackBaroAltitudeIndex] != nil {
			baroAltitude, assertionOK := data[trackBaroAltitudeIndex].(float64)
			if !assertionOK {
				return flightTrack, fmt.Errorf("%w: %v", errWaypointBaroAltitude, data[trackBaroAltitudeIndex])
			}

			waypoint.BaroAltitude = &baroAltitude
		}

		// TrueTrack index
		if data[trackTureTrackIndex] != nil {
			trueTrack, assertionOK := data[trackTureTrackIndex].(float64)
			if !assertionOK {
				return flightTrack, fmt.Errorf("%w: %v", errWaypointTrueTrack, data[trackTureTrackIndex])
			}

			waypoint.TrueTrack = &trueTrack
		}

		// Onground index
		if data[trackOnGroundIndex] != nil {
			onGround, assertionOK := data[trackOnGroundIndex].(bool)
			if !assertionOK {
				return flightTrack, fmt.Errorf("%w: %v", errWaypointOnGround, data[trackOnGroundIndex])
			}

			waypoint.OnGround = onGround
		}

		flightTrack.Path = append(flightTrack.Path, waypoint)
	}

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
