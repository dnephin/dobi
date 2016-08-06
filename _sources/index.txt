dobi Documentation
==================

A build automation tool for Docker applications.


Getting Started
---------------

This is a short guide for getting started with **dobi**.

1. :doc:`install` **dobi**.
2. Create a ``dobi.yaml`` at the root of your project repository.
3. Add a mount, image, and run resource to the ``dobi.yaml``. See :doc:`examples` and
   :doc:`config` for more information about how to define a resource.
4. Run a task using the name of the resource, and an optional action name. See
   :doc:`tasks` for a full list of tasks.

   .. code:: sh

       dobi TASK [TASK...]

Documentation
-------------

.. toctree::
    :maxdepth: 2

    install
    examples
    config
    variables
    tasks
