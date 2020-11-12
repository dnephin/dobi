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

.. include:: ../gen/config/image.rst


.. include:: ../gen/config/job.rst


.. include:: ../gen/config/mount.rst


.. include:: ../gen/config/alias.rst


.. include:: ../gen/config/compose.rst


.. include:: ../gen/config/env.rst


.. include:: ../gen/config/meta.rst


.. include:: ../gen/config/annotationFields.rst
