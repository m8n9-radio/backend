#!/bin/sh
set -e

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
