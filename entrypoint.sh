#!/bin/sh
set -e

echo "Waiting for postgres to start..."
until pg_isready -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER"; do
  sleep 2
done

echo "Postgres is ready"

echo "Running database migrations..."
/build/migrate/migrate

echo "Starting the app..."
exec /build/fuel-economy-go