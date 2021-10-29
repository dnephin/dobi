#!/usr/bin/env bash
# SUMMARY: Test example of capturing staout as env var
# LABELS:
# REPEAT:
# AUTHOR:
set -eu -o pipefail

cleanup() {
    echo "running cleanup"
    # These changes break autoclean. 
    # See https://github.com/dnephin/dobi/issues/227
    # and https://github.com/dnephin/dobi/issues/228
    dobi autoclean || true
}
trap "cleanup" EXIT

dobi dist

echo "image tagged with version"
docker inspect repo/myapp:3.4.5
