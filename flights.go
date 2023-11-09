package gopensky

import (
	"context"
	"fmt"
	"net/url"
)

// GetArrivalsByAirport retrieves flights for a certain airport which arrived within a given time interval [being, end].
// The given time interval must not be larger than seven days!
func GetArrivalsByAirport(ctx context.Context, airport string, begin int64, end int64) ([]FlighData, error) {
	var flighDataList []FlighData

	conn, err := getClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("client: %w", err)
	}

	if airport == "" {
		return nil, ErrEmptyAirportName
	}

	if begin <= 0 || end <= 0 {
		return nil, ErrInvalidUnixTime
	}

	requestParams := getFlightsRequestParams(airport, begin, end)

	response, err := conn.doGetRequest(ctx, "/flights/arrival", requestParams)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer response.Body.Close()

	if err := response.process(&flighDataList); err != nil {
		return nil, err
	}

	return flighDataList, nil
}

// GetDeparturesByAirport retrieves flights for a certain airport which departed
// within a given time interval [being, end].
// The given time interval must not be larger than seven days!
func GetDeparturesByAirport(ctx context.Context, airport string, begin int64, end int64) ([]FlighData, error) {
	var flighDataList []FlighData

	conn, err := getClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("client: %w", err)
	}

	if airport == "" {
		return nil, ErrEmptyAirportName
	}

	if begin <= 0 || end <= 0 {
		return nil, ErrInvalidUnixTime
	}

	requestParams := getFlightsRequestParams(airport, begin, end)

	response, err := conn.doGetRequest(ctx, "/flights/departure", requestParams)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer response.Body.Close()

	if err := response.process(&flighDataList); err != nil {
		return nil, err
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

	requestParams := getFlightsRequestParams("", begin, end)

	response, err := conn.doGetRequest(ctx, "/flights/all", requestParams)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer response.Body.Close()

	if err := response.process(&flighDataList); err != nil {
		return nil, err
	}

	return flighDataList, nil
}

func getFlightsRequestParams(airport string, begin int64, end int64) url.Values {
	requestParams := make(url.Values)
	if begin >= 0 {
		requestParams.Add("begin", fmt.Sprintf("%d", begin))
	}

	if end >= 0 {
		requestParams.Add("end", fmt.Sprintf("%d", end))
	}

	if airport != "" {
		requestParams.Add("airport", airport)
	}

	return requestParams
}
