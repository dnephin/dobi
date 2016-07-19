dobi Documentation
==================

A build automation tool for Docker applications.


Getting Started
---------------

**dobi** uses a YAML config file (:doc:`dobi.yaml <config>` by default) to store the project
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

    install
    examples
    config
    variables
    tasks
