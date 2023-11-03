package gopensky

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/navidys/gopensky/pkg/errorhandling"
	"github.com/rs/zerolog/log"
)

const (
	openSkyAPIURL = "https://opensky-network.org:443/api"
	clientKey     = valueKey("Client")
)

var (
	errDialContext = errors.New("unable to get dial context")
	errContextKey  = errors.New("invalid context key ")
)

type valueKey string

type APIResponse struct {
	*http.Response
	Request *http.Request
}

type Connection struct {
	auth   string
	uri    *url.URL
	client *http.Client
}

type ConnectionError struct {
	Err error
}

func (c ConnectionError) Error() string {
	return "unable to connect to api: " + c.Err.Error()
}

func (c ConnectionError) Unwrap() error {
	return c.Err
}

func newConnectionError(err error) error {
	return ConnectionError{Err: err}
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

// GetClient from context build by NewConnection().
func GetClient(ctx context.Context) (*Connection, error) {
	if c, ok := ctx.Value(clientKey).(*Connection); ok {
		return c, nil
	}

	return nil, fmt.Errorf("%w %s", errContextKey, clientKey)
}

// GetDialer returns raw Transport.DialContext from client.
func (c *Connection) GetDialer(ctx context.Context) (net.Conn, error) {
	client := c.client
	transport := client.Transport.(*http.Transport) //nolint:forcetypeassert

	if transport.DialContext != nil && transport.TLSClientConfig == nil {
		return transport.DialContext(ctx, c.uri.Scheme, c.uri.String()) //nolint:wrapcheck
	}

	return nil, errDialContext
}

func (c *Connection) DoGetRequest(ctx context.Context, httpBody io.Reader,
	endpoint string, queryParams url.Values,
) (*APIResponse, error) {
	requestURL := fmt.Sprintf("%s/%s", c.uri, endpoint)

	if len(queryParams) > 0 {
		params := queryParams.Encode()
		requestURL = fmt.Sprintf("%s?%s", requestURL, params)
		log.Debug().Msgf("do request params: %s", params)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, httpBody)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	if c.auth != "" {
		log.Debug().Msg("setting authorization")
		req.Header.Add("Authorization", "Basic "+c.auth)
	}

	log.Debug().Msgf("do get request: %s", req.URL)

	response, err := c.client.Do(req) //nolint:bodyclose

	return &APIResponse{response, req}, err //nolint:wrapcheck
}

func (h *APIResponse) IsInformational() bool {
	return h.Response.StatusCode/100 == 1
}

func (h *APIResponse) IsSuccess() bool {
	return h.Response.StatusCode/100 == 2 //nolint:gomnd
}

func (h *APIResponse) IsRedirection() bool {
	return h.Response.StatusCode/100 == 3 //nolint:gomnd
}

func (h *APIResponse) IsConflictError() bool {
	return h.Response.StatusCode == 409 //nolint:gomnd
}

func (h *APIResponse) IsClientError() bool {
	return h.Response.StatusCode == 409 //nolint:gomnd
}

func (h *APIResponse) IsServerError() bool {
	return h.Response.StatusCode/100 == 5 //nolint:gomnd
}

// Process drains the response body, and processes the HTTP status code
// Note: Closing the response.Body is left to the caller.
func (h APIResponse) Process(unmarshalInto interface{}) error {
	return h.ProcessWithError(unmarshalInto, &errorhandling.ModelError{})
}

// ProcessWithError drains the response body, and processes the HTTP status code
// Note: Closing the response.Body is left to the caller.
func (h APIResponse) ProcessWithError(unmarshalInto interface{}, unmarshalErrorInto interface{}) error {
	data, err := io.ReadAll(h.Response.Body)
	if err != nil {
		return fmt.Errorf("unable to process API response: %w", err)
	}

	if h.IsSuccess() || h.IsRedirection() {
		if unmarshalInto != nil {
			if err := json.Unmarshal(data, unmarshalInto); err != nil {
				return fmt.Errorf("unmarshalling into %#v, data %q: %w", unmarshalInto, string(data), err)
			}

			return nil
		}

		return nil
	}

	if h.IsConflictError() {
		return handleError(data, unmarshalErrorInto)
	}

	return handleError(data, &errorhandling.ModelError{})
}

func handleError(data []byte, unmarshalErrorInto interface{}) error {
	if err := json.Unmarshal(data, unmarshalErrorInto); err != nil {
		return fmt.Errorf("unmarshalling error into %#v, data %q: %w", unmarshalErrorInto, string(data), err)
	}

	return unmarshalErrorInto.(error) //nolint:forcetypeassert
}
