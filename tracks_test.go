package gopensky_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/navidys/gopensky"
)

var _ = Describe("Tracks", func() {
	Describe("GetTrackByAircraft", func() {
		It("retrieves the trajectory for a certain aircraft at a given time", func() {
			conn, err := gopensky.NewConnection(context.Background(), "", "")
			Expect(err).NotTo(HaveOccurred())

			_, err = gopensky.GetTrackByAircraft(conn, "", 1696755342)
			Expect(err).To(Equal(gopensky.ErrInvalidAircraftName))

			_, err = gopensky.GetTrackByAircraft(conn, "a835af", -1)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))
		})
	})
})
