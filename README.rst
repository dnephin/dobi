
dobi
====

A build system for Docker applications.

``dobi`` allows you to define the resources and tasks required to build,
test, and deploy your application.  Each resource may depend on other resource.
When a resource is stale ``dobi`` runs the appropriate task to re-build it.

Resources are defined in a ``dobi.yaml`` file, and are run using the command
``dobi <resource name>``.

``dobi`` differs from other build tools by making the Docker containers, images,
mounts, and environments a first-class resource.


Install
-------

Currently the only install option is ``go get``. After the first official
release a binary download will be provided.

.. code::

    go get github.com/dnephin/dobi

Usage
-----

See ``dobi --help`` for full usage instructions, and ``dobi.yaml`` in this repo
for an example configuration.

Every resource you define in a ``dobi.yaml`` is a runnable task, using the
resource name.

.. code::

    # Build the binary
    dobi binary

    # Run the unit tests
    dobi test-unit


Documentation
-------------

See `docs/ <./docs/index.rst>`_ for complete documentation.

Contributing
------------

.. image:: https://circleci.com/gh/dnephin/dobi/tree/master.svg?style=svg
    :target: https://circleci.com/gh/dnephin/dobi/tree/master

``dobi`` is still in early development. If you'd like to contribute, please open
an issue, or find an existing issue, and leave a comment saying you're working
on a feature.
