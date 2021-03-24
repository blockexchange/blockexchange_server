minetest blockexchange server software

![](https://github.com/blockexchange/blockexchange_server/workflows/docker/badge.svg)
![](https://github.com/blockexchange/blockexchange_server/workflows/test/badge.svg)
![](https://github.com/blockexchange/blockexchange_server/workflows/jshint/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/blockexchange/blockexchange_server/badge.svg)](https://coveralls.io/github/blockexchange/blockexchange_server)

![Docker Pulls](https://img.shields.io/docker/pulls/blockexchange/blockexchange)
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/blockexchange/blockexchange)
![Discord](https://img.shields.io/discord/736160589130235965)

Discord-server: https://discord.gg/BEf8hGEQz9

# Overview

Corresponding mod and more information: https://github.com/blockexchange/blockexchange

# Running

## Requirements

* Running postgres database

## Environment variables

Required:
* **PGUSER**
* **PGPASSWORD**
* **PGHOST**
* **PGPORT**
* **PGDATABASE**
* **BLOCKEXCHANGE_KEY** private key to sign the json web tokens with

Optional:
* **DISCORD_SCHEMA_FEED_URL** discord webhook for the "new schema" feed
* **BASE_URL** Application base url, used for redirects (without trailing slash)
* **GITHUB_APP_SECRET** Github app secret key
* **GITHUB_APP_ID** Github app ID
* **DISCORD_APP_SECRET** Discord app secret key
* **DISCORD_APP_ID** Discord app ID
* **MESEHUB_APP_SECRET** Mesehub app secret key
* **MESEHUB_APP_ID** Mesehub app ID
* **MATOMO_URL** Matomo tracker url
* **MATOMO_ID** Matomo tracker id
* **REDIS_HOST** redis host
* **REDIS_PORT** redis port

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

Prerequisites:
* docker
* docker-compose

```
# start the postgres database
docker-compose up -d postgres

# start the blockexchange server
docker-compose up -d blockexchange

# start the minetest server
docker-compose up -d minetest
```

The blockexchange UI is at http://127.0.0.1:8080 and the minetest server is reachable via 127.0.0.1:30000

# Database

Database model:

<img src="./doc/database.png"/>

## Serialized mapblock format

The on-disk format

```lua
-- examplary metadata
metadata = {
	size = {x=16,y=16,z=16},
	node_mapping = {
		["air"] = 126,
		["default:stone"] = 127
	},
	-- block metadata
	metadata = {
		meta = {
			["(0,0,0)"] = {
				inventory = {},
				fields = {}
			}
		},
		timers = {
			["(0,0,0)"] = {
				timeout = 2.0,
				elapsed = 1.4
			}
		}
	}
}

-- node_id's (2 bytes) / param1 (1 byte) / param2 (1 byte)
data = {
	0,0,0,0
	-- etc
}

-- database format
serialized_metadata = minetest.compress(minetest.write_json(metadata), "deflate")
serialized_data = minetest.compress(data, "deflate")
```

Serialized example as json, byte-array is encoded as base64 over the wire:
```json
{
	"data":"eJztwQENAAAMAqAHekijm8MNyEdVVVVVVVVVHX8AAAAAAAAAwLwCfjrAlw",
	"metadata":"eJw1ylsKgCAQRuG9/M8RBdLDbCYGnUJQk9Lognsvid4+OOeGl8SGE4NuJOtl3UAhO1cahMXI6DlGG+aajUycXSLNu4xWC0iptnvHzV5ShwPUD23X4PxxfSjlAYnOIYs",
	"offset_x":0.0,
	"offset_y":0.0,
	"offset_z":0.0,
	"schema_id":9.0
}
```

# License

Code: MIT

## Other assets

* `public/pics/default_mese_crystal.png` CC BY-SA 3.0 https://github.com/minetest/minetest_game
