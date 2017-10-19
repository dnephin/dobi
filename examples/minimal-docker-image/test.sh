#!/usr/bin/env bash
# SUMMARY: Test build minimal docker images
# LABELS:
# REPEAT:
# AUTHOR:
set -eu -o pipefail

cleanup() {
    echo "running cleanup"
    dobi autoclean
    rm -f output
}
trap "cleanup" EXIT

dobi run-dist | tee output
docker images

echo "image=builder creates an image"
docker inspect --type image minimal-dev:example-hello-root

echo "job=binary creates the binary"
ls ./dist/bin/hello

echo "image=dist-img creates an image"
docker inspect --type image example/hello:root

echo "job=run-dist outputs hello"
grep '^Hello, world!' output
