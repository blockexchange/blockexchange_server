#!/bin/sh

export PGUSER=postgres
export PGPASSWORD=enter
export PGHOST=127.0.0.1
export PGDATABASE=postgres
export PGPORT=5432

npm run watch
