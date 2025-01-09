package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type config struct {
	port int
}

type applicationContext struct {
	appCfg config
	logger *log.Logger
	mailer mailProp
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 80, "Web port")
	flag.Parse()

	appCtx := applicationContext{
		appCfg: cfg,
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		mailer: newMailer(),
	}

	appCtx.logger.Printf("Mailer env:%+v", appCtx.mailer)

	appCtx.logger.Println("Starting mail service on port", appCtx.appCfg.port)


	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", appCtx.appCfg.port),
		Handler: appCtx.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		appCtx.logger.Panic(err)
	}
}

func newMailer() mailProp {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	m := mailProp{
		domain:     os.Getenv("MAIL_DOMAIN"),
		host:       os.Getenv("MAIL_HOST"),
		port:       port,
		username:   os.Getenv("MAIL_USERNAME"),
		password:   os.Getenv("MAIL_PASSWORD"),
		encryption: os.Getenv("MAIL_ENCRYPTION"),
		fromName:   os.Getenv("MAIL_FROMNAME"),
		fromAddr:   os.Getenv("MAIL_FROMADDR"),
	}

	return m
}
