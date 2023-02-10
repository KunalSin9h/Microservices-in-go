package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	Conn      *amqp.Connection
	QueueName string
}

// The actual message been shared in queue
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		Conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (c *Consumer) setup() error {
	chnl, err := c.Conn.Channel()

	if err != nil {
		return nil
	}

	defer chnl.Close()

	return declareExchange(chnl)
}

// Listening to queues
func (c *Consumer) Listen(topics []string) error {
	chnl, err := c.Conn.Channel()
	if err != nil {
		return nil
	}
	defer chnl.Close()

	queue, err := declareRandomQueue(chnl)

	if err != nil {
		return nil
	}

	for _, topic := range topics {
		err := chnl.QueueBind(queue.Name, topic, "logs_topic", false, nil)
		if err != nil {
			return err
		}
	}

	message, err := chnl.Consume(queue.Name, "", true, false, false, false, nil) // message is (<-chan amqp.Delivery)

	if err != nil {
		return err
	}

	// HoldExecution
	holdExecution := make(chan bool)

	go func() {
		for d := range message {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			go handleRequest(payload) // This makes handling each request is separate go routine
			// which prevents code blocking for single request
		}
	}()

	log.Printf("Waiting for message on [Exchange, Queue] [logs_topic, %s]\n", queue.Name)

	<-holdExecution
	return nil
}

func handleRequest(load Payload) {
	switch load.Name {
	case "log", "event":
		err := logEvent(load)
		if err != nil {
			log.Println(err)
		}
	default:
		err := logEvent(load)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(load Payload) error {
	jsonByte, _ := json.MarshalIndent(load, "", "\t")
	res, err := http.Post("http://logger:5003/log", "application/json", bytes.NewBuffer(jsonByte))

	if err != nil {
		return err
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusBadRequest:
		return errors.New("bad request")
	case http.StatusInternalServerError:
		return errors.New("internal server error")
	default:
		return nil
	}
}
