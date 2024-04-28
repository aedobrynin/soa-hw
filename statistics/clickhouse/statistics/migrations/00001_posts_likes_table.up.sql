CREATE TABLE IF NOT EXISTS posts_likes
(
    user_id String,
    post_id String
)
ENGINE = MergeTree
PRIMARY KEY (user_id, post_id);

CREATE TABLE IF NOT EXISTS posts_likes_queue
(
    user_id String,
    post_id String  
)
ENGINE = Kafka
SETTINGS kafka_broker_list = 'kafka:9092',
        kafka_topic_list = 'posts_likes',
        kafka_group_name = 'posts_likes_consumer_clickhouse',
        kafka_format = 'JSON',
        kafka_max_block_size = 1048576;
--- TODO: myb use kafka_skip_broken_messages?

CREATE MATERIALIZED VIEW IF NOT EXISTS posts_likes_queue_mv TO posts_likes AS
SELECT user_id, post_id
FROM posts_likes_queue;
