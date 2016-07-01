Configuration File
==================

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

If an image depends on a run resource, the run resource will be executed first.
The run resource must exit before the image resource will be run.

An image resource can not depend on a mount.


Run
---
A run resource runs a process in a container.

Each run resource uses an image defined by an image resource.  By default, a
run resource is never considered up-to-date, it will always run.  If a run
resource has an ``artifact`` property, which is a path to a local file, the
last modified time of that file will be used. A run resource is considered
up-to-date if the modified time of the ``artifact`` is more recent then:

* the created time of the image it uses
* the last modified time of all files in any mounts used by the resource


The image resource used by a run resource is automatically added
as an implicit dependency.

If a run resource depends on another run resource, the dependency will be run first.

If a run resource uses ``mounts``, the mounts resources will be run first.

If a run resource uses a ``network`` resource it will be run first and the run
resource will join the default network for the environment.


Mount
-----
A mount resource creates a host or named mount. If the mount already exists
the resource is a no-op.

A mount can not depend on any resource.


Task Aliases
------------
A task alias is a list of other resource names. It will run each resource in the
order they are listed.


Service
-------

.. note::

    This resource hasn't been implemented yet.

A service resource runs a service using an image. The service is kept running
for the duration of the execution, and is shutdown when all other resources
are complete.

A service resource may depend on images, mounts, or run resources.


Meta
----
The meta section is the only section in the configuration that does not describe
a resource. It contains meta configuration details:

 * ``default`` - the default resource to run when no names are given on the
   command line
