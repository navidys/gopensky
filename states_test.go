package gopensky_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/navidys/gopensky"
)

var _ = Describe("States", func() {
	Describe("GetStates", func() {
		It("retrieve state vectors for a given time", func() {
			conn, err := gopensky.NewConnection(context.Background(), "", "")
			Expect(err).NotTo(HaveOccurred())

			_, err = gopensky.GetStates(conn, -1, nil, nil, false)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))
		})
	})

	Describe("getStateRequestParams", func() {
		It("generate state request parameters", func() {
			tests := []struct {
				time     int64
				icao24   []string
				bBox     *gopensky.BoundingBoxOptions
				extended bool
			}{
				{time: -1, icao24: nil, bBox: nil, extended: false},
				{time: 0, icao24: nil, bBox: nil, extended: false},
				{time: 0, icao24: []string{"icao24_a", "icao24_b"}, bBox: gopensky.NewBoundingBox(1.1, 1.2, 1, 1), extended: false},
				{time: 0, icao24: []string{"icao24_a", "icao24_b"}, bBox: gopensky.NewBoundingBox(2.2111, 2.1, 0, 1), extended: true},
			}

			for _, reqParams := range tests {
				urlVal := gopensky.GetStateRequestParams(reqParams.time, reqParams.icao24, reqParams.bBox, reqParams.extended)

				if reqParams.time > 0 {
					reqtime := fmt.Sprintf("%d", reqParams.time)
					Expect(urlVal.Get("time")).To(Equal(reqtime))
				}

				if len(reqParams.icao24) > 0 {
					for _, iaco24 := range reqParams.icao24 {
						Expect(urlVal["icao24"]).Should(ContainElement(ContainSubstring(iaco24)))
					}
				}

				if reqParams.bBox != nil {
					lamax := gopensky.FloatToString(reqParams.bBox.Lamax)
					lamin := gopensky.FloatToString(reqParams.bBox.Lamin)
					lomax := gopensky.FloatToString(reqParams.bBox.Lomax)
					lomin := gopensky.FloatToString(reqParams.bBox.Lomin)

					Expect(urlVal.Get("lamax")).To(Equal(lamax))
					Expect(urlVal.Get("lamin")).To(Equal(lamin))
					Expect(urlVal.Get("lomax")).To(Equal(lomax))
					Expect(urlVal.Get("lomin")).To(Equal(lomin))
				}

				if reqParams.extended {
					Expect(urlVal.Get("extended")).To(Equal("1"))
				}
			}
		})
	})
})
