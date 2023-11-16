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
				{time: 2, icao24: []string{"icao24_a", "icao24_b"}, bBox: gopensky.NewBoundingBox(1.1, 1.2, 1, 1), extended: false},
				{time: 2, icao24: []string{"icao24_a", "icao24_b"}, bBox: gopensky.NewBoundingBox(2.2111, 2.1, 0, 1), extended: true},
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

	Describe("decodeRawStateVector", func() {
		It("decode state vector array interface to struct", func() {
			testStringValue := "test_str"
			testIntValue := 1
			testTimeValue := int64(1)
			testFloatValue := -46.509700
			sensors := []int{1, 2}

			tests := []struct {
				have  []interface{}
				wants gopensky.StateVector
			}{
				{
					have:  []interface{}{testStringValue},
					wants: gopensky.StateVector{}},
				{
					have: []interface{}{
						testStringValue,
						testStringValue,
						testStringValue,
						float64(testTimeValue),
						float64(testTimeValue),
						testFloatValue,
						testFloatValue,
						testFloatValue,
						false,
						testFloatValue,
						testFloatValue,
						testFloatValue,
						sensors,
						testFloatValue,
						testStringValue,
						false,
						float64(testIntValue),
						float64(testIntValue),
					},
					wants: gopensky.StateVector{
						Icao24:         testStringValue,
						Callsign:       &testStringValue,
						OriginCountry:  testStringValue,
						TimePosition:   &testTimeValue,
						LastContact:    testTimeValue,
						Longitude:      &testFloatValue,
						Latitude:       &testFloatValue,
						BaroAltitude:   &testFloatValue,
						OnGround:       false,
						Velocity:       &testFloatValue,
						TrueTrack:      &testFloatValue,
						VerticalRate:   &testFloatValue,
						Sensors:        sensors,
						GeoAltitude:    &testFloatValue,
						Squawk:         &testStringValue,
						Spi:            false,
						PositionSource: testIntValue,
						Category:       testIntValue,
					},
				},
				{
					have: []interface{}{
						testStringValue,
						nil,
						testStringValue,
						nil,
						float64(testTimeValue),
						nil,
						nil,
						nil,
						false,
						nil,
						nil,
						nil,
						sensors,
						nil,
						nil,
						false,
						float64(testIntValue),
						float64(testIntValue),
					},
					wants: gopensky.StateVector{
						Icao24:         testStringValue,
						Callsign:       nil,
						OriginCountry:  testStringValue,
						TimePosition:   nil,
						LastContact:    testTimeValue,
						Longitude:      nil,
						Latitude:       nil,
						BaroAltitude:   nil,
						OnGround:       false,
						Velocity:       nil,
						TrueTrack:      nil,
						VerticalRate:   nil,
						Sensors:        sensors,
						GeoAltitude:    nil,
						Squawk:         nil,
						Spi:            false,
						PositionSource: testIntValue,
						Category:       testIntValue,
					},
				},
			}

			for _, decodeTest := range tests {
				stvec, err := gopensky.DecodeRawStateVector(decodeTest.have)

				if len(decodeTest.have) < 18 {
					Expect(err).To(HaveOccurred())
					Expect(stvec).To(BeNil())

					continue
				}

				Expect(err).ToNot(HaveOccurred())
				Expect(stvec).ToNot(BeNil())

				decodeData := *(stvec)

				// icao24
				Expect(decodeData.Icao24).To(Equal(decodeTest.wants.Icao24))

				// Callsign
				if decodeTest.wants.Callsign != nil {
					Expect(*(decodeData.Callsign)).To(Equal(testStringValue))
				} else {
					Expect(decodeData.Callsign).To(BeNil())
				}

				// OriginCountry
				Expect(decodeData.OriginCountry).To(Equal(decodeTest.wants.OriginCountry))

				// TimePosition
				if decodeTest.wants.TimePosition != nil {
					Expect(*(decodeData.TimePosition)).To(Equal(testTimeValue))
				} else {
					Expect(decodeData.TimePosition).To(BeNil())
				}

				// LastContact
				Expect(decodeData.LastContact).To(Equal(decodeTest.wants.LastContact))

				// Longitude
				if decodeTest.wants.Longitude != nil {
					Expect(*(decodeData.Longitude)).To(Equal(testFloatValue))
				} else {
					Expect(decodeData.Longitude).To(BeNil())
				}

				// Latitude
				if decodeTest.wants.Latitude != nil {
					Expect(*(decodeData.Latitude)).To(Equal(testFloatValue))
				} else {
					Expect(decodeData.Latitude).To(BeNil())
				}

				// BaroAltitude
				if decodeTest.wants.BaroAltitude != nil {
					Expect(*(decodeData.BaroAltitude)).To(Equal(testFloatValue))
				} else {
					Expect(decodeData.BaroAltitude).To(BeNil())
				}

				// OnGround
				Expect(decodeData.OnGround).To(Equal(decodeTest.wants.OnGround))

				// Velocity
				if decodeTest.wants.Velocity != nil {
					Expect(*(decodeData.Velocity)).To(Equal(testFloatValue))
				} else {
					Expect(decodeData.Velocity).To(BeNil())
				}

				// TrueTrack
				if decodeTest.wants.TrueTrack != nil {
					Expect(*(decodeData.TrueTrack)).To(Equal(testFloatValue))
				} else {
					Expect(decodeData.TrueTrack).To(BeNil())
				}

				// VerticalRate
				if decodeTest.wants.VerticalRate != nil {
					Expect(*(decodeData.VerticalRate)).To(Equal(testFloatValue))
				} else {
					Expect(decodeData.VerticalRate).To(BeNil())
				}

				// GeoAltitude
				if decodeTest.wants.GeoAltitude != nil {
					Expect(*(decodeData.GeoAltitude)).To(Equal(testFloatValue))
				} else {
					Expect(decodeData.GeoAltitude).To(BeNil())
				}

				// Squawk
				if decodeTest.wants.Squawk != nil {
					Expect(*(decodeData.Squawk)).To(Equal(testStringValue))
				} else {
					Expect(decodeData.Squawk).To(BeNil())
				}

				// Spi
				Expect(decodeData.Spi).To(Equal(decodeTest.wants.Spi))

				// PositionSource
				Expect(decodeData.PositionSource).To(Equal(decodeTest.wants.PositionSource))

				// Category
				Expect(decodeData.Category).To(Equal(decodeTest.wants.Category))
			}
		})
	})
})
