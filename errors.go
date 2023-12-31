package gopensky

import (
	"errors"
)

var (
	errContextKey = errors.New("invalid context key")

	errStateVecDataCount      = errors.New("invalid state vector data count")
	errStateVecIcao24         = errors.New("state vector icao24 assertion failed")
	errStateVecCallsign       = errors.New("state vector callsign assertion failed")
	errStateVecOriginCountry  = errors.New("state vector origin country assertion failed")
	errStateVecTimePosition   = errors.New("state vector time position assertion failed")
	errStateVecLastContact    = errors.New("state vector last contact assertion failed")
	errStateVecLongitude      = errors.New("state vector longitude assertion failed")
	errStateVecLatitude       = errors.New("state vector latitude assertion failed")
	errStateVecBaroAltitude   = errors.New("state vector baro altitude assertion failed")
	errStateVecOnGround       = errors.New("state vector on ground assertion failed")
	errStateVecVelocity       = errors.New("state vector velocity assertion failed")
	errStateVecTrueTrack      = errors.New("state vector true track assertion failed")
	errStateVecVerticalRate   = errors.New("state vector vertical rate assertion failed")
	errStateVecSensors        = errors.New("state vector sensors assertion failed")
	errStateVecGeoAltitude    = errors.New("state vector geo altitude assertion failed")
	errStateVecSquawk         = errors.New("state vector squawk assertion failed")
	errStateVecSpi            = errors.New("state vector spi assertion failed")
	errStateVecPositionSource = errors.New("state vector position source assertion failed")
	errStateVecCategory       = errors.New("state vector category assertion failed")

	errWaypointsDataCount   = errors.New("invalid track waypoints data count")
	errWaypointTime         = errors.New("track waypoints time assertion failed")
	errWaypointLatitude     = errors.New("track waypoints latitude assertion failed")
	errWaypointLongitude    = errors.New("track waypoints longitude assertion failed")
	errWaypointBaroAltitude = errors.New("track waypoints baro altitude assertion failed")
	errWaypointTrueTrack    = errors.New("track waypoints true track assertion failed")
	errWaypointOnGround     = errors.New("track waypoints on ground assertion failed")

	ErrInvalidAirportName  = errors.New("invalid airport name")
	ErrInvalidAircraftName = errors.New("invalid aircraft name")
	ErrInvalidUnixTime     = errors.New("invalid unix time")
)

type connectionError struct {
	err error
}

func (c connectionError) Error() string {
	return "unable to connect to api: " + c.err.Error()
}

func (c connectionError) Unwrap() error {
	return c.err
}

type httpModelError struct {
	Message string `json:"message"`
}

func (e httpModelError) Error() string {
	return e.Message
}
