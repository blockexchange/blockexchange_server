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

# test
json='{
	"description":"my schema",
	"tags": ["test", "stuff"],
	"size_x": 20,
	"size_y": 10,
	"size_z": 10,
	"part_length": 10
}'

json=$(curl -X POST "http://127.0.0.1:8080/api/schema" --data "$json" -H "Content-Type: application/json")
schema_id=$(echo $json | jq -r .id)

test "$schema_id" = "1" || exit -1
test "$(echo $json | jq -r .total_size)" = "0" || exit -1
test "$(echo $json | jq -r .total_parts)" = "0" || exit -1

json="{
	\"schema_id\": ${schema_id},
	\"offset_x\": 0,
	\"offset_y\": 0,
	\"offset_z\": 0,
	\"data\": \"return {}\"
}"

json=$(curl -X POST "http://127.0.0.1:8080/api/schemapart" --data "$json" -H "Content-Type: application/json")

json="{
	\"schema_id\": ${schema_id},
	\"offset_x\": 10,
	\"offset_y\": 0,
	\"offset_z\": 0,
	\"data\": \"return {}\"
}"

json=$(curl -X POST "http://127.0.0.1:8080/api/schemapart" --data "$json" -H "Content-Type: application/json")

json='{
	"modname_count": {
		"default": 2000,
		"air": 50000
	}
}'
json=$(curl -X POST "http://127.0.0.1:8080/api/schema/${schema_id}/complete" --data "$json" -H "Content-Type: application/json")


json=$(curl "http://127.0.0.1:8080/api/schema/${schema_id}")
test "$(echo $json | jq -r .id)" = "1" || exit -1
test "$(echo $json | jq -r .complete)" = "true" || exit -1
test "$(echo $json | jq -r .size_x)" = "20" || exit -1
test "$(echo $json | jq -r .size_y)" = "10" || exit -1
test "$(echo $json | jq -r .size_z)" = "10" || exit -1
test "$(echo $json | jq -r .total_size)" != "0" || exit -1
test "$(echo $json | jq -r .total_parts)" != "0" || exit -1

echo "Test complete!"
