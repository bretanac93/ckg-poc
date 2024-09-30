#!/usr/bin/env bash

echo "Creating topics in kafka..."
kafka-topics --create --if-not-exists --bootstrap-server kafka:9092  --topic local.orders --partitions 3 --replication-factor 1
kafka-topics --create --if-not-exists --bootstrap-server kafka:9092  --topic local.payments --partitions 3 --replication-factor 1

echo "Topics created in kafka:"
kafka-topics --bootstrap-server kafka:9092 --list

echo "Done!"
