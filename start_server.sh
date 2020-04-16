#!/bin/sh

export PGUSER=postgres
export PGPASSWORD=enter
export PGHOST=127.0.0.1
export PGDATABASE=postgres
export PGPORT=5432
export BLOCKEXCHANGE_NAME=dev-exchange
export BLOCKEXCHANGE_OWNER=nobody
export BLOCKEXCHANGE_KEY=abcdefSecretKey

#export MATOMO_URL=http://127.0.0.1:8080/
#export MATOMO_ID=666

npm run watch
