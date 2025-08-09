package gopensky

import (
	"context"
	"fmt"
	"net/url"
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

// Retrieve state vectors for a given time. If time = 0 the most recent ones are taken.
// It is possible to query a certain area defined by a bounding box of WGS84 coordinates.
// You can request the category of aircraft by setting extended to true.
func GetStates(ctx context.Context, time int64, icao24 []string,
	bBox *BoundingBoxOptions, extended bool,
) (*States, error) {
	var statesRep StatesResponse

	conn, err := getClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("client: %w", err)
	}

	if time < 0 {
		return nil, ErrInvalidUnixTime
	}

	requestParams := getStateRequestParams(time, icao24, bBox, extended)

	response, err := conn.doGetRequest(ctx, "/states/all", requestParams)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	errRespProcess := response.process(&statesRep)

	err = response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("response body close %w", err)
	}

	if errRespProcess != nil {
		return nil, errRespProcess
	}

	statesVecList := make([]StateVector, 0)

	for _, st := range statesRep.States {
		stvec, err := decodeRawStateVector(st)
		if err != nil {
			return nil, fmt.Errorf("decode state vector: %w", err)
		}

		if stvec != nil {
			statesVecList = append(statesVecList, *stvec)
		}
	}

	states := States{
		Time:   statesRep.Time,
		States: statesVecList,
	}

	return &states, nil
}

func getStateRequestParams(time int64, icao24 []string, bBox *BoundingBoxOptions, extended bool) url.Values {
	requestParams := make(url.Values)
	if time > 0 {
		requestParams.Add("time", fmt.Sprintf("%d", time)) //nolint:perfsprint
	}

	for _, icao24Item := range icao24 {
		requestParams.Add("icao24", icao24Item)
	}

	if extended {
		requestParams.Add("extended", "1")
	}

	if bBox != nil {
		requestParams.Add("lamax", floatToString(bBox.Lamax))
		requestParams.Add("lamin", floatToString(bBox.Lamin))
		requestParams.Add("lomax", floatToString(bBox.Lomax))
		requestParams.Add("lomin", floatToString(bBox.Lomin))
	}

	return requestParams
}

