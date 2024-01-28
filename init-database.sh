#!/bin/bash
set - e psql - v ON_ERROR_STOP = 1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE IF NOT EXISTS weatherunits (
        date date,
        humidity float8,
        temperature float8,
    );
EOSQL