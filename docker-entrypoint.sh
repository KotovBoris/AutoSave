#!/bin/sh

set -e

echo "Waiting for postgres..."
while ! nc -z ${DB_HOST} ${DB_PORT}; do
  sleep 1
done
echo "PostgreSQL started"

echo "Running migrations..."
migrate -path /migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

echo "Starting application..."
exec "$@"
