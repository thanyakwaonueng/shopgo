#!/bin/sh
[ -e "$PWD"/.env.dev ] && . "$PWD"/.env.dev

docker compose down -v
docker compose build
docker compose up -d

