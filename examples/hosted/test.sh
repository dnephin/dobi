#!/usr/bin/env bash
# SUMMARY: Test hosted environment
# LABELS:
# REPEAT:
# AUTHOR:
set -eu -o pipefail

cleanup() {
    echo "running cleanup"
    dobi autoclean
    rm -f output dockerignored-file uncommitted-file
}
trap "cleanup" EXIT

dobi run-dist | tee output
docker images

echo "image=dist-img creates an image"
docker inspect --type image example/hosted-hello:root

echo "job=run-dist outputs hello"
grep '^Hello, world!' output

echo "touch committed-file"
touch committed-file

echo "ensure image=dist-img stays fresh"
dist_img_id1=$(docker inspect --type image example/hosted-hello:root --format='{{ .Id }}')
dobi dist-img:build
dist_img_id2=$(docker inspect --type image example/hosted-hello:root --format='{{ .Id }}')
if [[ "$dist_img_id1" != "$dist_img_id2" ]]; then
    echo "ERROR: dist-img was incorrectly rebuilt"
    exit 1
fi

echo "touch uncommitted-file"
touch uncommitted-file

echo "ensure image=dist-img is stale"
dobi dist-img:build
dist_img_id3=$(docker inspect --type image example/hosted-hello:root --format='{{ .Id }}')
if [[ "$dist_img_id2" == "$dist_img_id3" ]]; then
    echo "ERROR: dist-img was not rebuilt when it should have been"
    exit 1
fi
