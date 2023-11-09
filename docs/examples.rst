Examples
========

All State Vectors
--------------------
Example for retrieving all states without authentication:

.. code-block:: go

    package main

    import (
        "context"
        "fmt"
        "os"

        "github.com/navidys/gopensky"
    )

    func main() {
        conn, err := gopensky.NewConnection(context.Background(), "", "")
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // retrieve all states information
        statesData, err := gopensky.GetStates(conn, 0, nil, nil, true)
        if err != nil {
            fmt.Println(err)
            os.Exit(2)
        }

        for _, state := range statesData.States {
            longitude := "<nil>"
            latitude := "<nil>"

            if state.Longitude != nil {
                longitude = fmt.Sprintf("%f", *state.Longitude)
            }

            if state.Latitude != nil {
                latitude = fmt.Sprintf("%f", *state.Latitude)
            }

            fmt.Printf("ICAO24: %s, Longitude: %s, Latitude: %s, Origin Country: %s \n",
                state.Icao24,
                longitude,
                latitude,
                state.OriginCountry,
            )
        }
    }

output:

.. code-block:: bash

    ICAO24: e49406, Longitude: -46.509700, Latitude: -23.144100, Origin Country: Brazil
    ICAO24: ad4f1c, Longitude: -81.208700, Latitude: 29.430400, Origin Country: United States
    ICAO24: 7c6b2d, Longitude: 171.803200, Latitude: -43.586000, Origin Country: Australia
    ICAO24: 4b1812, Longitude: 8.498200, Latitude: 47.461400, Origin Country: Switzerland
    ICAO24: 88044e, Longitude: 108.214900, Latitude: 18.868100, Origin Country: Thailand
    ICAO24: 440da1, Longitude: 17.950300, Latitude: 59.645800, Origin Country: Austria
    ICAO24: ab6fdd, Longitude: -96.840200, Latitude: 38.359400, Origin Country: United States
    ....
    ....
    ....


Arrivals by Airport
--------------------
Example of retrieving flights for a certain airport which arrived within a given time interval:

.. code-block:: go

    package main

    import (
        "context"
        "fmt"
        "os"
        "time"

        "github.com/navidys/gopensky"
    )

    func main() {
        conn, err := gopensky.NewConnection(context.Background(), "", "")
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // retrieve arrivals flights of:
        // airport: LFPG (Charles de Gaulle)
        // being time: 1696755342 (Sunday October 08, 2023 08:55:42 UTC)
        // end time: 1696928142 (Tuesday October 10, 2023 08:55:42 UTC)

        flightsData, err := gopensky.GetArrivalsByAirport(conn, "LFPG", 1696755342, 1696928142)
        if err != nil {
            fmt.Println(err)
            os.Exit(2)
        }

        for _, flightData := range flightsData {
            var (
              depAirport string
              callSign   string
            )

            if flightData.EstDepartureAirport != nil {
              depAirport = *flightData.EstDepartureAirport
            }

            if flightData.Callsign != nil {
              callSign = *flightData.Callsign
            }

            fmt.Printf("ICAO24: %s, callSign: %8s, Departed Airport: %4s, LastSeen: %s\n",
              flightData.Icao24,
              callSign,
              depAirport,
              time.Unix(flightData.LastSeen, 0),
            )
        }
    }

.. ::

output:

.. code-block:: bash

    ICAO24: 406544, callSign: EZY24ZB , Departed Airport: EGPH, LastSeen: 2023-10-10 07:33:07 +1100 AEDT
    ICAO24: 896180, callSign: UAE75   , Departed Airport:     , LastSeen: 2023-10-10 05:07:35 +1100 AEDT
    ICAO24: 738065, callSign: ELY327  , Departed Airport: LLBG, LastSeen: 2023-10-10 03:14:58 +1100 AEDT
    ICAO24: 4bc848, callSign: PGT412L , Departed Airport: LTFJ, LastSeen: 2023-10-10 01:31:15 +1100 AEDT
    ICAO24: 4891b6, callSign: ENT76YA , Departed Airport:     , LastSeen: 2023-10-09 20:52:38 +1100 AEDT
    ICAO24: 39856a, callSign: AFR55QA , Departed Airport: LFBO, LastSeen: 2023-10-09 20:45:12 +1100 AEDT
    ICAO24: 4ba9c9, callSign: THY5ER  , Departed Airport: LTFM, LastSeen: 2023-10-09 18:52:45 +1100 AEDT
    ICAO24: 738075, callSign:         , Departed Airport: LFPG, LastSeen: 2023-10-09 16:03:10 +1100 AEDT
    ICAO24: 39e68b, callSign: AFR61MP , Departed Airport: ESSA, LastSeen: 2023-10-09 07:23:04 +1100 AEDT
    ICAO24: 01020a, callSign: MSR801  , Departed Airport:     , LastSeen: 2023-10-09 05:46:24 +1100 AEDT
    ICAO24: 39e698, callSign: AFR94RP , Departed Airport: LOWW, LastSeen: 2023-10-09 04:51:45 +1100 AEDT
    ICAO24: 398569, callSign: AFR37LV , Departed Airport: LJLJ, LastSeen: 2023-10-09 02:03:00 +1100 AEDT

Departures by Airport
----------------------
Example of retrieving flights for a certain airport which departed within a given time interval:

.. code-block:: go

    package main

    import (
      "context"
      "encoding/json"
      "fmt"
      "os"

      "github.com/navidys/gopensky"
      "github.com/rs/zerolog/log"
    )

    func main() {
      conn, err := gopensky.NewConnection(context.Background(), "", "")
      if err != nil {
        fmt.Println(err)
        os.Exit(1)
      }

      // retrieve departed flights of:
      // airpor: LFPG (Charles de Gaulle)
      // being time: 1696755342 (Sunday October 08, 2023 08:55:42 UTC)
      // end time: 1696928142 (Tuesday October 10, 2023 08:55:42 UTC)

      flightsData, err := gopensky.GetDeparturesByAirport(conn, "LFPG", 1696755342, 1696928142)
      if err != nil {
        fmt.Println(err)
        os.Exit(2)
      }

      jsonResult, err := json.MarshalIndent(flightsData, "", "    ")
      if err != nil {
        log.Error().Msgf("%v", err)

        return
      }

      fmt.Printf("%s\n", jsonResult)
    }

.. ::

output:


.. code-block:: bash

  [
      {
          "icao24": "502ce6",
          "firstSeen": 1696927909,
          "estDepartureAirport": "LFPG",
          "lastSeen": 1696935895,
          "estArrivalAirport": "EVTA",
          "callsign": "BTI69R  ",
          "estDepartureAirportHorizDistance": 1958,
          "estDepartureAirportVertDistance": 86,
          "estArrivalAirportHorizDistance": 17721,
          "estArrivalAirportVertDistance": 3186,
          "departureAirportCandidatesCount": 91,
          "arrivalAirportCandidatesCount": 0
      },
      {
          "icao24": "39856f",
          "firstSeen": 1696927889,
          "estDepartureAirport": "LFPG",
          "lastSeen": 1696931551,
          "estArrivalAirport": "LEVT",
          "callsign": "AFR21YB ",
          "estDepartureAirportHorizDistance": 2312,
          "estDepartureAirportVertDistance": 71,
          "estArrivalAirportHorizDistance": 62514,
          "estArrivalAirportVertDistance": 2169,
          "departureAirportCandidatesCount": 91,
          "arrivalAirportCandidatesCount": 10
      },
      ...
      ...
      ...
    }
