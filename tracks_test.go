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

	Describe("parseFlightTrackResponse", func() {
		It("parses flight track api response", func() {
			testStringValue := "test_str"
			testTimeValue := int64(1)

			tests := []struct {
				have  gopensky.FlightTrackResponse
				wants gopensky.FlightTrack
			}{
				{
					have: gopensky.FlightTrackResponse{
						Icao24:    testStringValue,
						StartTime: float64(testTimeValue),
						EndTime:   float64(testTimeValue),
						Callsign:  &testStringValue,
						Path:      nil,
					},
					wants: gopensky.FlightTrack{
						Icao24:    testStringValue,
						StartTime: testTimeValue,
						EndTime:   testTimeValue,
						Callsign:  &testStringValue,
						Path:      nil,
					},
				},
				{
					have: gopensky.FlightTrackResponse{
						Icao24:    testStringValue,
						StartTime: 0,
						EndTime:   float64(testTimeValue),
						Callsign:  &testStringValue,
						Path:      nil,
					},
					wants: gopensky.FlightTrack{
						Icao24:    testStringValue,
						StartTime: 1,
						EndTime:   testTimeValue,
						Callsign:  &testStringValue,
						Path:      nil,
					},
				},
			}

			for _, respData := range tests {
				pdata, err := gopensky.ParseFlightTrackResponse(&respData.have)
				Expect(err).ToNot(HaveOccurred())

				Expect(pdata.Icao24).To(Equal(respData.wants.Icao24))
				Expect(pdata.StartTime).To(Equal(respData.wants.StartTime))
				Expect(pdata.EndTime).To(Equal(respData.wants.EndTime))

				if pdata.Path == nil {
					Expect(pdata.Path).To(BeNil())
				}
			}
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

	Describe("decodeWaypoint", func() {
		It("decode waypoint array interface to struct", func() {
			testTimeValue := int64(1)
			testFloatValue := -46.509700

			tests := []struct {
				have  []interface{}
				wants gopensky.WayPoint
			}{
				{have: []interface{}{testTimeValue}, wants: gopensky.WayPoint{}},
				{
					have: []interface{}{testTimeValue, testFloatValue, testFloatValue, testFloatValue, testFloatValue, false},
					wants: gopensky.WayPoint{
						Time:         testTimeValue,
						Latitude:     &testFloatValue,
						Longitude:    &testFloatValue,
						BaroAltitude: &testFloatValue,
						TrueTrack:    &testFloatValue,
						OnGround:     false,
					},
				},
				{
					have: []interface{}{testTimeValue, nil, nil, nil, nil, true},
					wants: gopensky.WayPoint{
						Time:         testTimeValue,
						Latitude:     nil,
						Longitude:    nil,
						BaroAltitude: nil,
						TrueTrack:    nil,
						OnGround:     true,
					},
				},
			}

			for _, decodeTest := range tests {
				path, err := gopensky.DecodeWaypoint(decodeTest.have)

				if len(decodeTest.have) < 6 {
					Expect(err).To(HaveOccurred())
					Expect(path).To(BeNil())

					continue
				}

				Expect(err).ToNot(HaveOccurred())
				Expect(path).ToNot(BeNil())

				// Time
				Expect(path.Time).To(Equal(decodeTest.wants.Time))

				// Latitude
				if decodeTest.wants.Latitude != nil {
					Expect(*(path.Latitude)).To(Equal(testFloatValue))
				} else {
					Expect(path.Latitude).To(BeNil())
				}

				// Longitude
				if decodeTest.wants.Longitude != nil {
					Expect(*(path.Longitude)).To(Equal(testFloatValue))
				} else {
					Expect(path.Longitude).To(BeNil())
				}

				// BaroAltitude
				if decodeTest.wants.BaroAltitude != nil {
					Expect(*(path.BaroAltitude)).To(Equal(testFloatValue))
				} else {
					Expect(path.BaroAltitude).To(BeNil())
				}

				// TrueTrack
				if decodeTest.wants.TrueTrack != nil {
					Expect(*(path.TrueTrack)).To(Equal(testFloatValue))
				} else {
					Expect(path.TrueTrack).To(BeNil())
				}

				// OnGround
				Expect(path.OnGround).To(Equal(decodeTest.wants.OnGround))
			}

		})
	})
})
