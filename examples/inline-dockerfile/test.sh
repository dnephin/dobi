#!/usr/bin/env bash
# SUMMARY: Test inline Dockerfile example
# LABELS:
# REPEAT:
# AUTHOR:
set -eu -o pipefail

cleanup() {
    echo "running cleanup"
    dobi autoclean
}
trap "cleanup" EXIT

dobi tree
