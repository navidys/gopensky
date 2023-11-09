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
  gopensky [command]

Available Commands:
  aircraft    Retrieve flights for a particular aircraft within a certain time interval
  arrivals    Retrieve flights for a certain airport which arrived within a given time interval
  departures  Retrieve flights for a certain airport which departed within a given time interval
  flights     Retrieve flights for a certain time interval
  help        Help about any command
  states      retrieve state vector information

Flags:
  -d, --debug             run in debug mode
  -h, --help              help for gopensky
  -j, --json              print json output
  -u, --username string   connection username
  -v, --version           print version and exit

Use "gopensky [command] --help" for more information about a command.
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

### Example get aircraft flights (table output)

Elon Musk's primary private jet (he has many and one of them has the ICAO24 transponder address `a835af`) flight data between Thu Aug 31 2023 23:11:04 and Fri Sep 29 2023 23:11:04.

```
$  gopensky aircraft -a a835af -b 1693523464 -e 1696029064

Header aliases:
  EDA   = estimated departure airport
  EAA   = estimated arrival airport
  EDAHD = estimated departure airport horizontal distance
  EDAVD = estimated departure airport vertical distance
  EAAHD = estimated arrival airport horizontal distance
  EAAVD = estimated arrival airport vertical distance
  DACC  = departure airport candidates count
  AACC  = arrival airport candidates count

ICAO24 EDA      EAA  CALLSIGN EDAHD EDAVD EAAHD EAAVD DACC AACC FIRST SEEN                    LAST SEEN
a835af N628TS   KUVA          154   154   14247 2417  425  189  2023-09-28 21:27:05 +0000 UTC 2023-09-28 21:46:02 +0000 UTC
a835af N628TS   02XS          0     0     1877  88    0    5    2023-09-19 00:00:06 +0000 UTC 2023-09-19 00:59:04 +0000 UTC
a835af N628TS   KSLC          0     0     3987  45    0    3    2023-09-16 00:00:05 +0000 UTC 2023-09-16 02:56:20 +0000 UTC
a835af N628TS   KIAD          25    25    1851  11    3    6    2023-09-13 02:14:03 +0000 UTC 2023-09-13 04:35:18 +0000 UTC
a835af N628TS   2TS2          34    34    3051  1415  2    13   2023-09-09 05:57:40 +0000 UTC 2023-09-09 08:40:39 +0000 UTC
a835af N628TS   KSJC          71    71    161   11    3    3    2023-09-07 19:21:04 +0000 UTC 2023-09-07 22:20:02 +0000 UTC

```
