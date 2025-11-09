#!/bin/bash
# run-migrations.sh

DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

echo "Running database migrations..."
echo "Database URL: postgres://${DB_USER}:****@${DB_HOST}:${DB_PORT}/${DB_NAME}"

for file in migrations/*.up.sql; do
    if [ -f "$file" ]; then
        echo "Applying migration: $(basename $file)"
        psql "${DB_URL}" -f "$file"
        if [ $? -ne 0 ]; then
            echo "Migration failed: $(basename $file)"
            exit 1
        fi
    fi
done

echo "Migrations completed successfully!"
