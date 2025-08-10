#!/bin/bash
DBSTRING="host=$DB_ADDR user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=disable"
goose postgres "$DBSTRING" up