#!/bin/bash
set -eux

mkdir -p "$PGDATA"
chmod 700 "$PGDATA"
chown -R postgres "$PGDATA"

chmod g+s /run/postgresql
chown -R postgres /run/postgresql

gosu postgres initdb
authMethod=trust
{ echo; echo "host all all 0.0.0.0/0 $authMethod"; } >> "$PGDATA/pg_hba.conf"

gosu postgres pg_ctl -D "$PGDATA" \
    -o "-c listen_addresses='localhost'" \
    -w start

gosu postgres psql < export.sql

gosu postgres pg_ctl -D "$PGDATA" -w stop
