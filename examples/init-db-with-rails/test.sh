#!/usr/bin/env bash
# SUMMARY: Test initialize a database
# LABELS:
# REPEAT:
# AUTHOR:
set -eu -o pipefail

cleanup() {
    echo "running cleanup"
    rm db/export.sql
    dobi autoclean
}
trap "cleanup" EXIT

dobi dev-environment

echo "image=rails creates an image with two tags"
docker inspect --type image \
    example/web:examplerailsdb \
    example/web:examplerailsdb-root

echo "compose=empty-db-env creates a container"
docker inspect --type container examplerailsdbexport_postgres_1

echo "job=export-models creates db/export.sql"
ls db/export.sql

echo "image=database-img creates an image with two tags"
docker inspect --type image \
    example/database:examplerailsdb \
    example/database:examplerailsdb-root
