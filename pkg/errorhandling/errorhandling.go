package errorhandling

// code has been imported from
// https://raw.githubusercontent.com/containers/podman/main/pkg/errorhandling/errorhandling.go.

import (
	"errors"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
)

// JoinErrors converts the error slice into a single human-readable error.
func JoinErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	// If there's just one error, return it.  This prevents the "%d errors
	// occurred:" header plus list from the multierror package.
	if len(errs) == 1 {
		return errs[0]
	}

	// `multierror` appends new lines which we need to remove to prevent
	// blank lines when printing the error.
	var multiE *multierror.Error
	multiE = multierror.Append(multiE, errs...)

	finalErr := multiE.ErrorOrNil()
	if finalErr == nil {
		return nil
	}

	return errors.New(strings.TrimSpace(finalErr.Error())) //nolint:goerr113
}

// ErrorsToString converts the slice of errors into a slice of corresponding
// error messages.
func ErrorsToStrings(errs []error) []string {
	if len(errs) == 0 {
		return nil
	}

	strErrs := make([]string, len(errs))
	for i := range errs {
		strErrs[i] = errs[i].Error()
	}

	return strErrs
}

// StringsToErrors converts a slice of error messages into a slice of
// corresponding errors.
func StringsToErrors(strErrs []string) []error {
	if len(strErrs) == 0 {
		return nil
	}

	errs := make([]error, len(strErrs))
	for i := range strErrs {
		errs[i] = errors.New(strErrs[i]) //nolint:goerr113
	}

	return errs
}

// Contains checks if err's message contains sub's message. Contains should be
// used iff either err or sub has lost type information (e.g., due to
// marshalling).  For typed errors, please use `errors.Contains(...)` or `Is()`
// in recent version of Go.
func Contains(err error, sub error) bool {
	return strings.Contains(err.Error(), sub.Error())
}

// ModelError is used in remote connections with podman.
type ModelError struct {
	// API root cause formatted for automated parsing
	// example: API root cause
	Because string `json:"cause"`
	// human error message, formatted for a human to read
	// example: human error message
	Message string `json:"message"`
	// HTTP response code
	// min: 400
	ResponseCode int `json:"response"`
}

func (e ModelError) Error() string {
	return e.Message
}

func (e ModelError) Cause() error {
	return errors.New(e.Because) //nolint:goerr113
}

func (e ModelError) Code() int {
	return e.ResponseCode
}

// Cause returns the most underlying error for the provided one. There is a
// maximum error depth of 100 to avoid endless loops. An additional error log
// message will be created if this maximum has reached.
func Cause(err error) error {
	cause := err

	const maxDepth = 100
	for i := 0; i <= maxDepth; i++ {
		res := errors.Unwrap(cause)
		if res == nil {
			return cause
		}

		cause = res
	}

	log.Error().Msgf("Max error depth of %d reached, cannot unwrap until root cause: %v", maxDepth, err)

	return cause
}
