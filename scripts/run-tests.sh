#!/bin/bash

# Make script exit on first error
set -e

# Start test database if not running
docker-compose up -d postgres_test

# Wait for database to be ready
echo "Waiting for test database to be ready..."
until docker-compose exec -T postgres_test pg_isready -U postgres; do
  sleep 1
done

# Run tests
GO_ENV=test go test -v ./...