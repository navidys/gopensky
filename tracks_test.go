package gopensky_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/h2non/gock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/navidys/gopensky"
)

var _ = Describe("Tracks", func() {
	Describe("GetTrackByAircraft", func() {
		It("retrieves the trajectory for a certain aircraft at a given time", func() {
			conn, err := gopensky.NewConnection(context.Background(), "", "")
			Expect(err).NotTo(HaveOccurred())

			defer gock.Off()
			gock.New(gopensky.OpenSkyAPIURL).
				Get("/tracks/all").
				Reply(200).
				File("mock_data/all_tracks.json")

			gclient, err := gopensky.GetClient(conn)
			Expect(err).NotTo(HaveOccurred())
			gock.InterceptClient(gclient)

			track, err := gopensky.GetTrackByAircraft(conn, "c060b9", 1689193028)
			Expect(err).NotTo(HaveOccurred())
			Expect(track.Icao24).To(Equal("c060b9"))
			Expect(*track.Callsign).To(Equal("POE2136"))
			Expect(track.StartTime).To(Equal(int64(1689193028)))
			Expect(track.EndTime).To(Equal(int64(1689197805)))
			Expect(track.Path).To(BeNil())
		})
	})

	Describe("GetTrackByAircraft - errors", func() {
		It("tests GetTrackByAircraft errors", func() {
			conn, err := gopensky.NewConnection(context.Background(), "", "")
			Expect(err).NotTo(HaveOccurred())

			gclient, err := gopensky.GetClient(conn)
			Expect(err).NotTo(HaveOccurred())
			gock.InterceptClient(gclient)

			_, err = gopensky.GetTrackByAircraft(context.Background(), "a835af", 1696755342)
			Expect(err.Error()).To(ContainSubstring("invalid context key"))

			_, err = gopensky.GetTrackByAircraft(conn, "", 1696755342)
			Expect(err).To(Equal(gopensky.ErrInvalidAircraftName))

			_, err = gopensky.GetTrackByAircraft(conn, "a835af", -1)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))

			defer gock.Off()

			gock.New(gopensky.OpenSkyAPIURL).
				Get("/tracks/all").
				Reply(200).
				BodyString("{'a': 2}")

			_, err = gopensky.GetTrackByAircraft(conn, "c060b9", 1689193028)
			Expect(err.Error()).To(ContainSubstring("unmarshalling"))

			gock.New(gopensky.OpenSkyAPIURL).
				Get("/tracks/invalid").
				Reply(200).
				BodyString("")

			_, err = gopensky.GetTrackByAircraft(conn, "c060b9", 1689193028)
			Expect(err.Error()).To(ContainSubstring("do request: Get"))

			for _, tfile := range []string{"tracks01.json", "tracks02.json", "tracks03.json", "tracks04.json"} {
				gock.New(gopensky.OpenSkyAPIURL).
					Get("/tracks/all").
					Reply(200).
					File("mock_data/errors/" + tfile)

				_, err = gopensky.GetTrackByAircraft(conn, "c060b9", 1689193028)
				Expect(errors.Unwrap(err).Error()).To(ContainSubstring("json: cannot unmarshal"))
			}
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
