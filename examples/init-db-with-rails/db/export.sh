#!/bin/bash

set -eux

# Wait for database to be ready
timeout 30 bash -e <<EOF
while ! pg_isready -t 10 --user postgres --host postgres --dbname postgres; do
    sleep 1
done
EOF

bin/rails db:migrate
pg_dump --user postgres --host postgres postgres > /db/export.sql
