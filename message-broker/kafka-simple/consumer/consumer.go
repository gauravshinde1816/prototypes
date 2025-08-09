package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func readMessage(r *kafka.Reader, partition int) {

	fmt.Printf("Consumer %v started \n", partition)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("failed to read message:", err)
		}

		fmt.Printf("message from p%d off=%d key=%s val=%s\n", m.Partition, m.Offset, m.Key, m.Value)

	}

}

func main() {

	r1 := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"127.0.0.1:9092"},
		Topic:       "my-topic",
		Partition:   0,
		StartOffset: kafka.LastOffset,
		MinBytes:    10e3,
		MaxBytes:    10e6,
	})

	r2 := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"127.0.0.1:9092"},
		Topic:       "my-topic",
		Partition:   1,
		StartOffset: kafka.LastOffset,
		MinBytes:    10e3,
		MaxBytes:    10e6,
	})

	r3 := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"127.0.0.1:9092"},
		Topic:       "my-topic",
		Partition:   2,
		StartOffset: kafka.LastOffset,
		MinBytes:    10e3,
		MaxBytes:    10e6,
	})

	defer r1.Close()
	defer r2.Close()
	defer r3.Close()

	go readMessage(r1, 0)
	go readMessage(r2, 1)
	go readMessage(r3, 2)

	select {}

}
