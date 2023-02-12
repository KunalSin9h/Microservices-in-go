package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(chnl *amqp.Channel) error {
	//                          ExName        ExType(Fanout, topic...etc)
	return chnl.ExchangeDeclare("logs_topic", "topic", true, false, false, false, nil)
}

func declareRandomQueue(chnl *amqp.Channel) (amqp.Queue, error) {
	return chnl.QueueDeclare("", true, false, true, false, nil)
}
