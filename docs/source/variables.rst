Variables
=========

Some fields in ``dobi.yaml`` support variable interpolation from values provided
by ``dobi``.  Variables name are wraped in braces, for example ``{env.USER}`` would
inject the value of the ``$USER`` environment variable.

Supported variables
-------------------

The following variables are made avariables:

* ``env.<environment variable>`` - a value of an environment variable
* ``git.sha`` - the git sha (with -dirty if there are un-committed changes)
* ``git.branch`` - the git branch name
* ``time.<format>`` - a date or time using `fmtdate
  <https://github.com/metakeule/fmtdate#placeholders>`_
* ``unique`` - a unique execution id generate from the project name and exec id
* ``exec-id`` - an execution id (without project name)
* ``project`` - the project name


Supported fields
----------------

The following resource fields support variables:

* ``run.env``
* ``image.tag``
* ``compose.files``
* ``compose.project``
