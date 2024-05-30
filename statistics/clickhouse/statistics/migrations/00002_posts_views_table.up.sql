CREATE TABLE IF NOT EXISTS posts_views
(
    user_id String,
    post_id String
)
ENGINE = MergeTree
ORDER BY post_id;

CREATE TABLE IF NOT EXISTS posts_views_queue
(
    user_id String,
    post_id String  
)
ENGINE = Kafka
SETTINGS kafka_broker_list = 'kafka:9092',
        kafka_topic_list = 'posts_views',
        kafka_group_name = 'posts_views_consumer_clickhouse',
        kafka_format = 'JSON',
        kafka_max_block_size = 1048576;
--- TODO: myb use kafka_skip_broken_messages?

CREATE MATERIALIZED VIEW IF NOT EXISTS posts_views_queue_mv TO posts_views AS
SELECT user_id, post_id
FROM posts_views_queue;
