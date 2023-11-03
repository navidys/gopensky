package gopensky

import (
	"errors"
)

const (
	stateVecIaco24Index = 0 + iota
	stateVecCallsignIndex
	stateVecCountryIndex
	stateVecTimePositionIndex
	stateVecLastContactIndex
	stateVecLongitudeIndex
	stateVecLatitudeIndex
	stateVecBaroAltitudeIndex
	stateVecOnGroundIndex
	stateVecVelocityIndex
	stateVecTrueTrackIndex
	stateVecVerticalRateIndex
	stateVecSensorsIndex
	stateVecGeoAltitudeIndex
	stateVecSquawkIndex
	stateVecSpiIndex
	stateVecPositionSourceIndex
	stateVecCategoryIndex
)

var (
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
)

type StatesResponse struct {
	// The time which the state vectors in this response are associated with.
	// All vectors represent the state of a vehicle with the interval.
	Time int64 `json:"time"`

	// The state vectors.
	States [][]interface{} `json:"states"`
}

type States struct {
	// The time which the state vectors in this response are associated with.
	// All vectors represent the state of a vehicle with the interval.
	Time int64 `json:"time"`

	// The state vectors.
	// States []StateVector `json:"states"`
	States []StateVector `json:"states"`
}

type StateVector struct {
	// Unique ICAO 24-bit address of the transponder in hex string representation.
	Icao24 string `json:"icao24"`

	// Callsign of the vehicle (8 chars). Can be null if no callsign has been received.
	Callsign *string `json:"callsign"`

	// Country name inferred from the ICAO 24-bit address.
	OriginCountry string `json:"origin_country"`

	// Unix timestamp (seconds) for the last position update.
	// Can be null if no position report was received by OpenSky within the past 15s.
	TimePosition *int64 `json:"time_position"`

	// Unix timestamp (seconds) for the last update in general.
	//  This field is updated for any new, valid message received from the transponder.
	LastContact int64 `json:"last_contact"`

	// WGS-84 longitude in decimal degrees. Can be null.
	Longitude *float64 `json:"longitude"`

	// WGS-84 latitude in decimal degrees. Can be null.
	Latitude *float64 `json:"latitude"`

	// Barometric altitude in meters. Can be null.
	BaroAltitude *float64 `json:"baro_altitude"`

	// Boolean value which indicates if the position was retrieved from a surface position report.
	OnGround bool `json:"on_ground"`

	// Velocity over ground in m/s. Can be null.
	Velocity *float64 `json:"velocity"`

	// True track in decimal degrees clockwise from north (north=0°). Can be null.
	TrueTrack *float64 `json:"true_track"`

	// Vertical rate in m/s.
	// A positive value indicates that the airplane is climbing, a negative value indicates that it descends.
	// Can be null.
	VerticalRate *float64 `json:"vertical_rate"`

	// IDs of the receivers which contributed to this state vector.
	//  Is null if no filtering for sensor was used in the request.
	Sensors []int `json:"sensors"`

	// Geometric altitude in meters. Can be null.
	GeoAltitude *float64 `json:"geo_altitude"`

	// The transponder code aka Squawk. Can be null.
	Squawk *string `json:"squawk"`

	// Whether flight status indicates special purpose indicator.
	Spi bool `json:"spi"`

	// Origin of this state’s position.
	// 0 = ADS-B
	// 1 = ASTERIX
	// 2 = MLAT
	// 3 = FLARM
	PositionSource int `json:"position_source"`

	// Aircraft category.
	// 0 = No information at all
	// 1 = No ADS-B Emitter Category Information
	// 2 = Light (< 15500 lbs)
	// 3 = Small (15500 to 75000 lbs)
	// 4 = Large (75000 to 300000 lbs)
	// 5 = High Vortex Large (aircraft such as B-757)
	// 6 = Heavy (> 300000 lbs)
	// 7 = High Performance (> 5g acceleration and 400 kts)
	// 8 = Rotorcraft
	// 9 = Glider / sailplane
	// 10 = Lighter-than-air
	// 11 = Parachutist / Skydiver
	// 12 = Ultralight / hang-glider / paraglider
	// 13 = Reserved
	// 14 = Unmanned Aerial Vehicle
	// 15 = Space / Trans-atmospheric vehicle
	// 16 = Surface Vehicle – Emergency Vehicle
	// 17 = Surface Vehicle – Service Vehicle
	// 18 = Point Obstacle (includes tethered balloons)
	// 19 = Cluster Obstacle
	// 20 = Line Obstacle
	Category int `json:"category"`
}

type BoundingBoxOptions struct {
	// Lower bound for the latitude in decimal degrees.
	Lamin float64

	// lower bound for the longitude in decimal degrees.
	Lomin float64

	// upper bound for the latitude in decimal degrees.
	Lamax float64

	// upper bound for the longitude in decimal degrees.
	Lomax float64
}
