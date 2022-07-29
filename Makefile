launch_kafka:
	docker pull spotify/kafka	
	docker rm -f kafka-docker-container
	docker run -d -p 2181:2181 -p 9092:9092 --name kafka-docker-container --env ADVERTISED_HOST=127.0.0.1 --env ADVERTISED_PORT=9092 spotify/kafka
	sleep 2
	docker exec kafka-docker-container /opt/kafka_2.11-0.10.1.0/bin/kafka-topics.sh --create --topic orders-topic --zookeeper localhost:2181 --partitions 1 --replication-factor 1
	docker exec -it kafka-docker-container /opt/kafka_2.11-0.10.1.0/bin/kafka-console-consumer.sh --topic orders-topic --from-beginning --bootstrap-server localhost:9092

run:
	KAFKA_ADDRESS=localhost:9092 TOPIC_NAME=orders-topic go run main.go
	
.PHONY: launch_kafka run