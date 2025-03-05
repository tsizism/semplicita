package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"syscall"

	"os/signal"
	"time"
	"trace/data"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	port     int
	rpcPort  int
	grpcPort int
}

var mongoClient *mongo.Client

type applicationContext struct {
	cfg       config
	logger    *log.Logger
	httpSrv   *http.Server
	mongo     *mongo.Client
	models    data.Models
	sessionId string
}

func main() {
	var appCfg config

	// go run .\cmd\api\. -port 8888
	flag.IntVar(&appCfg.port, "port", 80, "REST port")
	flag.Parse()

	appCfg.rpcPort = 5001
	appCfg.grpcPort = 50001

	appLogger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	appCtx := &applicationContext{
		cfg:       appCfg,
		logger:    appLogger,
		sessionId: uuid.NewString(),
	}

	// Connect to mongo and set disconnect handler
	err := appCtx.connectToMongo()

	if err != nil {
		appCtx.logger.Panic(err)
		return
	}

	// create ctx for disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // LIFO - 2

	defer func() { // LIFO - 1
		appCtx.cleanup("main defer", ctx)
		// appCtx.logger.Println("Main: Disconnect Monogo")
		// traceEvent.Data = "Exiting..."
		// err = appCtx.models.Insert(traceEvent)
		// if err = appCtx.mongo.Disconnect(ctx); err != nil {
		// 	appCtx.logger.Panic(err)
		// }
	}() // () execute anonymus func

	// End of Connect to mongo

	appCtx.models = data.New(appCtx.mongo, appCtx.logger)

	err = appCtx.InsertStartStopServiceEvent("Start")

	if err != nil {
		return
	}

	err = rpc.Register(new(RPCServer))

	if err != nil {
		appCtx.logger.Println("Failed to register RPC Server=", err)
		return
	}

	go appCtx.rpcListen()

	go appCtx.gRPCListen()

	// Wait for Ctrl-C and closeUp if happened
	// ch := make(chan os.Signal, 1) // terminate over Ctrl-C
	// signal.Notify(ch, os.Interrupt) Ctrl-C
	// signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	// go appCtx.closeUp(ch)

	go shutdown()
	forever := make(chan int)

	appCtx.serve()

	<-forever
}

func shutdown() {
	// Wait for Ctrl-C and closeUp if happened
	exitCh := make(chan os.Signal, 1) // terminate over Ctrl-C
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-exitCh
		fmt.Println("Shutdown")
		// Cleanup
		os.Exit(0)
	}()
}


func (appCtx *applicationContext) InsertStartStopServiceEvent(startStop string) error {
	traceEvent := data.TraceEntry{
		Src:  "trace",
		Via:  "",
		Data: fmt.Sprintf("%s trace-service sessionId=%s", startStop, appCtx.sessionId),
	}
	err := appCtx.models.Insert(traceEvent)

	if err != nil {
		appCtx.logger.Println("Failed to insert Start/Stop event", err)
		return err
	}

	return nil
}

// func (appCtx *applicationContext) closeUp(ch <-chan os.Signal) {
// 	<-ch // wait  signal
// 	appCtx.cleanup("CloseUp UPON signal", context.TODO())
// 	// os.Exit(1)	cleanup will unblock main()
// }

func (appCtx *applicationContext) cleanup(exitReason string, ctx context.Context) {
	appCtx.logger.Println("exitReason", exitReason)

	if appCtx.mongo != nil {
		err := appCtx.InsertStartStopServiceEvent("Stop")
		// traceEvent := data.TraceEntry {
		// 	Src: "trace",
		// 	Via: "",
		// 	Data: fmt.Sprint("Exiting:", appCtx.sessionId),
		// }
		// err := appCtx.models.Insert(traceEvent)

		if err != nil {
			appCtx.logger.Panic(err)
		}

		appCtx.logger.Println("Disconnect Mango")
		if err := appCtx.mongo.Disconnect(ctx); err != nil {
			appCtx.logger.Panic(err)
		}
		appCtx.mongo = nil
	}

	// unblock appCtx.httpSrv.ListenAndServe()
	if appCtx.httpSrv != nil {
		appCtx.logger.Println("Closing HTTP server")
		appCtx.httpSrv.Close()
		appCtx.httpSrv = nil
	}
}

func (appCtx *applicationContext) serve() {
	appCtx.httpSrv = &http.Server{
		Addr:    fmt.Sprintf(":%d", appCtx.cfg.port),
		Handler: appCtx.routes(),
	}

	appCtx.logger.Println("Starting Trace service on port", appCtx.cfg.port)

	err := appCtx.httpSrv.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		appCtx.logger.Panic(err)
	}

}

func (appCtx *applicationContext) rpcListen() error {
	appCtx.logger.Println("Start rpcListen")

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", appCtx.cfg.rpcPort))

	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			appCtx.logger.Println("Failed to Accept", err)
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}

func (appCtx *applicationContext) connectToMongo() error {
	// mongoURL := "mongodb://localhost:27017"  // client
	mongoURL := "mongodb://mongodb:27017"
	clientOptions := options.Client().ApplyURI(mongoURL)
	auth := options.Credential{Username: "admin", Password: "password"}
	clientOptions.SetAuth(auth)

	var err error

	mongoClient, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		appCtx.logger.Println("Faild to connect", err)
		return err
	}
	appCtx.mongo = mongoClient
	appCtx.logger.Printf("Connected to mongo using %+v", clientOptions.Auth)

	return nil
}
