#!/bin/sh
set -e

MIGRATION_FLAG="/build/migrations_done"

echo "Waiting for postgres to start..."
until pg_isready -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER"; do
  sleep 2
done

echo "Postgres is ready"

if [ ! -f "$MIGRATION_FLAG" ]; then
  echo "Running database migrations..."
  /build/migrate/migrate
  touch "$MIGRATION_FLAG"
  echo "Migrations completed"
else
  echo "Skipping migrations"
fi

echo "Starting the app..."
exec /build/fuel-economy-go
