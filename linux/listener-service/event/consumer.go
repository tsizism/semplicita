package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"

)

type qconsumer struct {
	exchangeName string
	conn *amqp.Connection
	logger *log.Logger
	queueName string
	// appCtx *main.ApplicationContext
}

func NewConsumer(conn *amqp.Connection, logger *log.Logger)  (qconsumer, error) {
	consumer := qconsumer {
		exchangeName : "logs_exchange",
		conn: conn,
		queueName : "Queen",
		logger: logger,
	}

	err := consumer.setup()

	if err != nil {
		return qconsumer{}, err
	}

	return consumer, nil
}

func (consumer *qconsumer) setup() error {
	consumer.logger.Println("Setting up consumer for exchange", consumer.exchangeName)
	qch, err := consumer.conn.Channel()

	if err != nil {
		return err
	}

	return declareExcange(qch, consumer.exchangeName)
}

type JSONPayload struct {
	Src string `json:"src"`
	Via string `json:"via"`
	Data string `json:"data"`
}

func (consumer *qconsumer) Listen(topics []string) error {
	qch, err := consumer.conn.Channel()
	
	if err != nil {
		return err
	}

	defer qch.Close()

	q, err := declareRandomQueue(qch)

	if err != nil {
		return err
	}

	for _, topic := range topics {
		err := qch.QueueBind(q.Name, topic, consumer.exchangeName, false, nil)

		if err != nil {
			return err
		}
	}

	deliveredMessages, err := qch.Consume(q.Name, "", true, false, false, false, nil)

	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for deliveredMessage := range deliveredMessages {
			// consumer.logger.Printf("deliveredMessage=%+v", deliveredMessage)
			consumer.logger.Printf("deliveredMessage Exchange=%s RoutingKey=%s", deliveredMessage.Exchange, deliveredMessage.RoutingKey)

			var payload JSONPayload
			json.Unmarshal(deliveredMessage.Body, &payload)

			go consumer.handlePayload(payload)
		}
	}()

	consumer.logger.Printf("Waiting for message [Exchange, Queue] [%s, %s]\n", consumer.exchangeName, q.Name)
	<-forever

	return nil
}

func (consumer *qconsumer)handlePayload(payload JSONPayload) {
	consumer.logger.Println("consumer.handlePayload", payload)

	switch payload.Src {
		case "log", "event" : 
			err := consumer.logEvent(payload)

			if err != nil {
				consumer.logger.Printf("handlePayload: Failed to logEvent %s", err)
				return
			}

		case "auth":
			//authenticate

		default:
			consumer.logger.Println("consumer.handlePayload default case")
			err := consumer.logEvent(payload)

			if err != nil {
				consumer.logger.Printf("handlePayload(default): Failed to logEvent %s", err)
				return
			}
	}

}

func (consumer *qconsumer)logEvent(payload JSONPayload) error {
	consumer.logger.Printf("logEvent: payload=%+v", payload)

	traceServiceEndpoint := "http://trace-service/trace"

	jsonData, _ := json.MarshalIndent(payload, "", "\t") // Marshal in prod

	request, err := http.NewRequest("POST", traceServiceEndpoint, bytes.NewBuffer(jsonData))

	if err != nil {
		consumer.logger.Printf("traceEvent: Failed to create endpoint request %s", err)
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		consumer.logger.Printf("traceEvent: Failed to call endpoint %s", err)
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		consumer.logger.Print("traceEvent:failed calling trace service")
		return errors.New("response.StatusCode != http.StatusAccepted")
	}

	return nil
}