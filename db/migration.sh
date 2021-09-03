#!/bin/bash

export MIGRATION_DIR=./db/migration
export DB_PORT="5432"

export DB_NAME="purchases"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_HOST="localhost"

export DB_SSL=disable
export PG_DSN="host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} user=${DB_USER} sslmode=${DB_SSL} password=${DB_PASSWORD}"

if [ "$1" = "--dryrun" ]; then
./bin/goose -dir ${MIGRATION_DIR} postgres "${PG_DSN}" status -v
else
./bin/goose -dir ${MIGRATION_DIR} postgres "${PG_DSN}" up -v
fi
