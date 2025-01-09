package event

import (
	"log"
	

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	exchangeName string
	conn *amqp.Connection
	logger *log.Logger	
}

func NewEventEmitter(conn *amqp.Connection, logger *log.Logger) (Emitter, error){
	emitter := Emitter {
		exchangeName : "logs_exchange",
		conn: conn,
		logger: logger,
	}

	err := emitter.setup()

	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}


func (e *Emitter) setup() error {
	e.logger.Println("Setting up Emitter for exchange", e.exchangeName)

	qch, err := e.conn.Channel()

	if err != nil {
		return err
	}

	defer qch.Close()

	return declareExcange(qch, e.exchangeName)
}


func (e *Emitter) Push(event string, severity string) error {
	qch, err := e.conn.Channel()

	if err != nil {
		return err
	}

	defer qch.Close()
	
	e.logger.Printf("Pushing to MQ channel exchange=%s, event=%s, sev=%s", e.exchangeName, event, severity)

	err = qch.Publish(e.exchangeName,
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(event),
		},
	)
	
	if err != nil {
		return  err
	}
	
	return nil
}

