package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type config struct {
    port int
	gRPCPort int
	maxConnCounts int
}

type applicationContext struct {
    cfg config
    logger *log.Logger
	connMQ *amqp.Connection
}

// 
// go build -o <your desired name>
//  curl http://localhost:4000
//d  run -it -p 8080:80 brokerapp
//dc up -d  (-d detached)
func main() {
    appCfg := config{
    	maxConnCounts: 2,
		gRPCPort: 50001,
    }

	defaultPort := 8080
    flag.IntVar(&appCfg.port, "port", defaultPort, "API server port")
    flag.Parse()

    appLogger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
    appLogger.Printf("Starting broker service on port %d\n",appCfg.port)

    appCtx := &applicationContext {
        cfg : appCfg,
        logger : appLogger,
    }

    addr := fmt.Sprintf(":%d", appCtx.cfg.port)

	conn, err := appCtx.connectMQ()

	if err != nil {
		appCtx.logger.Fatalln(err)
		os.Exit(1)
	}

	defer conn.Close()

	appCtx.connMQ = conn

    srv := &http.Server{
        Addr: addr,
        Handler: appCtx.routes(),
    }

    err = srv.ListenAndServe()

    if err != nil {
        appLogger.Panic(err)
    }
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

			if counts > appCtx.cfg.maxConnCounts  { //
				fmt.Println(err)
				return nil, err
			}

			backoff := time.Duration(math.Pow(float64(counts), 2)) * time.Second
			appCtx.logger.Printf("Backing off for %d sec", backoff/time.Second)
			time.Sleep(backoff)

		} else {
			conn = c
			appCtx.logger.Println("Conneted to RabittMQ")
			break	
		}
	}

	return conn, nil
}