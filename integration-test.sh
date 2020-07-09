#!/bin/bash

set -e

# build
docker build . -t blockexchange

# setup
docker network create blockexchange
docker run --name blockexchange_pg --rm -d \
 -e POSTGRES_PASSWORD=enter \
 --network blockexchange \
 postgres &

# wait for psql start
sleep 1
until docker exec -i blockexchange_pg pg_isready
do
 sleep 1
done


docker run --name blockexchange_server --rm -d \
 -e PGUSER=postgres \
 -e PGPASSWORD=enter \
 -e PGHOST=blockexchange_pg \
 -e PGDATABASE=postgres \
 -e PGPORT=5432 \
 -e BLOCKEXCHANGE_KEY=blah \
 --network blockexchange \
 blockexchange

function cleanup {
	# cleanup
	docker stop blockexchange_server
  docker stop blockexchange_pg
  docker stop blockexchange_selenium
  docker network rm blockexchange
}

trap cleanup EXIT

# wait for nodejs start
until docker exec -i blockexchange_server curl -v http://127.0.0.1:8080/
do
  sleep 1
done

### Test web client

docker build e2e -t bx_e2e

docker run --name blockexchange_selenium --rm -d \
 --network blockexchange \
 selenium/standalone-chrome

until docker exec -i blockexchange_server curl -v http://blockexchange_selenium:4444/
do
  sleep 1
done

docker run --name blockexchange_e2e --network blockexchange --rm bx_e2e

### Test minetest client

git clone https://github.com/blockexchange/blockexchange.git || echo ok

CFG=/tmp/minetest.conf
MTDIR=/tmp/mt
WORLDDIR=${MTDIR}/worlds/world

cat <<EOF > ${CFG}
 blockexchange.url = http://blockexchange_server:8080
 secure.http_mods = blockexchange
EOF

mkdir -p ${WORLDDIR}
chmod 777 ${MTDIR} -R || echo ok
docker run --rm -i \
	-v ${CFG}:/etc/minetest/minetest.conf:ro \
	-v ${MTDIR}:/var/lib/minetest/.minetest \
  -v $(pwd)/blockexchange:/var/lib/minetest/.minetest/worlds/world/worldmods/blockexchange \
  -v $(pwd)/blockexchange/test/test_mod:/var/lib/minetest/.minetest/worlds/world/worldmods/blockexchange_test \
  --network blockexchange \
	registry.gitlab.com/minetest/minetest/server:5.2.0

test -f ${WORLDDIR}/integration_test.json && exit 0 || exit 1


echo "Test complete!"
