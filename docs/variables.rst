Variables
=========

Some fields in the ``dobi.yaml`` file accept variables provided by ``dobi``
to alter the behaviour of a task. Variables are in the form ``{key}``. The
following variables are made avariables:

* ``env.<environment variable>`` - a value of an environment variable
* ``git.sha`` - the git sha (with -dirty if there are un-committed changes)
* ``git.branch`` - the git branch name
* ``unique`` - a unique execution id generate from the project name and exec id
* ``exec-id`` - an execution id (without project name)
* ``project`` - the project name
