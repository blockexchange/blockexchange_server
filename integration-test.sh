#!/bin/bash

# build
docker build . -t blockexchange

# setup
docker run --name blockexchange_pg --rm \
 -e POSTGRES_PASSWORD=enter \
 --network host \
 postgres &

bash -c 'while !</dev/tcp/localhost/5432; do sleep 1; done;'

docker run --name blockexchange_server --rm \
 -e PGUSER=postgres \
 -e PGPASSWORD=enter \
 -e PGHOST=127.0.0.1 \
 -e PGDATABASE=postgres \
 -e PGPORT=5432 \
 --network host \
 blockexchange &

function cleanup {
	# cleanup
	docker stop blockexchange_server
	docker stop blockexchange_pg
}

trap cleanup EXIT

bash -c 'while !</dev/tcp/localhost/8080; do sleep 1; done;'

# TODO download mod and start minetest

echo "Test complete!"
