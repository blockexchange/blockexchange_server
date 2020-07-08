minetest blockexchange server software

![](https://github.com/blockexchange/blockexchange_server/workflows/integration-test/badge.svg)
![](https://github.com/blockexchange/blockexchange_server/workflows/docker/badge.svg)
![](https://github.com/blockexchange/blockexchange_server/workflows/jshint_frontend/badge.svg)
![](https://github.com/blockexchange/blockexchange_server/workflows/jshint_backend/badge.svg)
![](https://github.com/blockexchange/blockexchange_server/workflows/test/badge.svg)

![Docker Pulls](https://img.shields.io/docker/pulls/blockexchange/blockexchange)
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/blockexchange/blockexchange)

# Overview

Corresponding mod and more information: https://github.com/blockexchange/blockexchange

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

Image: https://hub.docker.com/r/blockexchange/blockexchange

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

## docker-compose usage

A `docker-compose` example:

```yml
version: "2"

services:
 blockexchange:
  image: blockexchange/blockexchange
  restart: always
  depends_on:
   - postgres
  environment:
   - PGUSER=postgres
   - PGPASSWORD=enter
   - PGHOST=postgres
   - PGDATABASE=postgres
   - PGPORT=5432
   - BLOCKEXCHANGE_NAME=My-Blockexchange
   - BLOCKEXCHANGE_OWNER=yourname
  ports:
   - "8080:8080"

 postgres:
  image: postgres:12
  restart: always
  environment:
   POSTGRES_PASSWORD: enter
  volumes:
   - "./data/postgres:/var/lib/postgresql/data"
```


# Development

Web- and backend development

```
# start the postgres database
./dev/start_pg.sh

# start the nodejs server
./start_server.sh
```

Go to http://127.0.0.1:8080

# License

Code: MIT

Textures:
* public/textures/default*.png CC BY-SA 3.0 https://github.com/minetest/minetest_game/
