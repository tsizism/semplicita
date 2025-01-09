package main

import "net/http"

func (appCtx applicationContext) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/authenticate", appCtx.authenticateHandler)

	return mux
}


