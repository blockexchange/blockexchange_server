#!/bin/sh
cat dump.sql | docker exec -i blockexchange_pg psql -U postgres
