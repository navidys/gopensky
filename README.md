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

Here is an example program once you've added the library, visit [Golang OpenSky Network API](https://navidys.github.io/gopensky/) for more examples.

```
import (
	"context"
	"fmt"
	"os"

	"github.com/navidys/gopensky"
)

func main() {
	conn, err := gopensky.NewConnection(context.Background(), "username", "password")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// retrieve all states information
	statesData, err := gopensky.GetStates(conn, 0, nil, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	for _, state := range statesData.States {
		fmt.Printf("ICAO24: %s, Origin Country: %s, Longitude: %v, Latitude: %v \n", state.Icao24, state.OriginCountry, state.Longitude, state.Latitude)
	}
}
```

## License

Licensed under the [Apache 2.0](LICENSE) license.