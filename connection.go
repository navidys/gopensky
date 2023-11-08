package gopensky

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
)

const (
	openSkyAPIURL = "https://opensky-network.org:443/api"
	clientKey     = valueKey("Client")
)

type valueKey string

type apiResponse struct {
	*http.Response
	Request *http.Request
}

type Connection struct {
	auth   string
	uri    *url.URL
	client *http.Client
}

func newConnectionError(err error) error {
	return connectionError{err: err}
}

func NewConnection(ctx context.Context, username string, password string) (context.Context, error) {
	_url, err := url.Parse(openSkyAPIURL)
	if err != nil {
		perr := fmt.Errorf("invalid url %s: %w", openSkyAPIURL, err)

		return nil, newConnectionError(perr)
	}

	connection := Connection{
		uri: _url,
	}

	if username != "" {
		connection.auth = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	}

	dialContext := func(ctx context.Context, _, _ string) (net.Conn, error) {
		return net.Dial("tcp", _url.Host) //nolint:wrapcheck
	}

	connection.client = &http.Client{
		Transport: &http.Transport{
			DialContext:        dialContext,
			DisableCompression: true,
		},
	}

	return context.WithValue(ctx, clientKey, &connection), nil
}

// getClient from context build by NewConnection().
func getClient(ctx context.Context) (*Connection, error) {
	if c, ok := ctx.Value(clientKey).(*Connection); ok {
		return c, nil
	}

	return nil, fmt.Errorf("%w %s", errContextKey, clientKey)
}

func (c *Connection) doGetRequest(ctx context.Context, httpBody io.Reader,
	endpoint string, queryParams url.Values,
) (*apiResponse, error) {
	requestURL := fmt.Sprintf("%s/%s", c.uri, endpoint)

	if len(queryParams) > 0 {
		params := queryParams.Encode()
		requestURL = fmt.Sprintf("%s?%s", requestURL, params)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, httpBody)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	if c.auth != "" {
		req.Header.Add("Authorization", "Basic "+c.auth)
	}

	response, err := c.client.Do(req) //nolint:bodyclose

	return &apiResponse{response, req}, err //nolint:wrapcheck
}

func (h *apiResponse) isInformational() bool {
	return h.Response.StatusCode/100 == 1
}

func (h *apiResponse) isSuccess() bool {
	return h.Response.StatusCode/100 == 2 //nolint:gomnd
}

func (h *apiResponse) isRedirection() bool {
	return h.Response.StatusCode/100 == 3 //nolint:gomnd
}

// process drains the response body, and processes the HTTP status code
// Note: Closing the response.Body is left to the caller.
func (h apiResponse) process(unmarshalInto interface{}) error {
	return h.processWithError(unmarshalInto)
}

// processWithError drains the response body, and processes the HTTP status code
// Note: Closing the response.Body is left to the caller.
func (h apiResponse) processWithError(unmarshalInto interface{}) error {
	data, err := io.ReadAll(h.Response.Body)
	if err != nil {
		return fmt.Errorf("unable to process API response: %w", err)
	}

	if h.isSuccess() || h.isRedirection() {
		if unmarshalInto != nil {
			if err := json.Unmarshal(data, unmarshalInto); err != nil {
				return fmt.Errorf("unmarshalling into %#v, data %q: %w", unmarshalInto, string(data), err)
			}

			return nil
		}

		return nil
	}

	if h.isInformational() {
		return nil
	}

	return handleError(h.Response.StatusCode, data)
}

func handleError(statusCode int, data []byte) error {
	errorModel := httpModelError{
		Message: fmt.Sprintf("%s %s", http.StatusText(statusCode), data),
	}

	return errorModel
}
