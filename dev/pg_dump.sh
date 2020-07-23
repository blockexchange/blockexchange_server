#!/bin/sh
docker exec -it blockexchange_pg pg_dump -U postgres > dump.sql
