package gopensky_test

import (
	"context"
	"fmt"

	"github.com/h2non/gock"
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

			gclient, err := gopensky.GetClient(conn)
			Expect(err).NotTo(HaveOccurred())
			gock.InterceptClient(gclient)

			defer gock.Off()

			gock.New(gopensky.OpenSkyAPIURL).
				Get("/flights/arrival").
				Reply(200).
				BodyString("s")

			_, err = gopensky.GetArrivalsByAirport(conn, "KEWR", 1696755342, 1696928142)
			Expect(err.Error()).To(ContainSubstring("unmarshalling"))

			_, err = gopensky.GetArrivalsByAirport(conn, "", 1696755342, 1696928142)
			Expect(err).To(Equal(gopensky.ErrInvalidAirportName))

			_, err = gopensky.GetArrivalsByAirport(conn, "KEWR", 0, 1696928142)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))

			_, err = gopensky.GetArrivalsByAirport(conn, "KEWR", 1696755342, -1)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))

			_, err = gopensky.GetArrivalsByAirport(context.Background(), "KEWR", 1696755342, 1696928142)
			Expect(err.Error()).To(ContainSubstring("invalid context key"))

			gock.New(gopensky.OpenSkyAPIURL).
				Get("/flights/arrival").
				Reply(200).
				File("mock_data/flights_data.json")

			flightData, err := gopensky.GetArrivalsByAirport(conn, "KEWR", 1696755342, 1696928142)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(flightData)).To(Equal(3))

			Expect(flightData[0].Icao24).To(Equal("c060b9"))
			Expect(*(flightData[0].Callsign)).To(Equal("POE2136"))
			Expect(flightData[0].FirstSeen).To(Equal(int64(1689193028)))
			Expect(flightData[0].LastSeen).To(Equal(int64(1689197805)))
			Expect(*(flightData[0].EstArrivalAirport)).To(Equal("KEWR"))
			Expect(flightData[0].EstDepartureAirport).To(BeNil())
			Expect(flightData[0].ArrivalAirportCandidatesCount).To(Equal(3))
			Expect(flightData[0].EstDepartureAirportHorizDistance).To(Equal(int64(357)))
			Expect(flightData[0].EstDepartureAirportVertDistance).To(Equal(int64(24)))
			Expect(flightData[0].EstArrivalAirportHorizDistance).To(Equal(int64(591)))
			Expect(flightData[0].EstArrivalAirportVertDistance).To(Equal(int64(14)))
			Expect(flightData[0].DepartureAirportCandidatesCount).To(Equal(1))
		})
	})

	Describe("GetDeparturesByAirport", func() {
		It("retrieves flights for a certain airport which departed", func() {
			conn, err := gopensky.NewConnection(context.Background(), "", "")
			Expect(err).NotTo(HaveOccurred())

			gclient, err := gopensky.GetClient(conn)
			Expect(err).NotTo(HaveOccurred())
			gock.InterceptClient(gclient)

			defer gock.Off()

			gock.New(gopensky.OpenSkyAPIURL).
				Get("/flights/departure").
				Reply(200).
				BodyString("s")

			_, err = gopensky.GetDeparturesByAirport(conn, "KEWR", 1696755342, 1696928142)
			Expect(err.Error()).To(ContainSubstring("unmarshalling"))

			_, err = gopensky.GetDeparturesByAirport(conn, "", 1696755342, 1696928142)
			Expect(err).To(Equal(gopensky.ErrInvalidAirportName))

			_, err = gopensky.GetDeparturesByAirport(conn, "KEWR", 0, 1696928142)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))

			_, err = gopensky.GetDeparturesByAirport(conn, "KEWR", 1696755342, -1)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))

			_, err = gopensky.GetDeparturesByAirport(context.Background(), "KEWR", 1696755342, 1696928142)
			Expect(err.Error()).To(ContainSubstring("invalid context key"))

			gock.New(gopensky.OpenSkyAPIURL).
				Get("/flights/departure").
				Reply(200).
				File("mock_data/flights_data.json")

			flightData, err := gopensky.GetDeparturesByAirport(conn, "KEWR", 1696755342, 1696928142)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(flightData)).To(Equal(3))

			Expect(flightData[1].Icao24).To(Equal("c060b9"))
			Expect(*(flightData[1].Callsign)).To(Equal("RPA3462"))
			Expect(flightData[1].FirstSeen).To(Equal(int64(1689192822)))
			Expect(flightData[1].LastSeen).To(Equal(int64(1689196463)))
			Expect(*(flightData[1].EstDepartureAirport)).To(Equal("KEWR"))
			Expect(flightData[1].EstArrivalAirport).To(BeNil())
			Expect(flightData[1].ArrivalAirportCandidatesCount).To(Equal(6))
			Expect(flightData[1].EstDepartureAirportHorizDistance).To(Equal(int64(788)))
			Expect(flightData[1].EstDepartureAirportVertDistance).To(Equal(int64(9)))
			Expect(flightData[1].EstArrivalAirportHorizDistance).To(Equal(int64(201)))
			Expect(flightData[1].EstArrivalAirportVertDistance).To(Equal(int64(30)))
			Expect(flightData[1].DepartureAirportCandidatesCount).To(Equal(1))
		})
	})

	Describe("GetFlightsByInterval", func() {
		It("retrieves flights for a certain time interval", func() {
			conn, err := gopensky.NewConnection(context.Background(), "", "")
			Expect(err).NotTo(HaveOccurred())

			gclient, err := gopensky.GetClient(conn)
			Expect(err).NotTo(HaveOccurred())
			gock.InterceptClient(gclient)

			defer gock.Off()

			gock.New(gopensky.OpenSkyAPIURL).
				Get("/flights/all").
				Reply(200).
				BodyString("s")

			_, err = gopensky.GetFlightsByInterval(conn, 1696755342, 1696928142)
			Expect(err.Error()).To(ContainSubstring("unmarshalling"))

			_, err = gopensky.GetFlightsByInterval(conn, 0, 1696928142)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))

			_, err = gopensky.GetFlightsByInterval(conn, 1696755342, -1)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))

			_, err = gopensky.GetFlightsByInterval(context.Background(), 1696755342, 1696928142)
			Expect(err.Error()).To(ContainSubstring("invalid context key"))

			gock.New(gopensky.OpenSkyAPIURL).
				Get("/flights/all").
				Reply(200).
				File("mock_data/flights_data.json")

			flightData, err := gopensky.GetFlightsByInterval(conn, 1696755342, 1696928142)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(flightData)).To(Equal(3))

			Expect(flightData[0].Icao24).To(Equal("c060b9"))
			Expect(*(flightData[0].Callsign)).To(Equal("POE2136"))
			Expect(flightData[0].FirstSeen).To(Equal(int64(1689193028)))
			Expect(flightData[0].LastSeen).To(Equal(int64(1689197805)))
			Expect(*(flightData[0].EstArrivalAirport)).To(Equal("KEWR"))
			Expect(flightData[0].EstDepartureAirport).To(BeNil())
			Expect(flightData[0].ArrivalAirportCandidatesCount).To(Equal(3))
			Expect(flightData[0].EstDepartureAirportHorizDistance).To(Equal(int64(357)))
			Expect(flightData[0].EstDepartureAirportVertDistance).To(Equal(int64(24)))
			Expect(flightData[0].EstArrivalAirportHorizDistance).To(Equal(int64(591)))
			Expect(flightData[0].EstArrivalAirportVertDistance).To(Equal(int64(14)))
			Expect(flightData[0].DepartureAirportCandidatesCount).To(Equal(1))
		})
	})

	Describe("GetFlightsByAircraft", func() {
		It("retrieves flights for a particular aircraft within a certain time interval", func() {
			conn, err := gopensky.NewConnection(context.Background(), "", "")
			Expect(err).NotTo(HaveOccurred())

			gclient, err := gopensky.GetClient(conn)
			Expect(err).NotTo(HaveOccurred())
			gock.InterceptClient(gclient)

			defer gock.Off()
			gock.New(gopensky.OpenSkyAPIURL).
				Get("/flights/aircraft").
				Reply(200).
				BodyString("s")

			_, err = gopensky.GetFlightsByAircraft(conn, "c060b9", 1696755342, 1696928142)
			Expect(err.Error()).To(ContainSubstring("unmarshalling"))

			gock.New(gopensky.OpenSkyAPIURL).
				Get("/flights/aircraft").
				Reply(200).
				File("mock_data/flights_data.json")

			_, err = gopensky.GetFlightsByAircraft(conn, "", 0, 1696928142)
			Expect(err).To(Equal(gopensky.ErrInvalidAircraftName))

			_, err = gopensky.GetFlightsByAircraft(conn, "c060b9", 0, 1696928142)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))

			_, err = gopensky.GetFlightsByAircraft(conn, "c060b9", 1696755342, -1)
			Expect(err).To(Equal(gopensky.ErrInvalidUnixTime))

			_, err = gopensky.GetFlightsByAircraft(context.Background(), "c060b9", 1696755342, 1696928142)
			Expect(err.Error()).To(ContainSubstring("invalid context key"))

			flightData, err := gopensky.GetFlightsByAircraft(conn, "c060b9", 1696755342, 1696928142)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(flightData)).To(Equal(3))

			Expect(flightData[2].Icao24).To(Equal("c060b9"))
			Expect(*(flightData[2].Callsign)).To(Equal("N401TD"))
			Expect(flightData[2].FirstSeen).To(Equal(int64(1689192818)))
			Expect(flightData[2].LastSeen).To(Equal(int64(1689198430)))
			Expect(*(flightData[2].EstArrivalAirport)).To(Equal("KEWR"))
			Expect(flightData[2].EstDepartureAirport).To(BeNil())
			Expect(flightData[2].ArrivalAirportCandidatesCount).To(Equal(4))
			Expect(flightData[2].EstDepartureAirportHorizDistance).To(Equal(int64(13461)))
			Expect(flightData[2].EstDepartureAirportVertDistance).To(Equal(int64(24)))
			Expect(flightData[2].EstArrivalAirportHorizDistance).To(Equal(int64(204)))
			Expect(flightData[2].EstArrivalAirportVertDistance).To(Equal(int64(8)))
			Expect(flightData[2].DepartureAirportCandidatesCount).To(Equal(1))
		})
	})
})
