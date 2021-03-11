#!/bin/sh

export PGUSER=postgres
export PGPASSWORD=enter
export PGHOST=127.0.0.1
export PGDATABASE=postgres
export PGPORT=5432
export BLOCKEXCHANGE_NAME=My-Blockexchange
export BLOCKEXCHANGE_OWNER=yourname
export BLOCKEXCHANGE_KEY=abcdefSecretKey
export LOGLEVEL=debug
export BASE_URL=http://localhost:8080

test -f secret-env.sh && . ./secret-env.sh

go run .
