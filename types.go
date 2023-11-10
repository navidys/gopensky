package gopensky

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

const (
	trackTimeIndex = 0 + iota
	trackLatitudeIndex
	trackLongitudeIndex
	trackBaroAltitudeIndex
	trackTureTrackIndex
	trackOnGroundIndex
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

	// Callsign of the vehicle (8 chars). Can be nil if no callsign has been received.
	Callsign *string `json:"callsign"`

	// Country name inferred from the ICAO 24-bit address.
	OriginCountry string `json:"originCountry"`

	// Unix timestamp (seconds) for the last position update.
	// Can be nil if no position report was received by OpenSky within the past 15s.
	TimePosition *int64 `json:"timePosition"`

	// Unix timestamp (seconds) for the last update in general.
	//  This field is updated for any new, valid message received from the transponder.
	LastContact int64 `json:"lastContact"`

	// WGS-84 longitude in decimal degrees. Can be nil.
	Longitude *float64 `json:"longitude"`

	// WGS-84 latitude in decimal degrees. Can be nil.
	Latitude *float64 `json:"latitude"`

	// Barometric altitude in meters. Can be nil.
	BaroAltitude *float64 `json:"baroAltitude"`

	// Boolean value which indicates if the position was retrieved from a surface position report.
	OnGround bool `json:"onGround"`

	// Velocity over ground in m/s. Can be nil.
	Velocity *float64 `json:"velocity"`

	// True track in decimal degrees clockwise from north (north=0°). Can be nil.
	TrueTrack *float64 `json:"trueTrack"`

	// Vertical rate in m/s.
	// A positive value indicates that the airplane is climbing, a negative value indicates that it descends.
	// Can be nil.
	VerticalRate *float64 `json:"verticalTate"`

	// IDs of the receivers which contributed to this state vector.
	//  Is nil if no filtering for sensor was used in the request.
	Sensors []int `json:"sensors"`

	// Geometric altitude in meters. Can be nil.
	GeoAltitude *float64 `json:"geoAltitude"`

	// The transponder code aka Squawk. Can be nil.
	Squawk *string `json:"squawk"`

	// Whether flight status indicates special purpose indicator.
	Spi bool `json:"spi"`

	// Origin of this state’s position.
	// 0 = ADS-B
	// 1 = ASTERIX
	// 2 = MLAT
	// 3 = FLARM
	PositionSource int `json:"positionSource"`

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

// NewBoundingBox returns new bounding box options for states information gathering.
func NewBoundingBox(lamin float64, lomin float64, lamax float64, lomax float64) *BoundingBoxOptions {
	return &BoundingBoxOptions{
		Lamin: lamin,
		Lomin: lomin,
		Lamax: lamax,
		Lomax: lomax,
	}
}
