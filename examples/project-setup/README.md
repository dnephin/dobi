
# Example - Project setup

Some projects need customization before they can run. Often these
requirements are documented in a README, and require a user to copy a
default file into place and edit it manually.  With `dobi` these steps can be
automated and enforced as a dependency of other tasks.

This example demonstrates creating a `.env` file from user input and using the
`setup` task as a dependency for other tasks.

To run a project task:

    dobi app

The first time you run the app, you'll be prompted for some options. Once the
`.env` file exists you won't be prompted anymore.
