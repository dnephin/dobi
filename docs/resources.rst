Resources
=========

Every section in a ``dobi.yaml`` configuration file defines a resource. Each
resource must be one of the following types.

.. contents::
    :backlinks: none
    :depth: 2


Image
-----
An image resource builds an image from a Dockerfile, or pulls an image from a
registry.

An image is considered up-to-date if all files in the build context have a
modified time older than the created time of the current image.

If an image depends on another image resource, the dependency will be built
first (if it is not up-to-date).

If an image depends on a command, the command will be run first. The
command resource must exit before the image resource will be run.

An image resource can not depend on a volume.


Command
-------
A command resource runs a process in a container.

Each command uses an image defined by an image resource.  By default, a command
is never considered up-to-date, it will always run.  If a container resource has
an ``artifact`` property, which is a path to a local file, the last modified
time of that file will be used. A command resource is considered up-to-date if
the modified time of the ``artifact`` is more recent then:

* the created time of the image it uses
* the last modified time of all files in any volumes used by the resource


The image resource used by a command resource is automatically added
as an implicit dependency of the command.

If a container depends on another container, the container will be run first.

If a uses ``volumes``, the volumes resources will be run first.

If a container uses a ``network`` resource it will be run first and the container
will join the default network for the environment.


Volume
------
A volume resource creates a host or named volume. If the volume already exists
the resource is a no-op.

A volume can not depend on any resource.


Environment
-----------

.. note::

    This resource hasn't been implemented yet.

An environment resource runs multiple containers as defined by a Compose file.

An environment may depend on images, volumes, or commands.
