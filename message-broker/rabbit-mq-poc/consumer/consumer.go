package main

import (
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var URL = "amqp://guest:guest@localhost:5672/"
var Q1_name = "Q1"
var Q2_name = "Q2"

func main() {
	var conn, err = amqp.Dial(URL)

	if err != nil {
		panic(err)
	}

	ch1, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	go func() {
		msgs, err := ch1.Consume("Q1", "", true, false, false, false, nil)

		fmt.Println("Consumer 1 for Q1 started")
		finalString := ""
		if err != nil {
			panic(err)
		}

		for msg := range msgs {
			fmt.Println("Message received for Consumer 1", string(msg.Body))
			if string(msg.Body) == "CHECK" {
				if finalString == "12" {
					fmt.Println("Messages are in order. Final string is", finalString, "PASS")

				} else {
					fmt.Println("Messages are not in order. Final string is", finalString, "FAIL")
					os.Exit(1)
				}
			} else {
				finalString += string(msg.Body)
			}

		}

	}()

	go func() {
		msgs2, err := ch1.Consume("Q2", "", true, false, false, false, nil)
		fmt.Println("Consumer 1 for Q2 started")
		finalString := ""
		if err != nil {
			panic(err)
		}

		for msg := range msgs2 {
			fmt.Println("Message received for Consumer 2", string(msg.Body))
			if string(msg.Body) == "CHECK" {
				if finalString == "12" {
					fmt.Println("Messages are in order. Final string is", finalString, "PASS")

				} else {
					fmt.Println("Messages are not in order. Final string is", finalString, "FAIL")
					os.Exit(1)
				}
			} else {
				finalString += string(msg.Body)
			}

		}
	}()

	go func() {
		msgs2, err := ch1.Consume("Q2", "", true, false, false, false, nil)

		fmt.Println("Consumer 2 for Q2 started")
		finalString := ""
		if err != nil {
			panic(err)
		}

		for msg := range msgs2 {
			fmt.Println("Message received for Consumer 3", string(msg.Body))
			if string(msg.Body) == "CHECK" {
				if finalString == "12" {
					fmt.Println("Messages are in order. Final string is", finalString, "PASS")

				} else {
					fmt.Println("Messages are not in order. Final string is", finalString, "FAIL")
					os.Exit(1)
				}
			} else {
				finalString += string(msg.Body)
			}

		}
	}()

	select {}

}
