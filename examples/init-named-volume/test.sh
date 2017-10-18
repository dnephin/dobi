#!/usr/bin/env bash
# SUMMARY: Test initialize named volume example
# LABELS:
# REPEAT:
# AUTHOR:
set -eu -o pipefail

cleanup() {
    echo "running cleanup"
    dobi autoclean
}
trap "cleanup" EXIT

dobi populate

echo "the volume contains the correct file"
docker volume inspect example-init-volume-data
docker run --rm \
    -v example-init-volume-data:/data \
    alpine:3.6 ls -1 data | \
    grep "newfile"

dobi view | grep newfile
