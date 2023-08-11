#!/bin/bash

set -e

./app/tooling/wait-for-it/wait-for-it.sh docker_postgres:5432 -t 30
cd /app/sql/migrations
goose postgres $DB_URL up

exec "$@"