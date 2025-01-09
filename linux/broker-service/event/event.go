package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)



func declareExcange(qch *amqp.Channel, exchangeName string) error {
	return qch.ExchangeDeclare(
		exchangeName, 
		"topic", 		// type
		true, 			// durable
		false, 			// auto-deleted
		false,			// internal
		false,			// no wait
		nil, 			// args
	)
}


// func declareRandomQueue(qch *amqp.Channel) (amqp.Queue, error) {
// 	return qch.QueueDeclare("",
// 		false,
// 		false,
// 		true, //exclusive
// 		false,
// 		nil,
// 	)
// }