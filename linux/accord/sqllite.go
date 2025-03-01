package main

// go get github.com/mattn/go-sqllite3

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	// _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var once sync.Once


func main2() {
	fmt.Printf("sqlite3")

	run()
	run()
	run()
	run()	
}
func run() {	
	once.Do( func() {
		log.Println("opening db conn ..")

		db , err := sql.Open("sqlite3", "./sqlite3.db"); if err != nil {
			log.Fatal( err)		
		}
		defer db.Close()

		log.Println("db conn opned")
	})

	println("bla")
}


