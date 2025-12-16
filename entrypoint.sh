#!/bin/sh
set -e

echo "========================================="
echo "=== Waiting for PostgreSQL..."
echo "========================================="

# Wait for PostgreSQL to be ready
until nc -z "${DB_HOST}" 5432; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done

echo "✓ PostgreSQL is up and accepting connections"
echo ""

echo "========================================="
echo "=== Running database migrations..."
echo "========================================="
if /app/bin migrate up; then
    echo "✓ Migrations completed successfully!"
else
    echo "✗ Migration failed!"
    exit 1
fi
echo "========================================="
echo "=== Starting application..."
echo "========================================="

exec "$@"
