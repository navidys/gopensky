# Go OpenSKY Network API
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

Here is an example program of retrieving flights for a Charles de Gaulle airport between Sunday October 08, 2023 08:55:42 and Tuesday October 10, 2023 08:55:42.

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
```

output:

```
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
```

## License

Licensed under the [Apache 2.0](LICENSE) license.
