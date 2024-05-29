DROP TABLE IF EXISTS posts_likes_queue;

CREATE TABLE IF NOT EXISTS posts_likes_queue
(
    user_id String,
    post_id String,
    post_author_id String
)
ENGINE = Kafka
SETTINGS kafka_broker_list = 'kafka:9092',
        kafka_topic_list = 'posts_likes',
        kafka_group_name = 'posts_likes_consumer_clickhouse',
        kafka_format = 'JSON',
        kafka_max_block_size = 1048576;
--- TODO: myb use kafka_skip_broken_messages?

CREATE TABLE IF NOT EXISTS posts_likes_indexed_by_post_author
(
    user_id String,
    post_id String,
    post_author_id String
)
ENGINE = MergeTree
ORDER BY post_author_id;

CREATE MATERIALIZED VIEW IF NOT EXISTS posts_likes_queue_mv TO posts_likes_indexed_by_post_author AS
SELECT user_id, post_id, post_author_id
FROM posts_likes_queue;
