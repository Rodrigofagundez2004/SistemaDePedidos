#!/bin/bash
# Solo corre en un volumen nuevo y solo si usás el perfil "db" (Postgres en contenedor).
set -e
for db in users_db products_db orders_db payments_db; do
  if ! psql -tAc "SELECT 1 FROM pg_database WHERE datname='$db'" --username "$POSTGRES_USER" --dbname postgres | grep -q 1; then
    psql -c "CREATE DATABASE $db" --username "$POSTGRES_USER" --dbname postgres
  fi
done
