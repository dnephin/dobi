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

    go get github.com/dnephin/dobi

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
