CREATE SCHEMA IF NOT EXISTS posts;

CREATE TABLE IF NOT EXISTS posts.posts (
   id            UUID   NOT NULL PRIMARY KEY,
   author_id     UUID   NOT NULL
   content       TEXT   NOT NULL,
   created_ts    TIMESTAMPTZ NOT NULL DEFAULT now(),
   updates_ts    TIMESTAMPTZ NOT NULL DEFAULT now(),
);
