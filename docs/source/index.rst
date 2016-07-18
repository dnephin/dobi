dobi Documentation
==================

A build automation tool for Docker applications.


Getting Started
---------------

**dobi** uses a YAML config file (``dobi.yaml`` by default) to store the project
tasks required to build, test, and deploy your application.  The first step is
to create a ``dobi.yaml``.

Once you've created a ``dobi.yaml`` you can install the **dobi** command and run
the tasks.

.. code:: sh

    dobi <resource>

Documentation
-------------

.. toctree::
    :maxdepth: 2

    examples
    config
    variables
    tasks


Install
-------

There are two install options:

Download the binary
~~~~~~~~~~~~~~~~~~~

Binaries are available for Linux, OSX, and Windows. Download a binary from
`github.com/dnephin/dobi/releases <https://github.com/dnephin/dobi/releases>`_

Install from source
~~~~~~~~~~~~~~~~~~~

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


Indices
=======
* :ref:`genindex`
* :ref:`search`
