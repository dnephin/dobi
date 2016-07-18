Config Reference
================

Every section in a :file:`dobi.yaml` configuration file defines a resource (with the
exception of `meta`_, which is configuration for **dobi**).

Each section in the file has the following form:

.. code-block:: yaml

    type=name:
        field: value
        ...

Each resource must be one of the following resource types:

.. include:: gen/config/image.rst


.. include:: gen/config/run.rst

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


.. include:: gen/config/mount.rst


.. include:: gen/config/alias.rst


.. include:: gen/config/compose.rst


.. include:: gen/config/meta.rst
