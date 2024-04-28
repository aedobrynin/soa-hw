#!/bin/bash

kafka-topics.sh --create --if-not-exists --topic $KAFKA_POSTS_LIKES_TOPIC_NAME --bootstrap-server kafka:9092
kafka-topics.sh --create --if-not-exists --topic $KAFKA_POSTS_VIEWS_TOPIC_NAME --bootstrap-server kafka:9092
