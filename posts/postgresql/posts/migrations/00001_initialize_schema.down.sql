DROP TRIGGER IF EXISTS posts_posts_set_updated_ts on posts.posts;
DROP FUNCTION IF EXISTS trigger_set_updated_ts();
DROP TABLE IF EXISTS posts.posts;
DROP SCHEMA posts;
