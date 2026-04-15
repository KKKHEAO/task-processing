#sh

docker exec -it kafka_task bash

kafka-topics.sh \
--bootstrap-server localhost:9092 \
--create \
--topic tasks.created \
--partitions 3 \
--replication-factor 1

kafka-topics.sh \
--bootstrap-server localhost:9092 \
--create \
--topic tasks.retry.1m \
--partitions 3 \
--replication-factor 1

kafka-topics.sh \
--bootstrap-server localhost:9092 \
--create \
--topic tasks.retry.5m \
--partitions 3 \
--replication-factor 1

kafka-topics.sh \
--bootstrap-server localhost:9092 \
--create \
--topic tasks.retry.30m \
--partitions 3 \
--replication-factor 1

kafka-topics.sh \
--bootstrap-server localhost:9092 \
--create \
--topic tasks.dlq \
--partitions 3 \
--replication-factor 1