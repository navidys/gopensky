package gopensky_test

import (
	"errors"
	"net/http"

	"github.com/navidys/gopensky"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Errors", func() {
	Describe("connectionError ", func() {
		It("Error()", func() {
			connErr := gopensky.NewConnectionError(errors.New("test error"))
			Expect(connErr.Error()).To(Equal("unable to connect to api: test error"))
		})
	})

	Describe("httpModelError ", func() {
		It("Error()", func() {
			httpError := gopensky.HandleError(http.StatusNotFound, []byte("test data"))
			Expect(httpError.Error()).To(Equal(http.StatusText(http.StatusNotFound) + " test data"))
		})
	})
})
