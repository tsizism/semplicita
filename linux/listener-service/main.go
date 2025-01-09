package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type config struct {
	maxConnCounts int
}

type applicationContext struct {
	logger *log.Logger
	cfg    config
}

func main() {
	// conn to rabbitmq

	appCfg := config{maxConnCounts: 2}

	appCtx := applicationContext{cfg: appCfg}
	appCtx.logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

	conn, err := appCtx.connectMQ()

	if err != nil {
		appCtx.logger.Fatalln(err)
		os.Exit(1)
	}

	defer closeConn(appCtx.logger, conn)

	// listen
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(conn, appCtx.logger)
	if err != nil {
		panic(err)
	}

	// watch q and cosnume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func closeConn(logger *log.Logger, conn *amqp.Connection) {
	logger.Println("Closing Connection ...")
	conn.Close()
	logger.Println("Connection Closed")
}

func (appCtx *applicationContext) connectMQ() (*amqp.Connection, error) {
	appCtx.logger.Println("Conneting to RabittMQ ..")
	var conn *amqp.Connection
	counts := 0

	for {
		// c, err := amqp.Dial("amqp://guest:guest@localhost")
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")

		if err != nil {
			appCtx.logger.Printf("RabbitMQ not yet ready, count=%d\n", counts)
			counts++

			if counts > appCtx.cfg.maxConnCounts { //
				fmt.Println(err)
				return nil, err
			}

			backoff := time.Duration(math.Pow(float64(counts), 2)) * time.Second
			appCtx.logger.Printf("Backing off for %d msec", backoff/time.Second)
			time.Sleep(backoff)

		} else {
			conn = c
			appCtx.logger.Println("Conneted to RabittMQ")
			break
		}
	}

	return conn, nil
}
