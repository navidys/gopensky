# gopensky cmd

## Installation

* From source:

    ```
    $ make binary
    ```

* RPM package (COPR upstream build):

    ```
    $ sudo dnf copr enable navidys/gopensky
    $ sudo install gopensky
    ```

## Usage

```
$ gopensky -h
Query opensky network live API information (ADS-B and Mode S data)

Usage:
  gopensky [flags]
  gopensky [command]

Available Commands:
  arrivals    Retrieve flights for a certain airport which arrived within a given time interval
  help        Help about any command
  states      retrieve state vector information

Flags:
  -d, --debug             run in debug mode
  -h, --help              help for gopensky
  -j, --json              print json output
  -u, --username string   connection username
  -v, --version           print version and exit

```

### Example get states (json output)

```
$ gopensky states -j
{
    "time": 1699160183,
    "states": [
        {
            "icao24": "c067ae",
            "callsign": "ACA327  ",
            "origin_country": "Canada",
            "time_position": 1699160183,
            "last_contact": 1699160188,
            "longitude": -114.0126,
            "latitude": 51.1315,
            "baro_altitude": null,
            "on_ground": true,
            "velocity": 0,
            "true_track": 67.5,
            "vertical_rate": null,
            "sensors": null,
            "geo_altitude": null,
            "squawk": null,
            "spi": false,
            "position_source": 0,
            "category": 0
        },
        {
            ....
        }
    ]
}
```

### Example get states (table output)


```
$  gopensky states
```

![Screenshot](./docs/_static/gopensky-query.png)
