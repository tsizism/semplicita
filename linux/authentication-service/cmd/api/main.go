package main

import (
	"authentication/data"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type config struct {
	port int
	// DB *sql.DB
	// Models data.Models
	repo data.Repository
	Client* http.Client
}

type applicationContext struct {
	cfg    config
	logger *log.Logger
}

func main() {
	defaultPort := 8081
	appCfg :=  config {
		Client: http.DefaultClient,
	}
	flag.IntVar(&appCfg.port, "port", defaultPort, "Authentication service port")
	flag.Parse()

	appLogger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	appCtx := &applicationContext {
		cfg: appCfg,
		logger: appLogger,
	}

	conn := appCtx.connectToDB()
	appCfg.repo = data.NewPostgressRepository(conn)

	// appCtx.cfg.DB = conn
	// appCtx.cfg.Models = data.New(conn)

	if conn == nil {
		appCtx.logger.Panic("Failed to connect to Postgress")
	}

	appCtx.logger.Printf("Starting Authentication service port=%d", appCtx.cfg.port)

	srv := &http.Server {
		Addr: fmt.Sprintf(":%d", appCtx.cfg.port),
		Handler: appCtx.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		appCtx.logger.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {

	db, err :=sql.Open("pgx", dsn)
	
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (appCtx *applicationContext) connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")	

	for i:=0; i < 10; i++ {
		conn, err := openDB(dsn)

		if err == nil {
			appCtx.logger.Println("Postgress connected")
			return conn
		}

		appCtx.logger.Println("Postgress not ready, wait 2 sec and try again")
		time.Sleep(2 * time.Second)
	}

	appCtx.logger.Println("Gave up connecting DB")

	return nil
}

