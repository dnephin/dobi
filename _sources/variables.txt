Config Variables
================

Some fields in ``dobi.yaml`` support variable interpolation from values provided
by **dobi** or the shell environment.  Variables are wrapped in braces, for example
``{env.USER}`` would inject the value of the ``$USER`` environment variable.

Format
------

Environment variables have the following format:


.. code-block:: default

    "{" [section.]variable[:default] "}"

**{}**
    All variables are wrapped in braces

**section**
    Some variables are grouped into sections (like **git** or **env**)

**default**
    Variables can have default values. The value after the last colon is taken
    as the default value. An empty default value makes the variable act like an
    optional variable.

Example
~~~~~~~

Use the value of ``$VERSION`` from the host use the variable:

.. code-block:: none

    {env.VERSION}

Use a default value of ``v1.0``:

.. code-block:: none

    {env.VERSION:v1.0}

Use a variable with an empty default value as an optional value:

.. code-block:: none

    {env.VERSION:}


Supported Variables
-------------------

The supported variables are:

==================  ===========================================================
Variable            Description
==================  ===========================================================
``env.<variable>``  value of an environment variable
``exec-id``         execution id (without project name)

``fs.cwd``          current working directory
``fs.projectdir``   directory which contains the ``dobi.yaml``

``git.branch``      current git branch name
``git.sha``         current git sha
``git.short-sha``   first 10 characters of the current git sha
``project``         project name
``time.<format>``   a date or time using `fmtdate
                    <https://github.com/metakeule/fmtdate#placeholders>`_
                    (note: if your time format includes a ``:`` you must add
                    another ``:`` to the end of the format, otherwise the string
                    after the final ``:`` will be taken as the default value)
``unique``          a unique execution id generate from the project name and exec
                    id
``user.name``       username of the active user
``user.uid``        uid of the active user
``user.gid``        primary gid of the active user
``user.home``       home directory of the active user
``user.group``      primary group name of the active user
==================  ===========================================================


Config Fields
-------------

Variables are only supported in specific fields:

+----------------+-----------------------------------------------------------+
| Resource       | Field                                                     |
+================+===========================================================+
| env            | files                                                     |
|                +-----------------------------------------------------------+
|                | variables                                                 |
+----------------+-----------------------------------------------------------+
| job            | env                                                       |
|                +-----------------------------------------------------------+
|                | user                                                      |
|                +-----------------------------------------------------------+
|                | net-mode                                                  |
|                +-----------------------------------------------------------+
|                | working-dir                                               |
+----------------+-----------------------------------------------------------+
| image          | tag                                                       |
|                +-----------------------------------------------------------+
|                | image                                                     |
|                +-----------------------------------------------------------+
|                | args                                                      |
+----------------+-----------------------------------------------------------+
| compose        | files                                                     |
|                +-----------------------------------------------------------+
|                | project                                                   |
+----------------+-----------------------------------------------------------+
| mount          | path                                                      |
|                +-----------------------------------------------------------+
|                | bind                                                      |
+----------------+-----------------------------------------------------------+
| meta           | exec-id                                                   |
+----------------+-----------------------------------------------------------+
