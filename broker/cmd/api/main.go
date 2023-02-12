package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const PORT = "5001"

type Config struct {
	RabbitMQ *amqp.Connection
}

func main() {

	rabbitmqConnection, err := connectToRabbitMQ()

	if err != nil {
		panic("Not able to connect to RabbitMQ Server")
	}

	app := Config{
		RabbitMQ: rabbitmqConnection,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	log.Printf("Started Broker Server at port: %s", PORT)
	log.Fatal(server.ListenAndServe())
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
