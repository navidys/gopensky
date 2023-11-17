package gopensky

import (
	"context"
	"fmt"
	"net/http"
)

var (
	DecodeRawStateVector     = decodeRawStateVector
	DecodeWaypoint           = decodeWaypoint
	FloatToString            = floatToString
	GetFlightsRequestParams  = getFlightsRequestParams
	GetTracksRequestParams   = getTracksRequestParams
	GetStateRequestParams    = getStateRequestParams
	ParseFlightTrackResponse = parseFlightTrackResponse
	OpenSkyAPIURL            = openSkyAPIURL
	NewConnectionError       = newConnectionError
	HandleError              = handleError
)

func GetClient(ctx context.Context) (*http.Client, error) {
	if c, ok := ctx.Value(clientKey).(*Connection); ok {
		return c.client, nil
	}

	return nil, fmt.Errorf("%w %s", errContextKey, clientKey)
}
