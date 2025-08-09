package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {

	w := kafka.Writer{
		Addr:     kafka.TCP("127.0.0.1:9092"),
		Topic:    "my-topic",
		Balancer: &kafka.LeastBytes{},
		// MaxAttempts: 1, // NO implicit retries → no dupes
	}

	messages := []kafka.Message{
		{
			Key:       []byte("Key 1"),
			Value:     []byte("Message 1 on Partition 0"),
			Partition: 0,
		},
		{
			Key:       []byte("Key 2"),
			Value:     []byte("Message 2 on Partition 0"),
			Partition: 0,
		},

		{
			Key:       []byte("Key 1"),
			Value:     []byte("Message 1 on Partition 1"),
			Partition: 1,
		},
		{
			Key:       []byte("Key 2"),
			Value:     []byte("Message 2 on Partition 1"),
			Partition: 1,
		},

		{
			Key:       []byte("Key 1"),
			Value:     []byte("Message 1 on Partition 2"),
			Partition: 2,
		},
		{
			Key:       []byte("Key 2"),
			Value:     []byte("Message 2 on Partition 2"),
			Partition: 2,
		},
	}

	time.Sleep(10 * time.Second)
	fmt.Println("Producer started writing")
	err := w.WriteMessages(context.Background(), messages...)

	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	fmt.Println("Producer finished")

}
