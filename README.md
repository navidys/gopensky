# Go OpenSKY Network API
![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/navidys/gopensky)](https://pkg.go.dev/github.com/navidys/gopensky)
[![Go Report](https://goreportcard.com/badge/github.com/navidys/gopensky)](https://goreportcard.com/report/github.com/navidys/gopensky)

This is the golang implementation of the OpenSky network's live API.
The API lets you retrieve live airspace information (ADS-B and Mode S data) for research and non-commerical purposes.

For documentation and examples visit [Golang OpenSky Network API](https://navidys.github.io/gopensky/).

A `gopensky` binary command line is also available to query the opensky network api, check [installation guide](./INSTALL.md) for RPM package installation or building from source.

`NOTE:` there are some limitation sets for anonymous and OpenSky users, visit following links for more information:
* [OpenSky Network Rest API documentation](https://openskynetwork.github.io/opensky-api/)
* [OpenSky Network Website](https://opensky-network.org/).

## Getting Started

Install the latest version of the library:

```
$ go get github.com/navidys/gopensky
```

Next, include gopensky in you application:

```
import "github.com/navidys/gopensky"
```

## Features

* `GetStates` - retrieve state vectors for a given time.
* `GetArrivalsByAirport` - retrieves flights for a certain airport which arrived within a given time interval.
* `GetDeparturesByAirport` - retrieves flights for a certain airport which departed within a given time interval.
* `GetFlightsByInterval` - retrieves flights for a certain time interval.
* `GetFlightsByAircraft` - retrieves flights for a particular aircraft within a certain time interval.
* `GetTrackByAircraft` - retrieves the trajectory for a certain aircraft at a given time.

## Examples

Here is an example program of retrieving Elon Musk's primary private jet (he has many and one of them has the ICAO24 transponder address `a835af`) flights between Thu Aug 31 2023 23:11:04 and Fri Sep 29 2023 23:11:04.

Visit [Golang OpenSky Network API](https://navidys.github.io/gopensky/) for more examples.

```
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

	// Retrieve Elon Musk's primary private jet flight data, he has many and one of
	// them has the ICAO24 transponder address a835af
	// begin time: 1693523464 (Thu Aug 31 2023 23:11:04)
	// end time: 1696029064 (Fri Sep 29 2023 23:11:04)

	flightsData, err := gopensky.GetFlightsByAircraft(conn, "a835af", 1693523464, 1696029064)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	for _, flightData := range flightsData {
		var (
			departedAirport string
			arrivalAriport  string
		)

		if flightData.EstDepartureAirport != nil {
			departedAirport = *flightData.EstDepartureAirport
		}

		if flightData.EstArrivalAirport != nil {
			arrivalAriport = *flightData.EstArrivalAirport
		}

		fmt.Printf("ICAO24: %s, Departed: %4s, Arrival: %4s, LastSeen: %s\n",
			flightData.Icao24,
			departedAirport,
			arrivalAriport,
			time.Unix(flightData.LastSeen, 0),
		)
	}
}
```

output:

```
ICAO24: a835af, Departed: KAUS, Arrival: KUVA, LastSeen: 2023-09-29 07:46:02 +1000 AEST
ICAO24: a835af, Departed:     , Arrival: 02XS, LastSeen: 2023-09-19 10:59:04 +1000 AEST
ICAO24: a835af, Departed:     , Arrival: KSLC, LastSeen: 2023-09-16 12:56:20 +1000 AEST
ICAO24: a835af, Departed: KAUS, Arrival: KIAD, LastSeen: 2023-09-13 14:35:18 +1000 AEST
ICAO24: a835af, Departed: KSJC, Arrival: 2TS2, LastSeen: 2023-09-09 18:40:39 +1000 AEST
ICAO24: a835af, Departed: KAUS, Arrival: KSJC, LastSeen: 2023-09-08 08:20:02 +1000 AEST
```

## License

Licensed under the [Apache 2.0](LICENSE) license.
