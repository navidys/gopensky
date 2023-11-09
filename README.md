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

| Function  | Description |
| --------- | ----------- |
| GetStates |  Retrieve state vectors for a given time.
| GetArrivalsByAirport | Retrieves flights for a certain airport which arrived within a given time interval.
| GetDeparturesByAirport | Retrieves flights for a certain airport which departed within a given time interval.
| GetFlightsByInterval | Retrieves flights for a certain time interval.

## Examples

Here is an example program of retrieving flights between SMonday, October 9, 2023 6:19:28 and Monday, October 9, 2023 7:19:28.

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

	// retrieve flights
	// being time: 1696832368 (Monday, October 9, 2023 6:19:28)
	// end time: 1696835968 (Monday, October 9, 2023 7:19:28)

	flightsData, err := gopensky.GetFlightsByInterval(conn, 1696832368, 1696835968)
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
ICAO24: 008833, Departed: FAGC, Arrival: FAFF, LastSeen: 2023-10-09 17:24:52 +1100 AEDT
ICAO24: 008de4, Departed: FAOR, Arrival:     , LastSeen: 2023-10-09 17:48:42 +1100 AEDT
ICAO24: 008dea, Departed: FAOR, Arrival:     , LastSeen: 2023-10-09 17:16:38 +1100 AEDT
ICAO24: 008df9, Departed: FAOR, Arrival: FANC, LastSeen: 2023-10-09 17:37:19 +1100 AEDT
ICAO24: 009893, Departed: FAOR, Arrival: FABS, LastSeen: 2023-10-09 18:09:45 +1100 AEDT
ICAO24: 00af2e, Departed: FAOR, Arrival: FANC, LastSeen: 2023-10-09 17:38:33 +1100 AEDT
ICAO24: 00b097, Departed: FAOR, Arrival: FAWN, LastSeen: 2023-10-09 17:12:13 +1100 AEDT
ICAO24: 0100e4, Departed:     , Arrival:     , LastSeen: 2023-10-09 17:20:43 +1100 AEDT
ICAO24: 01010b, Departed:     , Arrival:     , LastSeen: 2023-10-09 17:27:32 +1100 AEDT
ICAO24: 0101ba, Departed:     , Arrival: HE13, LastSeen: 2023-10-09 17:37:38 +1100 AEDT
ICAO24: 0101cd, Departed:     , Arrival: HE28, LastSeen: 2023-10-09 17:54:49 +1100 AEDT
ICAO24: 01022e, Departed: EDDK, Arrival:     , LastSeen: 2023-10-09 18:19:19 +1100 AEDT
...
...
...
```

## License

Licensed under the [Apache 2.0](LICENSE) license.
