#!/bin/sh
cat dump.sql | sudo docker exec -i blockexchange_pg psql -U postgres
