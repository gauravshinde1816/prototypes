package main

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Q1_name = "Q1"
var Q2_name = "Q2"

func sendMessage(ch *amqp.Channel, queueName string, message string) {
	ch.PublishWithContext(context.Background(), "", queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
	fmt.Println("Message sent to queue", queueName)
}

func main() {

	URL := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(URL)

	if err != nil {
		fmt.Println("Error connecting to RabbitMQ", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		panic(err)
	}

	sendMessage(ch, Q1_name, "1")
	sendMessage(ch, Q1_name, "2")
	sendMessage(ch, Q1_name, "CHECK")

	sendMessage(ch, Q2_name, "1")
	sendMessage(ch, Q2_name, "2")
	sendMessage(ch, Q2_name, "CHECK")

	select {}

}
