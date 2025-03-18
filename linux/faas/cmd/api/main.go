package main

// Fintech as-a-Service

//Tools environment: GOPATH=C:\Users\tsizi\go, GOTOOLCHAIN=auto
//$env:GOPATH
//GOPROXY
//C:\Users\tsizi\go

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

type config struct {
	port int
}

type applicationContext struct {
	logger *log.Logger
	cfg    config
}

type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

const DEFAULT_PORT = 8083
const DEFAULT_RPC_PORT = 5003

func main() {
	s := time.Now()
	// s := time.Now(); defer func(){t := time.Since(s).Round(time.Millisecond);fmt.Printf("FAAS duration=%s", t)}()
	
	print(`Starting Module`)
	buildInfo, _ := debug.ReadBuildInfo()
	fmt.Printf(" '%+v': defaultPort=%d ...\n", buildInfo.Main.Path, DEFAULT_PORT)

	cfg := config{port: DEFAULT_PORT}
	// flag.IntVar(&cfg.port, "port", DEFAULT_PORT, "Yahoo Port")
	// flag.Parse()

	appCtx := applicationContext{
		logger: log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime),
		cfg:    cfg,
	}

	appCtx.logger.Printf(`Started Yahoo service on port %d`, DEFAULT_RPC_PORT)

	rpcServer := NewRPCServer(appCtx.logger); defer rpcServer.Shutdown()

	err := rpc.Register(rpcServer)
	if err != nil {
		appCtx.logger.Fatalf(`Error registering RPC server: %v`, err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", DEFAULT_RPC_PORT))
	if err != nil {
		appCtx.logger.Fatalf(`Error starting HTTP listener: %v`, err)
	} else {
		defer listener.Close()
	}

	go awaitShutdown(rpcServer, s)

	rpc.Accept(listener)
}

func awaitShutdown(r *RPCServer, s time.Time) {
	// Wait for Ctrl-C and closeUp if happened
	exitCh := make(chan os.Signal, 1) // terminate over Ctrl-C
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-exitCh
		r.Shutdown()
		t := time.Since(s).Round(time.Millisecond)
		fmt.Printf("FAAS duration=%s", t)
		os.Exit(0)
	}()
}
