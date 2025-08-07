package gopensky

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

const (
	trackTimeIndex = 0 + iota
	trackLatitudeIndex
	trackLongitudeIndex
	trackBaroAltitudeIndex
	trackTureTrackIndex
	trackOnGroundIndex
)

type FlightTrackResponse struct {
	Icao24    string          `json:"icao24"`
	StartTime float64         `json:"startTime"`
	EndTime   float64         `json:"endTime"`
	Callsign  *string         `json:"callsign"`
	Path      [][]interface{} `json:"path"`
}

type FlightTrack struct {
	// Unique ICAO 24-bit address of the transponder in hex string representation.
	Icao24 string `json:"icao24"`

	// Time of the first waypoint in seconds since epoch (Unix time).
	StartTime int64 `json:"startTime"`

	// Time of the last waypoint in seconds since epoch (Unix time).
	EndTime int64 `json:"endTime"`

	// Callsign (8 characters) that holds for the whole track. Can be nil.
	Callsign *string `json:"callsign"`

	// Waypoints of the trajectory (description below).
	Path []WayPoint `json:"path"`
}

type WayPoint struct {
	// Time which the given waypoint is associated with in seconds since epoch (Unix time).
	Time int64 `json:"time"`

	// WGS-84 latitude in decimal degrees. Can be nil.
	Latitude *float64 `json:"latitude"`

	// WGS-84 longitude in decimal degrees. Can be nil.
	Longitude *float64 `json:"longitude"`

	// Barometric altitude in meters. Can be nil.
	BaroAltitude *float64 `json:"baroAltitude"`

	// True track in decimal degrees clockwise from north (north=0°). Can be nil.
	TrueTrack *float64 `json:"trueTrack"`

	// Boolean value which indicates if the position was retrieved from a surface position report.
	OnGround bool `json:"onGround"`
}

// GetTrackByAircraft retrieves the trajectory for a certain aircraft at a given time.
func GetTrackByAircraft(ctx context.Context, icao24 string, time int64) (FlightTrack, error) {
	var (
		flightTrack         FlightTrack
		flightTrackResponse FlightTrackResponse
	)

	if icao24 == "" {
		return flightTrack, ErrInvalidAircraftName
	}

	if time < 0 {
		return flightTrack, ErrInvalidUnixTime
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

	errRespProcess := response.process(&flightTrackResponse)

	if err := response.Body.Close(); err != nil {
		return flightTrack, fmt.Errorf("response body close %w", err)
	}

	if errRespProcess != nil {
		return flightTrack, errRespProcess
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
	startTime := time.Unix(int64(response.StartTime), 0).Unix()
	if startTime <= 0 {
		flightTrack.StartTime = 1
	} else {
		flightTrack.StartTime = startTime
	}

	for _, waypointData := range response.Path {
		waypoint, err := decodeWaypoint(waypointData)
		if err != nil {
			return flightTrack, fmt.Errorf("decode waypoint: %w", err)
		}

		flightTrack.Path = append(flightTrack.Path, *waypoint)
	}

	return flightTrack, nil
}

func decodeWaypoint(data []interface{}) (*WayPoint, error) { //nolint:funlen,cyclop
	if len(data) < trackOnGroundIndex {
		return nil, errWaypointsDataCount
	}

	var waypoint WayPoint

	// Time index
	if data[trackTimeIndex] != nil {
		wtime, assertionOK := data[trackTimeIndex].(int64)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errWaypointTime, data[trackTimeIndex])
		}

		waypoint.Time = wtime
	}

	// Latitude index
	if data[trackLatitudeIndex] != nil {
		latitude, assertionOK := data[trackLatitudeIndex].(float64)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errWaypointLatitude, data[trackLatitudeIndex])
		}

		waypoint.Latitude = &latitude
	}

	// Longitude index
	if data[trackLongitudeIndex] != nil {
		longitude, assertionOK := data[trackLongitudeIndex].(float64)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errWaypointLongitude, data[trackLongitudeIndex])
		}

		waypoint.Longitude = &longitude
	}

	// BaroAltitude index
	if data[trackBaroAltitudeIndex] != nil {
		baroAltitude, assertionOK := data[trackBaroAltitudeIndex].(float64)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errWaypointBaroAltitude, data[trackBaroAltitudeIndex])
		}

		waypoint.BaroAltitude = &baroAltitude
	}

	// TrueTrack index
	if data[trackTureTrackIndex] != nil {
		trueTrack, assertionOK := data[trackTureTrackIndex].(float64)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errWaypointTrueTrack, data[trackTureTrackIndex])
		}

		waypoint.TrueTrack = &trueTrack
	}

	// Onground index
	if data[trackOnGroundIndex] != nil {
		onGround, assertionOK := data[trackOnGroundIndex].(bool)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errWaypointOnGround, data[trackOnGroundIndex])
		}

		waypoint.OnGround = onGround
	}

	return &waypoint, nil
}

func getTracksRequestParams(time int64, icao24 string) url.Values {
	requestParams := make(url.Values)
	if time >= 0 {
		requestParams.Add("time", fmt.Sprintf("%d", time)) //nolint:perfsprint
	}

	if icao24 != "" {
		requestParams.Add("icao24", icao24)
	}

	return requestParams
}
