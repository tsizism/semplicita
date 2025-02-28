package main

import (
	"fmt"
	"net/http"
	"time"
)

func main1() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
		fmt.Fprintf(w, "Hello, this is a dummy service!") // func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
	})

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}