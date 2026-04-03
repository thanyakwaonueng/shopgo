#!/bin/sh
# 1. Load env
[ -e "$PWD"/.env.dev ] && . "$PWD"/.env.dev

# 2. Clean start
docker compose down -v
docker compose build

# 3. Start ONLY the database first
echo "Starting Database..."
docker compose up -d db

# 4. Wait for Postgres to be ready
echo "Waiting for Postgres to initialize (Asia/Bangkok)..."
until docker compose exec db pg_isready -U postgres > /dev/null 2>&1; do
  echo "Postgres is unavailable - sleeping"
  sleep 2
done

echo "Postgres is READY! Starting API..."

# 5. Now start the API
docker compose up -d api

# 6. Show logs so you can see the "Starting ShopGo API..." message
docker compose logs -f api
