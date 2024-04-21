CREATE SCHEMA IF NOT EXISTS core;

CREATE TABLE IF NOT EXISTS core.users (
   id            UUID   NOT NULL PRIMARY KEY,
   login         TEXT   NOT NULL UNIQUE,
   password_hash TEXT   NOT NULL,
   name          TEXT,
   surname       TEXT,
   email         TEXT,
   phone         TEXT
);

CREATE INDEX IF NOT EXISTS idx__core__login
   ON core.users(login);
