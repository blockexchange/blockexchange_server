minetest blockexchange server software

![](https://github.com/blockexchange/blockexchange_server/workflows/integration-test/badge.svg)
![](https://github.com/blockexchange/blockexchange_server/workflows/docker/badge.svg)
![](https://github.com/blockexchange/blockexchange_server/workflows/jshint_frontend/badge.svg)
![](https://github.com/blockexchange/blockexchange_server/workflows/jshint_backend/badge.svg)
![](https://github.com/blockexchange/blockexchange_server/workflows/test/badge.svg)


State: **WIP**

Docker: https://hub.docker.com/r/blockexchange/blockexchange

# Running

## Requirements

* Running postgres database

## Environment variables

* **PGUSER**
* **PGPASSWORD**
* **PGHOST**
* **PGPORT**
* **PGDATABASE**
* **BLOCKEXCHANGE_KEY** private key to sign the json web tokens with

## Docker usage

This example shows a simple throw-away setup.

Start a postgres server:
```bash
docker run --rm -it \
 -e POSTGRES_PASSWORD=enter \
 --network host \
 postgres
```

Start the server part:
```bash
docker run -it --rm \
 -e PGUSER=postgres \
 -e PGPASSWORD=enter \
 -e PGHOST=127.0.0.1 \
 -e PGDATABASE=postgres \
 -e PGPORT=5432 \
 -e BLOCKEXCHANGE_KEY=blah \
 --network host \
 blockexchange/blockexchange
```

Go to http://127.0.0.1:8080
