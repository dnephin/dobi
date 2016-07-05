
dobi
====

A build automation tool for Docker applications.

``dobi`` allows you to define the resources required to build,
test, and deploy your application.  Each resource may depend on other resource.
When a resource is stale ``dobi`` runs the appropriate tasks to re-build it.

Resources are defined in a ``dobi.yaml`` file, and are run using the command
``dobi <resource name>``. Multiple tasks can be joined together with an
``alias`` resource to define the high level operations.

``dobi`` differs from other build automation tools by making the Docker
containers, images, mounts, and environments a first-class resource.

See `Getting Started <./docs/index.rst>`_

.. image:: https://circleci.com/gh/dnephin/dobi/tree/master.svg?style=svg
    :target: https://circleci.com/gh/dnephin/dobi/tree/master

Features
--------

* only re-build out-of-date resources, so operations are fast.
* everything runs in a container, so operations are portable and reliable.
* dependencies are automatically run first, and out-of-date dependencies force
  dependents to be re-built.
* tasks are grouped together to create high-level operations. Developer
  operations are encoded in an easy to read configuration file.


Install
-------

Download the binary
~~~~~~~~~~~~~~~~~~~

From `dnephin/dobi/releases <https://github.com/dnephin/dobi/releases>`_

Install from source
~~~~~~~~~~~~~~~~~~~

.. code::

    go get github.com/dnephin/dobi

Usage
-----

See ``dobi --help`` for full usage.

Examples
--------

* `dobi <https://github.com/dnephin/dobi/blob/master/dobi.yaml>`_ - a Golang
  command line tool (``dobi`` uses itself!)


Documentation
-------------

See `docs <./docs/index.rst>`_ for complete documentation.

Contributing
------------

``dobi`` is still in early development. If you'd like to contribute, please open
an issue, or find an existing issue, and leave a comment saying you're working
on a feature.
