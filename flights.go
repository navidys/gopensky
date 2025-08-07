package gopensky

import (
	"context"
	"fmt"
	"net/url"
)

type FlighData struct {
	// Unique ICAO 24-bit address of the transponder in hex string representation.
	// All letters are lower case.
	Icao24 string `json:"icao24"`

	// Estimated time of departure for the flight as Unix time (seconds since epoch).
	FirstSeen int64 `json:"firstSeen"`

	// ICAO code of the estimated departure airport.
	// Can be nil if the airport could not be identified.
	EstDepartureAirport *string `json:"estDepartureAirport"`

	// Estimated time of arrival for the flight as Unix time (seconds since epoch).
	LastSeen int64 `json:"lastSeen"`

	// ICAO code of the estimated arrival airport.
	// Can be nil if the airport could not be identified.
	EstArrivalAirport *string `json:"estArrivalAirport"`

	// Callsign of the vehicle (8 chars).
	// Can be nil if no callsign has been received.
	// If the vehicle transmits multiple callsigns during the flight,
	// we take the one seen most frequently.
	Callsign *string `json:"callsign"`

	// Horizontal distance of the last received airborne position
	// to the estimated departure airport in meters.
	EstDepartureAirportHorizDistance int64 `json:"estDepartureAirportHorizDistance"`

	// Vertical distance of the last received airborne position
	// to the estimated departure airport in meters.
	EstDepartureAirportVertDistance int64 `json:"estDepartureAirportVertDistance"`

	// Horizontal distance of the last received airborne position
	// to the estimated arrival airport in meters.
	EstArrivalAirportHorizDistance int64 `json:"estArrivalAirportHorizDistance"`

	// Vertical distance of the last received airborne position to
	// the estimated arrival airport in meters.
	EstArrivalAirportVertDistance int64 `json:"estArrivalAirportVertDistance"`

	// Number of other possible departure airports.
	// These are airports in short distance to estDepartureAirport.
	DepartureAirportCandidatesCount int `json:"departureAirportCandidatesCount"`

	// Number of other possible departure airports.
	ArrivalAirportCandidatesCount int `json:"arrivalAirportCandidatesCount"`
}

// GetArrivalsByAirport retrieves flights for a certain airport which arrived within a given time interval [being, end].
// The given time interval must not be larger than seven days!
func GetArrivalsByAirport(ctx context.Context, airport string, begin int64, end int64) ([]FlighData, error) {
	var flighDataList []FlighData

	conn, err := getClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("client: %w", err)
	}

	if airport == "" {
		return nil, ErrInvalidAirportName
	}

	if begin <= 0 || end <= 0 {
		return nil, ErrInvalidUnixTime
	}

	requestParams := getFlightsRequestParams(airport, "", begin, end)

	response, err := conn.doGetRequest(ctx, "/flights/arrival", requestParams)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	errRespProcess := response.process(&flighDataList)

	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("response body close %w", err)
	}

	if errRespProcess != nil {
		return nil, errRespProcess
	}

	return flighDataList, nil
}

// GetDeparturesByAirport retrieves flights for a certain airport which departed
// The given time interval must not be larger than seven days!
func GetDeparturesByAirport(ctx context.Context, airport string, begin int64, end int64) ([]FlighData, error) {
	var flighDataList []FlighData

	conn, err := getClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("client: %w", err)
	}

	if airport == "" {
		return nil, ErrInvalidAirportName
	}

	if begin <= 0 || end <= 0 {
		return nil, ErrInvalidUnixTime
	}

	requestParams := getFlightsRequestParams(airport, "", begin, end)

	response, err := conn.doGetRequest(ctx, "/flights/departure", requestParams)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	errRespProcess := response.process(&flighDataList)

	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("response body close %w", err)
	}

	if errRespProcess != nil {
		return nil, errRespProcess
	}

	return flighDataList, nil
}

// GetFlightsByInterval retrieves flights for a certain time interval [begin, end].
// The given time interval must not be larger than two hours!
func GetFlightsByInterval(ctx context.Context, begin int64, end int64) ([]FlighData, error) {
	var flighDataList []FlighData

	conn, err := getClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("client: %w", err)
	}

	if begin <= 0 || end <= 0 {
		return nil, ErrInvalidUnixTime
	}

	requestParams := getFlightsRequestParams("", "", begin, end)

	response, err := conn.doGetRequest(ctx, "/flights/all", requestParams)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	errRespProcess := response.process(&flighDataList)

	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("response body close %w", err)
	}

	if errRespProcess != nil {
		return nil, errRespProcess
	}

	return flighDataList, nil
}

// GetFlightsByAircraft retrieves flights for a particular aircraft within a certain time interval.
// Resulting flights departed and arrived within [begin, end].
// The given time interval must not be larger than 30 days!
func GetFlightsByAircraft(ctx context.Context, icao24 string, begin int64, end int64) ([]FlighData, error) {
	var flighDataList []FlighData

	conn, err := getClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("client: %w", err)
	}

	if icao24 == "" {
		return nil, ErrInvalidAircraftName
	}

	if begin <= 0 || end <= 0 {
		return nil, ErrInvalidUnixTime
	}

	requestParams := getFlightsRequestParams("", icao24, begin, end)

	response, err := conn.doGetRequest(ctx, "/flights/aircraft", requestParams)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	errRespProcess := response.process(&flighDataList)

	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("response body close %w", err)
	}

	if errRespProcess != nil {
		return nil, errRespProcess
	}

	return flighDataList, nil
}

func getFlightsRequestParams(airport string, aircraft string, begin int64, end int64) url.Values {
	requestParams := make(url.Values)
	if begin >= 0 {
		requestParams.Add("begin", fmt.Sprintf("%d", begin)) //nolint:perfsprint
	}

	if end >= 0 {
		requestParams.Add("end", fmt.Sprintf("%d", end)) //nolint:perfsprint
	}

	if aircraft != "" {
		requestParams.Add("icao24", aircraft)
	}

	if airport != "" {
		requestParams.Add("airport", airport)
	}

	return requestParams
}