func decodeRawStateVector(data []interface{}) (*StateVector, error) { //nolint:funlen,cyclop,gocognit,gocyclo
	var assertionOK bool

	stVector := StateVector{}
	recvDataCount := len(data)

	if recvDataCount < stateVecCategoryIndex {
		return nil, errStateVecDataCount
	}

	// Icao24 index
	stVector.Icao24, assertionOK = data[stateVecIaco24Index].(string)
	if !assertionOK {
		return nil, fmt.Errorf("%w: %v", errStateVecIcao24, data[stateVecIaco24Index])
	}

	// Callsign index
	if data[stateVecCallsignIndex] != nil {
		val, assertionOK := data[stateVecCallsignIndex].(string)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errStateVecCallsign, data[stateVecCallsignIndex])
		}

		stVector.Callsign = &val
	}

	// OriginCountry index
	stVector.OriginCountry, assertionOK = data[stateVecCountryIndex].(string)
	if !assertionOK {
		return nil, fmt.Errorf("%w: %v", errStateVecOriginCountry, data[stateVecCountryIndex])
	}

	// TimePosition index
	if data[stateVecTimePositionIndex] != nil {
		val, assertionOK := data[stateVecTimePositionIndex].(float64)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errStateVecTimePosition, data[stateVecTimePositionIndex])
		}

		timePos := int64(val)

		stVector.TimePosition = &timePos
	}

	// LastContact index
	lastContact, assertionOK := data[stateVecLastContactIndex].(float64)
	if !assertionOK {
		return nil, fmt.Errorf("%w: %v", errStateVecLastContact, data[stateVecLastContactIndex])
	}

	stVector.LastContact = int64(lastContact)

	// Longitude index
	if data[stateVecLongitudeIndex] != nil {
		stVecLongitude, assertionOK := data[stateVecLongitudeIndex].(float64)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errStateVecLongitude, data[stateVecLongitudeIndex])
		}

		stVector.Longitude = &stVecLongitude
	}

	// Latitude index
	if data[stateVecLatitudeIndex] != nil {
		stVecLatitude, assertionOK := data[stateVecLatitudeIndex].(float64)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errStateVecLatitude, data[stateVecLatitudeIndex])
		}

		stVector.Latitude = &stVecLatitude
	}

	// BaroAltitude index
	if data[stateVecBaroAltitudeIndex] != nil {
		stVectorBaroAltitude, assertionOK := data[stateVecBaroAltitudeIndex].(float64)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errStateVecBaroAltitude, data[stateVecBaroAltitudeIndex])
		}

		stVector.BaroAltitude = &stVectorBaroAltitude
	}

	// OnGround index
	stVector.OnGround, assertionOK = data[stateVecOnGroundIndex].(bool)
	if !assertionOK {
		return nil, fmt.Errorf("%w: %v", errStateVecOnGround, data[stateVecOnGroundIndex])
	}

	// Velocity index
	if data[stateVecVelocityIndex] != nil {
		stVectorVelocity, assertionOK := data[stateVecVelocityIndex].(float64)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errStateVecVelocity, data[stateVecVelocityIndex])
		}

		stVector.Velocity = &stVectorVelocity
	}

	// TrueTrack index
	if data[stateVecTrueTrackIndex] != nil {
		stVectorTrueTrack, assertionOK := data[stateVecTrueTrackIndex].(float64)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errStateVecTrueTrack, data[stateVecTrueTrackIndex])
		}

		stVector.TrueTrack = &stVectorTrueTrack
	}

	// VerticalRate index
	if data[stateVecVerticalRateIndex] != nil {
		stVectorVerticalRate, assertionOK := data[stateVecVerticalRateIndex].(float64)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errStateVecVerticalRate, data[stateVecVerticalRateIndex])
		}

		stVector.VerticalRate = &stVectorVerticalRate
	}

	// Sensors index
	if data[stateVecSensorsIndex] != nil {
		stVector.Sensors, assertionOK = data[stateVecSensorsIndex].([]int)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errStateVecSensors, data[stateVecSensorsIndex])
		}
	}

	// GeoAltitude index
	if data[stateVecGeoAltitudeIndex] != nil {
		stVectorGeoAltitude, assertionOK := data[stateVecGeoAltitudeIndex].(float64)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errStateVecGeoAltitude, data[stateVecGeoAltitudeIndex])
		}

		stVector.GeoAltitude = &stVectorGeoAltitude
	}

	// Squawk index
	if data[stateVecSquawkIndex] != nil {
		stVectorSquawk, assertionOK := data[stateVecSquawkIndex].(string)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errStateVecSquawk, data[stateVecSquawkIndex])
		}

		stVector.Squawk = &stVectorSquawk
	}

	// Spi index
	stVector.Spi, assertionOK = data[stateVecSpiIndex].(bool)
	if !assertionOK {
		return nil, fmt.Errorf("%w: %v", errStateVecSpi, data[stateVecSpiIndex])
	}

	// PositionSource index
	stVectorPositionSource, assertionOK := data[stateVecPositionSourceIndex].(float64)
	if !assertionOK {
		return nil, fmt.Errorf("%w: %v", errStateVecPositionSource, data[stateVecPositionSourceIndex])
	}

	stVector.PositionSource = int(stVectorPositionSource)

	// Category index
	if recvDataCount == stateVecCategoryIndex+1 {
		stVectorCategory, assertionOK := data[stateVecCategoryIndex].(float64)
		if !assertionOK {
			return nil, fmt.Errorf("%w: %v", errStateVecCategory, data[stateVecCategoryIndex])
		}

		stVector.Category = int(stVectorCategory)
	}

	return &stVector, nil
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
