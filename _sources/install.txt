Install
=======


There are two install options:

Download the binary
-------------------

Binaries are available for Linux, OSX, and Windows. Download a binary from
`github.com/dnephin/dobi/releases <https://github.com/dnephin/dobi/releases>`_

Install from source
-------------------

.. code:: sh

    git clone git@github.com:dnephin/dobi && cd dobi
    docker run -ti --rm -w $(pwd) -v $(pwd):$(pwd) -e DOCKER_HOST \
        -v /var/run/docker.sock:/var/run/docker.sock \
        dnephin/dobi:0.5 deps binary

The binaries will be in ``./dist/bin``

Usage
-----

To run a task use the name of the resource from the ``dobi.yml`` config. For
example to run a ``run=test`` resource:

.. code:: sh

    dobi test

Many resources also have other actions such as `:rm`:

.. code:: sh

    dobi test:rm


See :doc:`tasks` for a full of actions.

See ``dobi --help`` for full usage.
