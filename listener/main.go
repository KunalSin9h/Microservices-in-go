package main

import (
	"log"
	"math"
	"time"

	"listener/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try  connect to rabbitmq server
	rabbitConn, err := connectToRabbitMQ() // rabbitConn -> *amqp.Connection

	if err != nil {
		log.Fatal(err)
	}

	defer rabbitConn.Close()

	consumer, err := event.NewConsumer(rabbitConn) // consumer -> event.Consumer

	if err != nil {
		log.Fatal(err)
	}

	topics := []string{"log.INFO", "log.WARNING", "log.ERROR"}
	log.Fatal(consumer.Listen(topics))
}

func connectToRabbitMQ() (conn *amqp.Connection, err error) {
	var count int32
	var countLimit int32 = 5
	var sleepTime = 1 * time.Second

	for {
		//                    "-------Connection-String----------"
		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672")
		if err != nil {
			count++
			log.Printf("Trying to connect with RabbitMQ...(%d attempts more remaining)\n", countLimit-count)
		} else {
			break
		}
		if count >= countLimit {
			log.Println("Could't connect to RabbitMQ!")
			break
		}
		// Delay will be like 1sec -> 4sec -> 9sec -> 16sec -> ...
		sleepTime = time.Duration(math.Pow(float64(count), 2)) * time.Second
		log.Printf("Retrying to connect after %v seconds\n", sleepTime)
		time.Sleep(sleepTime)
		continue
	}

	return
}
