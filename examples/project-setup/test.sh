#!/usr/bin/env bash
# SUMMARY: Test project setup example
# LABELS:
# REPEAT:
# AUTHOR:
set -eu -o pipefail

cleanup() {
    echo "running cleanup"
    dobi autoclean
    rm -f .env
}
trap "cleanup" EXIT


expect script.exp
set -x
grep username=myuser .env
grep port=8889 .env
