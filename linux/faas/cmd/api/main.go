package main

// Fintech as-a-Service

//Tools environment: GOPATH=C:\Users\tsizi\go, GOTOOLCHAIN=auto
//$env:GOPATH
//C:\Users\tsizi\go

import (
	"faas/fintechapi"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"runtime/debug"
)

type config struct {
	port int
}

type applicationContext struct{
	logger *log.Logger
	cfg	config
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
	print(`Starting Module`)
    buildInfo, _ := debug.ReadBuildInfo()
    fmt.Printf(" '%+v': defaultPort=%d ...\n", buildInfo.Main.Path, DEFAULT_PORT)

	cfg := config{}
	flag.IntVar(&cfg.port, "port", DEFAULT_PORT, "Yahoo Port")
	flag.Parse()

	appCtx := applicationContext{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		cfg: cfg,
	}

	appCtx.logger.Printf(`Started Yahoo service on port %d`, DEFAULT_RPC_PORT)

	rpcServer := fintechapi.NewRPCServer(fintechapi.StockAPI{})

	err := rpc.Register(rpcServer); if err != nil {
		appCtx.logger.Fatalf(`Error registering RPC server: %v`, err)
	}

	// rpc.HandleHTTP()

    listener, err :=  net.Listen("tcp", fmt.Sprintf(":%d", DEFAULT_RPC_PORT)); if err != nil {
		appCtx.logger.Fatalf(`Error starting HTTP listener: %v`, err)
	}	

	defer listener.Close()

	rpc.Accept(listener)

}