#!/usr/bin/env bash
# SUMMARY: Test example of capturing staout as env var
# LABELS:
# REPEAT:
# AUTHOR:
set -eu -o pipefail

cleanup() {
    echo "running cleanup"
    dobi autoclean
}
trap "cleanup" EXIT

dobi dist

echo "image tagged with version"
docker inspect repo/myapp:3.4.5
