# AutoSave Project

This repository contains the full monorepo for the AutoSave project, including the frontend and backend.

## ?? Structure

- `/frontend`: Contains the frontend application (Vue/React/etc.).
- `/backend`: Contains the Go backend application.

## ?? Backend

For instructions on how to run the backend, see [`backend/README.md`](backend/README.md).

### Quick Start (from root)

```bash
# Build and start all backend services in Docker
docker-compose up -d --build

# Apply database migrations
docker-compose --profile tools run --rm migrate
```

## ?? Frontend

(Instructions to be added here)
