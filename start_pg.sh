#!/bin/sh
sudo docker run --name blockexchange_pg --rm -it -e POSTGRES_PASSWORD=enter -p 5432:5432 postgres
