Tasks
=====

Each resource defined in the ``dobi.yml`` provides one or more tasks. Each
resource has a default task, which is usually the "create" or "build" action
associated with that resource (build an image, run a container, etc).

Each resource also defines a ``remove`` task which can be used to remove
anything that was created by the create action of the resource.

To run a task use the name of the resource, an optionally an action name.

.. code-block:: sh

    # Run the test resource
    dobi test

    # Run the remove action for the builder resource
    dobi builder:rm

To list all the tasks in a project run

.. code-block:: sh

    dobi list


Image Tasks
-----------

`image <./config.html#image>`_ resources have the following tasks:

``:build`` *(default)*
~~~~~~~~~~~~~~~~~~~~~~

Build a Docker image from a Dockerfile. The image is tagged using the **image**
field and the first tag from the list of **tags** in the image resource. If the
**tags** field is not set, the value of ``{unique}`` will be used as the time. See
:doc:`variables` for more information about how to set the unique value.


.. note::

   For every image built by **dobi** a small file is created in the ``./.dobi/images/``
   directory (relative to the directory which contains the ``dobi.yaml``). This
   file is used to track the modified time of the image. This file is necessary
   because Docker does not store the "last built" time of an image, only the first
   time it was built. If the image is built again, and the build is completely cached,
   no new image id gets created, so the "created time" of the image is actually the
   original created time.

   Without this file images often appear as stale, because the original created time
   of the image is earlier than the last attempted build. This will happen when a file
   in the image context is modified, but that modification doesn't change the docker
   build cache (which is the case if the modified file is never added to the image
   using ``COPY`` or ``ADD``). By saving the image id in a local file, and using the
   modified time of that file as the "last build time", **dobi** is able to skip
   the build of an image in many cases.

   If Docker adds a "last modified" time to the image data, **dobi** will be able
   to use that time instead of tracking the time itself.


``:tag``
~~~~~~~~

Tag the image with all the tags in the **tags** field.

The ``:tag`` action always depends on the ``:build`` action for the image.

``:push``
~~~~~~~~~

Push the image tags to a registry.

The ``:push`` action always depends on the ``:tag`` action for the image.


``:pull``
~~~~~~~~~


.. note::

    This action is planned, but not implemented yet.


Pull the image from a registry.


``:remove``
~~~~~~~~~~~

:alias: ``:rm``

Remove all the image tags, and the image.


Job Tasks
---------

`job <./config.html#job>`_ resources have the following tasks:

``:run`` *(default)*
~~~~~~~~~~~~~~~~~~~~

Run a process in a container.

``:remove``
~~~~~~~~~~~

:alias: ``:rm``

Remove the container (if it exists), and remove the artifact (if one is defined).

Mount Tasks
-----------

`mount <./config.html#mount>`_ resources have the following tasks:

``:create`` *(default)*
~~~~~~~~~~~~~~~~~~~~~~~

Create the host directory to be bind mounted, if it doesn't already exist.


``:remove``
~~~~~~~~~~~

:alias: ``:rm``

Does nothing. This action exists because all resources have have a remove task.

Alias Tasks
-----------

`alias <./config.html#alias>`_ resources have the following tasks:

``:run`` *(default)*
~~~~~~~~~~~~~~~~~~~~~

Run all the tasks in the list of tasks.

``:remove``
~~~~~~~~~~~

:alias: ``:rm``

Remove runs the remove task for all the resources in the task list in
reverse order.


Compose Tasks
-------------

`compose <./config.html#compose>`_ resources have the following tasks:

``:up`` *(default)*
~~~~~~~~~~~~~~~~~~~

Up runs ``docker-compose up -d`` with the files and project name from
the resource to create a new isolated environment.

``:down``
~~~~~~~~~

:alias: ``:rm``, ``:remove``

Down runs ``docker-compose down`` to remove all the containers and networks created
by Compose.

``:attach``
~~~~~~~~~~~

Attach runs ``docker-compose up`` and attaches to the logs.
