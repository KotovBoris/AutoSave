#!/bin/sh

# Health check for the API service
curl -f http://localhost:${APP_PORT:-8080}/health || exit 1
