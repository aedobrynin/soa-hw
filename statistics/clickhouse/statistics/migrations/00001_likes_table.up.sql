CREATE TABLE IF NOT EXISTS likes
(
    user_id String,
    post_id String
)
ENGINE = MergeTree
PRIMARY KEY (user_id, post_id);

CREATE TABLE IF NOT EXISTS likes_queue
(
    user_id String,
    post_id String  
)
ENGINE = Kafka
SETTINGS kafka_broker_list = 'kafka:9092',
        kafka_topic_list = 'likes',
        kafka_group_name = 'likes_consumer_clickhouse',
        kafka_format = 'JSON',
        kafka_max_block_size = 1048576;

CREATE MATERIALIZED VIEW IF NOT EXISTS likes_queue_mv TO likes AS
SELECT user_id, post_id
FROM likes_queue;
