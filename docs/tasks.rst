Tasks
=====

Each resource defined in the ``dobi.yml`` provides one or more tasks. Each
resource has a default task, which is usually the "create" or "build" action
associated with that resource (build an image, run a container, etc).

Each resource also defines a "remove" task which can be used to remove the
artifact the resource created.

Tasks are run using ``dobi RESOURCE[:ACTION]`` where ``RESOURCE:ACTION`` is the
name of the task.


.. contents::
    :backlinks: none
    :depth: 2


Image Tasks
-----------

Each image resource has the following tasks.

Build (default)
~~~~~~~~~~~~~~~

Build (``:build``) a Docker image from a Dockerfile.


Push
~~

Push (``:push``) the image to a registry.

The ``push`` task always depends on the ``build`` task for the image.


Pull
~~~~


Pull (``:pull``) the image from a registry.


Remove
~~~~~~

Remove (``:rm``) the image.


Compose Tasks
-------------

Up (default)
~~~~~~~~~~~~

Up (``:up``) runs ``docker-compose up -d`` with the files and project name from
the resource to create a new isolated environment.

Down
~~~~

Down (``:down``, or ``:rm``) removes all the containers and networks created by
Compose.

Attach
~~~~~~

Attach (``:attach``) runs ``docker-compose up`` and attaches to the logs.


Run Tasks
---------

Run (default)
~~~~~~~~~~~~~

Run (``:run``) runs a conatiner.


Mount Tasks
-----------

Create (default)
~~~~~~~~~~~~~~~~

Create (``:create``) creates the host directory to be bind mounted.


Alias Tasks
-----------

Run (default)
~~~~~~~~~~~~~

Run (``:run``) runs all the tasks in the list of tasks

Remove
~~~~~~

Remove (``:rm``) runs the remove task for all the resources in the task list in
reverse order.
