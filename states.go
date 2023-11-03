package gopensky

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rs/zerolog/log"
)

// Retrieve state vectors for a given time. If time = 0 the most recent ones are taken.
func GetStates(ctx context.Context, time int64, icao24 []string, bBox *BoundingBoxOptions) (*States, error) {
	var statesRep StatesResponse

	conn, err := GetClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("client: %w", err)
	}

	requestParams := make(url.Values)
	if time >= 0 {
		requestParams.Add("time", fmt.Sprintf("%d", time))
	}

	for _, icao24Item := range icao24 {
		requestParams.Add("icao24", icao24Item)
	}

	response, err := conn.DoGetRequest(ctx, nil, "/states/all", requestParams)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer response.Body.Close()

	if err := response.Process(&statesRep); err != nil {
		return nil, err
	}

	statesVecList := make([]StateVector, 0)

	for _, st := range statesRep.States {
		stvec, err := decodeRawStateVector(st)
		if err != nil {
			log.Error().Msgf("cannot decode received data: %v", err)

			continue
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
