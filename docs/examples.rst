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
        // airpor: LFPG (Charles de Gaulle)
        // being time: 1696755342 (Sunday October 08, 2023 08:55:42 UTC)
        // end time: 1696928142 (Tuesday October 10, 2023 08:55:42 UTC)

        flightsData, err := gopensky.GetArrivalsByAirport(conn, "LFPG", 1696755342, 1696928142)
        if err != nil {
            fmt.Println(err)
            os.Exit(2)
        }

        for _, flightData := range flightsData {
            var depAirport string
            if flightData.EstDepartureAirport != nil {
                depAirport = *flightData.EstDepartureAirport
            }

            fmt.Printf("ICAO24: %s, Departure Airport: %4s, LastSeen: %s\n",
                flightData.Icao24,
                depAirport,
                time.Unix(flightData.LastSeen, 0),
            )
        }
    }

.. ::

output:

.. code-block:: bash

    ICAO24: 406544, Departure Airport: EGPH, LastSeen: 2023-10-10 07:33:07 +1100 AEDT
    ICAO24: 896180, Departure Airport:     , LastSeen: 2023-10-10 05:07:35 +1100 AEDT
    ICAO24: 738065, Departure Airport: LLBG, LastSeen: 2023-10-10 03:14:58 +1100 AEDT
    ICAO24: 4bc848, Departure Airport: LTFJ, LastSeen: 2023-10-10 01:31:15 +1100 AEDT
    ICAO24: 4891b6, Departure Airport:     , LastSeen: 2023-10-09 20:52:38 +1100 AEDT
    ICAO24: 39856a, Departure Airport: LFBO, LastSeen: 2023-10-09 20:45:12 +1100 AEDT
    ICAO24: 4ba9c9, Departure Airport: LTFM, LastSeen: 2023-10-09 18:52:45 +1100 AEDT
    ICAO24: 738075, Departure Airport: LFPG, LastSeen: 2023-10-09 16:03:10 +1100 AEDT
    ICAO24: 39e68b, Departure Airport: ESSA, LastSeen: 2023-10-09 07:23:04 +1100 AEDT
    ICAO24: 01020a, Departure Airport:     , LastSeen: 2023-10-09 05:46:24 +1100 AEDT
    ICAO24: 39e698, Departure Airport: LOWW, LastSeen: 2023-10-09 04:51:45 +1100 AEDT
    ICAO24: 398569, Departure Airport: LJLJ, LastSeen: 2023-10-09 02:03:00 +1100 AEDT
