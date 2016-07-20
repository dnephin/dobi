
Examples
========

A great way to get started with **dobi** is to look over some examples.


Full Configs
------------

* `dobi <https://github.com/dnephin/dobi/blob/master/dobi.yaml>`_ - a Golang
  command line tool (``dobi`` uses itself!). This example demonstrates:

  * ``dist-img`` - building a minimal image (that doesn't contain build dependencies)
    for distributing an application.
  * building multiple images for different tasks (``builder``, ``linter``, ``releaser``)
  * ``watch`` - watching code for changes and auto-running unit tests and
    ``docs-watch`` for watching docs changes and auto-building the docs.
  * ``shell`` (and ``docs-shell``) - getting an interactive shell that contains
    all the dependencies required to build or test the project.
  * combinding tasks with an alias so they can be run together (``test``, and
    ``all``).



Use Cases
---------

* `examples/minimal-docker-image
  <https://github.com/dnephin/dobi/blob/master/examples/minimal-docker-image/>`_
  - building a minimal docker image, that does not contain build/compile
  dependencies. This is sometimes referred to as "squashing" an image.
* `examples/tag-images
  <https://github.com/dnephin/dobi/blob/master/examples/tag-images/>`_
  - tag images with metadata (git sha, git branch, datetime, and version)
* `examples/project-setup
  <https://github.com/dnephin/dobi/blob/master/examples/project-setup/>`_
  - prompt users for configuration and generate a ``.env`` file if it doesn't
  exist.
* `examples/init-db-with-rails
  <https://github.com/dnephin/dobi/blob/master/examples/init-db-with-rails/>`_
  - load database tables and fixtures from rails models and create a database
  image and development environment.
