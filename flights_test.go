package gopensky_test

import (
	"context"
	"fmt"

	"github.com/navidys/gopensky"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flights", func() {
	Describe("getFlightsRequestParams", func() {
		It("generate flights request parameters", func() {
			tests := []struct {
				airport  string
				aircraft string
				begin    int64
				end      int64
			}{
				{airport: "", aircraft: "", begin: -1, end: -1},
				{airport: "EDDF", aircraft: "", begin: 1693523464, end: -1},
				{airport: "EDDF", aircraft: "", begin: 1693523464, end: 1696029064},
				{airport: "", aircraft: "a835af", begin: 1693523464, end: 1696029064},
				{airport: "EDDF", aircraft: "a835af", begin: 1693523464, end: 1696029064},
			}

			for _, reqParams := range tests {
				urlVal := gopensky.GetFlightsRequestParams(reqParams.airport, reqParams.aircraft, reqParams.begin, reqParams.end)
				Expect(urlVal.Get("airport")).To(Equal(reqParams.airport))
				Expect(urlVal.Get("icao24")).To(Equal(reqParams.aircraft))

				if reqParams.begin >= 0 {
					begin := fmt.Sprintf("%d", reqParams.begin)
					Expect(urlVal.Get("begin")).To(Equal(begin))
				} else {
					Expect(urlVal.Get("begin")).To(Equal(""))
				}

				if reqParams.end >= 0 {
					end := fmt.Sprintf("%d", reqParams.end)
					Expect(urlVal.Get("end")).To(Equal(end))
				} else {
					Expect(urlVal.Get("end")).To(Equal(""))
				}
			}
		})
	})

	Describe("GetArrivalsByAirport", func() {
		It("retrieves flights for a certain airport which arrived within a given time interval", func() {
			conn, err := gopensky.NewConnection(context.Background(), "", "")
			Expect(err).NotTo(HaveOccurred())

			_, err = gopensky.GetArrivalsByAirport(conn, "", 1696755342, 1696928142)
			Expect(err).To(Equal(gopensky.ErrInvalidAirportName))

			_, err = gopensky.GetArrivalsByAirport(conn, "LFPG", 0, 1696928142)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))

			_, err = gopensky.GetArrivalsByAirport(conn, "LFPG", 1696755342, -1)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))
		})
	})

	Describe("GetDeparturesByAirport", func() {
		It("retrieves flights for a certain airport which departed", func() {
			conn, err := gopensky.NewConnection(context.Background(), "", "")
			Expect(err).NotTo(HaveOccurred())

			_, err = gopensky.GetDeparturesByAirport(conn, "", 1696755342, 1696928142)
			Expect(err).To(Equal(gopensky.ErrInvalidAirportName))

			_, err = gopensky.GetDeparturesByAirport(conn, "LFPG", 0, 1696928142)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))

			_, err = gopensky.GetDeparturesByAirport(conn, "LFPG", 1696755342, -1)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))
		})
	})

	Describe("GetFlightsByInterval", func() {
		It("retrieves flights for a certain time interval", func() {
			conn, err := gopensky.NewConnection(context.Background(), "", "")
			Expect(err).NotTo(HaveOccurred())

			_, err = gopensky.GetFlightsByInterval(conn, 0, 1696928142)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))

			_, err = gopensky.GetFlightsByInterval(conn, 1696755342, -1)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))
		})
	})

	Describe("GetFlightsByAircraft", func() {
		It("retrieves flights for a particular aircraft within a certain time interval", func() {
			conn, err := gopensky.NewConnection(context.Background(), "", "")
			Expect(err).NotTo(HaveOccurred())

			_, err = gopensky.GetFlightsByAircraft(conn, "", 0, 1696928142)
			Expect(err).To(Equal(gopensky.ErrInvalidAircraftName))

			_, err = gopensky.GetFlightsByAircraft(conn, "a835af", 0, 1696928142)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))

			_, err = gopensky.GetFlightsByAircraft(conn, "a835af", 1696755342, -1)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))
		})
	})
})
