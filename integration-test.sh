#!/bin/bash

# build
docker build . -t blockexchange

# setup
docker network create blockexchange
docker run --name blockexchange_pg --rm \
 -e POSTGRES_PASSWORD=enter \
 --network blockexchange \
 postgres &

# wait for psql start
sleep 1
until docker exec -i blockexchange_pg pg_isready
do
 sleep 1
done


docker run --name blockexchange_server --rm \
 -e PGUSER=postgres \
 -e PGPASSWORD=enter \
 -e PGHOST=blockexchange_pg \
 -e PGDATABASE=postgres \
 -e PGPORT=5432 \
 -e BLOCKEXCHANGE_KEY=blah \
 --network blockexchange \
 blockexchange &

function cleanup {
	# cleanup
	docker stop blockexchange_server
	docker stop blockexchange_pg
  docker network rm blockexchange
}

trap cleanup EXIT

# wait for nodejs start
sleep 1
until docker exec -i blockexchange_server curl -v http://127.0.0.1:8080/
do
  sleep 1
done

git clone https://github.com/blockexchange/blockexchange.git

CFG=/tmp/minetest.conf
MTDIR=/tmp/mt
WORLDDIR=${MTDIR}/worlds/world

cat <<EOF > ${CFG}
 blockexchange.url = http://blockexchange_server:8080
 secure.http_mods = blockexchange
EOF

mkdir -p ${WORLDDIR}
chmod 777 ${MTDIR} -R
docker run --rm -i \
	-v ${CFG}:/etc/minetest/minetest.conf:ro \
	-v ${MTDIR}:/var/lib/minetest/.minetest \
  -v $(pwd)/blockexchange:/var/lib/minetest/.minetest/worlds/world/worldmods/blockexchange \
  -v $(pwd)/blockexchange/test/test_mod:/var/lib/minetest/.minetest/worlds/world/worldmods/blockexchange_test \
  --network blockexchange \
	registry.gitlab.com/minetest/minetest/server:5.2.0

test -f ${WORLDDIR}/integration_test.json && exit 0 || exit 1


echo "Test complete!"
