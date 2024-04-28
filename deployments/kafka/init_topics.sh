#!/bin/bash

kafka-topics.sh --create --if-not-exists --topic $KAFKA_LIKES_TOPIC_NAME --bootstrap-server kafka:9092
