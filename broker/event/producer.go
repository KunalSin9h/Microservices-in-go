package event

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	Conn *amqp.Connection
}

func NewProducer(conn *amqp.Connection) (Producer, error) {
	pro := Producer{
		Conn: conn,
	}

	err := pro.setup()

	if err != nil {
		return Producer{}, err
	}

	return pro, nil
}

func (p *Producer) setup() error {
	chnl, err := p.Conn.Channel() // Opens a channel
	if err != nil {
		return err
	}

	defer chnl.Close()

	return declareExchange(chnl)
}

func (p *Producer) Push(event, severity string) error {
	chnl, err := p.Conn.Channel() // Opens a channel

	if err != nil {
		return err
	}
	defer chnl.Close()
	//                                                   ExchangeName Routing Key                 Message
	err = chnl.PublishWithContext(context.Background(), "logs_topic", severity, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(event),
	})

	return err
}
