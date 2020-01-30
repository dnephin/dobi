
dobi
====

A build automation tool for Docker applications.

Keep your project tasks organized, portable, repeatable, and fast with ``dobi``.
Define the resources and tasks required to build, test, and release your project in
a ``dobi.yaml`` and run them with ``dobi TASK``.

See `Getting Started <https://dnephin.github.io/dobi/>`_

.. image:: https://img.shields.io/github/release/dnephin/dobi.svg?maxAge=7200
    :target: https://github.com/dnephin/dobi/releases/latest

.. image:: https://circleci.com/gh/dnephin/dobi/tree/master.svg?style=shield
    :target: https://circleci.com/gh/dnephin/dobi/tree/master

.. image:: https://goreportcard.com/badge/github.com/dnephin/dobi
    :target: https://goreportcard.com/report/github.com/dnephin/dobi

.. image:: https://badges.gitter.im/dnephin/dobi.png
    :target: https://gitter.im/dnephin-dobi/Lobby

Features
--------

Key features of ``dobi``:

* **optimal** - tasks are only run when the resource is stale. If a resource
  hasn't changed the task is skipped.
* **portable** - all tasks run in a container, so developers are free to use
  different operating systems and environments.
* **repeatable** - tasks are defined in a ``dobi.yaml`` so new contributers can
  get started quickly, and a task will always produce the same results.
  Variables are supported, but must be explicitly defined, so there's no hidden
  environment variables that could change the behaviour of a task.
* **organized** - tasks can be chained together using an ``alias`` resource to
  produce entire workflows like ``test`` or ``release``, which may involve
  multiple independent tasks.
* **dependencies** - tasks can depend on other tasks using ``depends``. When a
  task is run, its dependencies are checked first, and run if they are stale.


Install
-------

The one liner::

    curl -L -o /usr/local/bin/dobi "https://github.com/dnephin/dobi/releases/download/v0.13.0/dobi-$(uname -s)"; chmod +x /usr/local/bin/dobi

For a Windows binary, and more install options, see `Install <https://dnephin.github.io/dobi/install.html>`_

Examples
--------

See `Examples <https://dnephin.github.io/dobi/examples.html>`_

Documentation
-------------

See `Documentation <https://dnephin.github.io/dobi/>`_


Contributing
------------

``dobi`` is still in early development. If you'd like to contribute, please open
an issue, or find an existing issue, and leave a comment saying you're working
on a feature.
