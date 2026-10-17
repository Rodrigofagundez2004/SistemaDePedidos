-- Crea las 4 bases del proyecto en TU PostgreSQL.
-- Ejecutar una sola vez, conectado a la base "postgres":
--   psql -U postgres -h localhost -f docker/postgres/create-databases.sql
-- o pegalo en el Query Tool de pgAdmin (conectado a la base "postgres").
-- Si alguna ya existe, va a dar error en esa línea; el resto se crea igual.

CREATE DATABASE users_db;
CREATE DATABASE products_db;
CREATE DATABASE orders_db;
CREATE DATABASE payments_db;
