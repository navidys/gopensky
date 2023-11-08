Go OpenSky API Functions
==============================

.. _FUNC_CONNECTION:

func :ref:`NewConnection <FUNC_CONNECTION>`
--------------------------------------------

    Creates a new connection context to OpenSky Network live API server.


    .. code-block:: go

        func NewConnection(ctx context.Context, username string, password string) (context.Context, error)


    :Parameters:
        - **ctx** (`context.Context <https://pkg.go.dev/context#Context>`_) - connection context.
        - **username** (string) - an OpenSky username (Anonymous connection will be use by providing empty username).
        - **password** (string) - an OpenSky password for the given username.

    :Returns: context.Context, error

.. _FUNC_GET_STATES:

func :ref:`GetStates <FUNC_GET_STATES>`
--------------------------------------------

    Retrieve state vectors for a given time.

    .. code-block:: go

        func GetStates(ctx context.Context, time int64, icao24 []string, bBox *BoundingBoxOptions, extended bool) (*States, error)


    :Parameters:
        - **ctx** (`context.Context <https://pkg.go.dev/context#Context>`_) - connection context.
        - **time** (int64) - time as Unix time stamp (seconds since epoch) or datetime. The datetime must be in UTC!. If ``time = 0`` the most recent ones are taken.
        - **icao24** ([]string)  - optionally retrieve only state vectors for the given ICAO24 address(es). The parameter an array of str containing multiple addresses.
        - **bBox** (:ref:`*BoundingBoxOptions<TYPE_BBOX_OPTIONS>`) - optionally retrieve state vectors within a bounding box. Use :ref:`NewBoundingBox<BBOX_FUNC>` function to create a new one.
        - **extended** (bool) - set to ``true`` to request the category of aircraft

    :Returns: :ref:`*States<TYPE_STATES>`, error


.. _FUNC_GET_ARRIVALS_BY_AIRPORT:

func :ref:`GetArrivalsByAirport <FUNC_GET_ARRIVALS_BY_AIRPORT>`
----------------------------------------------------------------

    Retrieves flights for a certain airport which arrived within a given time interval [being, end].

    The given time interval must not be larger than seven days!

    .. code-block:: go

        func GetArrivalsByAirport(ctx context.Context, airport string, begin int64, end int64) ([]FlightData, error)


    :Parameters:
        - **ctx** (`context.Context <https://pkg.go.dev/context#Context>`_) - connection context.
        - **airport** (string) - ICAO identifier for the airport.
        - **begin** (int64) - Start of time interval to retrieve flights for as Unix time (seconds since epoch).
        - **end** (int64)  - End of time interval to retrieve flights for as Unix time (seconds since epoch).

    :Returns: :ref:`[]FlightData<TYPE_FLIGHT_DATA>`, error

.. _BBOX_FUNC:

func :ref:`NewBoundingBox <BBOX_FUNC>`
--------------------------------------------

    Creates a new bounding (min_latitude, max_latitude, min_longitude, max_longitude) box option.

    .. code-block:: go

        func NewBoundingBox (lamin float64, lomin float64, lamax float64, lomax float64) *BoundingBoxOptions

    :Parameters:
        - **lamin** (float64) - lower bound for the latitude in WGS84 decimal degrees.
        - **lomin** (float64) - lower bound for the longitude in in WGS84 decimal degrees.
        - **lamax** (float64) - upper bound for the latitude in WGS84 decimal degrees.
        - **lomax** (float64) - upper bound for the longitude in in WGS84 decimal degrees.

    :Returns: :ref:`*BoundingBoxOptions<TYPE_BBOX_OPTIONS>`
