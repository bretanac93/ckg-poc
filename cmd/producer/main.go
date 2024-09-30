package main

import (
	"github.com/bretanac93/ckg/internal/order"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	// define topics to listen from. Also wildcards are supported, something like "local.*"
	// that's unsafe tho. I prefer to declare all topics explicitly.
	topic := "local.orders"

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":     "127.0.0.1:29092",
		"broker.address.family": "v4",
	})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	for i := 0; i < 1000000; i++ {
		k, v := generateMessage()
		if err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            k,
			Value:          v,
		}, nil); err != nil {
			panic(err)
		}
	}

	producer.Flush(15 * 1000)
}

func generateMessage() ([]byte, []byte) {
	o := order.Random()

	payload, err := o.ToJSON()
	if err != nil {
		panic(err)
	}

	return []byte(o.ID.String()), payload
}
