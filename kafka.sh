#!/bin/sh

docker exec kafka_task /opt/bitnami/kafka/bin/kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create --topic tasks.created --partitions 3 --replication-factor 1

docker exec kafka_task /opt/bitnami/kafka/bin/kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create --topic tasks.retry.1m --partitions 3 --replication-factor 1

docker exec kafka_task /opt/bitnami/kafka/bin/kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create --topic tasks.retry.5m --partitions 3 --replication-factor 1

docker exec kafka_task /opt/bitnami/kafka/bin/kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create --topic tasks.retry.30m --partitions 3 --replication-factor 1

docker exec kafka_task /opt/bitnami/kafka/bin/kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create --topic tasks.dlq --partitions 3 --replication-factor 1