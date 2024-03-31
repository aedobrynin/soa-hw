CREATE SCHEMA IF NOT EXISTS posts;

CREATE TABLE IF NOT EXISTS posts.posts (
   id            UUID   NOT NULL PRIMARY KEY,
   author_id     UUID   NOT NULL,
   content       TEXT   NOT NULL,
   created_ts    TIMESTAMPTZ NOT NULL DEFAULT now(),
   updated_ts    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION trigger_set_updated_ts()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_ts = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER posts_posts_set_updated_ts
   BEFORE UPDATE ON posts.posts
   FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_ts();
