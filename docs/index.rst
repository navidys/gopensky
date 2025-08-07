Go OpenSky Network API documentation
====================================

**gopensky**  is the Go implementation of the OpenSky network's live API.
It lets you retrieve live airspace information (ADS-B and Mode S data) for research and non-commerical purposes.


There are some limitation sets for anonymous and OpenSky users, visit following links for more information:

    * `OpenSky Network Rest API documentation <https://openskynetwork.github.io/opensky-api/>`_
    * `OpenSky Network Website <https://opensky-network.org/>`_

Installation
--------------

Use ``go get`` to install the latest version of the library:

.. code-block:: bash

    $ go get github.com/navidys/gopensky

Next, include gopensky in you application:

.. code-block:: go

    import "github.com/navidys/gopensky"


Further Reading
---------------

.. toctree::
    :maxdepth: 1

    Introduction <self>
    Go API Functions<goapi_functions>
    Go API Types <goapi_types>
    Examples <examples>
