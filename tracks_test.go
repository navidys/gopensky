package gopensky_test

import (
	"context"
	"fmt"

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

	Describe("getTracksRequestParams", func() {
		It("generate track request parameters", func() {
			tests := []struct {
				time   int64
				icao24 string
			}{
				{time: -1, icao24: "icao24_a"},
				{time: -1, icao24: ""},
				{time: 0, icao24: "icao24_b"},
				{time: 2, icao24: "icao24_c"},
			}

			for _, reqParams := range tests {
				urlVal := gopensky.GetTracksRequestParams(reqParams.time, reqParams.icao24)
				Expect(urlVal.Get("icao24")).To(Equal(reqParams.icao24))

				if reqParams.time >= 0 {
					reqTime := fmt.Sprintf("%d", reqParams.time)
					Expect(urlVal.Get("time")).To(Equal(reqTime))
				} else {
					Expect(urlVal.Get("time")).To(Equal(""))
				}
			}
		})
	})
})
