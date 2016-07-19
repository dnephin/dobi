
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

See `Getting Started <https://dnephin.github.io/dobi/>`_

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

See `Install <https://dnephin.github.io/dobi/#install>`_.

Examples
--------

See `Examples <https://dnephin.github.io/dobi/examples.html>`_.

Documentation
-------------

See `Documentantion <https://dnephin.github.io/dobi/>`_


Contributing
------------

``dobi`` is still in early development. If you'd like to contribute, please open
an issue, or find an existing issue, and leave a comment saying you're working
on a feature.
