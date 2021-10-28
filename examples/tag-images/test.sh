#!/usr/bin/env bash
# SUMMARY: Test example of tagging images
# LABELS:
# REPEAT:
# AUTHOR:
set -eu -o pipefail

remove_all_tags() {
    local image=$1
    docker images -q "$image" | uniq | xargs docker rmi -f
}

cleanup() {
    echo "running cleanup"
    rm -f output
    remove_all_tags example/tagged-app
    remove_all_tags example/tagged-db
    dobi autoclean
}
trap "cleanup" EXIT

export APP_VERSION=testing
dobi tag-images

echo "image=app tags 5 images"
[[ $(docker images -q example/tagged-app | wc -l) == 5 ]]

echo "image=db tags 5 images"
[[ $(docker images -q example/tagged-db | wc -l) == 5 ]]
