#!/bin/bash

# Create database and user if not exists
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Create autosave user if not exists
    DO
    \$do\$
    BEGIN
        IF NOT EXISTS (
            SELECT FROM pg_catalog.pg_roles
            WHERE rolname = 'autosave'
        ) THEN
            CREATE ROLE autosave WITH LOGIN PASSWORD 'autosave_password';
        END IF;
    END
    \$do\$;

    -- Grant privileges
    GRANT ALL PRIVILEGES ON DATABASE autosave_db TO autosave;
    GRANT ALL ON SCHEMA public TO autosave;
EOSQL
