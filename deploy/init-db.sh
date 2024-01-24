#!/bin/bash

set -xe

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL

DROP TABLE IF EXISTS customers;

CREATE table customers(
    id text PRIMARY KEY,
    email text UNIQUE
);
EOSQL