Examples
========

All State Vectors
--------------------
Example for retrieving all states without authentication:

.. code-block:: go

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
        statesData, err := gopensky.GetStates(conn, 0, nil, nil)
        if err != nil {
            fmt.Println(err)
            os.Exit(2)
        }

        for _, state := range statesData.States {
            fmt.Printf("ICAO24: %s, Origin Country: %s, Longitude: %v, Latitude: %v \n",
                state.Icao24,
                state.OriginCountry,
                state.Longitude,
                state.Latitude,
            )
        }
    }
