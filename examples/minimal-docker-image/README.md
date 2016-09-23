
# Example - Building minimal docker images

This example demonstrates how to use ``dobi`` to build a minimal docker image,
that does not contain build/compile dependencies. This is sometimes referred
to as "squashing" an image. Building a small image is desirable when you want to
distribute an application as a docker image.

To create the image run:

    dobi dist-img

Then run the image with:

    docker run -ti --rm example/hello:$USER

or

    dobi run-dist


To remove the binary and the image run:

    dobi dist-img:rm builder:rm
