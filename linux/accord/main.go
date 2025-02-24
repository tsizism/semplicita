package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

var rawOrders = []string{
	`{"productCode":1111, "quantity": 5, "status": 0}`,
	`{"productCode":2222, "quantity": 42.3, "status": 1}`,
	`{"productCode":3333, "quantity": 19, "status": 2}`,
	`{"productCode":4444, "quantity": 8, "status": 3}`,
}

var orders = []order{}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go recieveOrders(&wg)
	wg.Wait()
	fmt.Println(orders)
}

func recieveOrders(wg *sync.WaitGroup) {
	for _, rawOrder := range rawOrders {
		var newOrder order

		err := json.Unmarshal([]byte(rawOrder), &newOrder)

		if err != nil {
			log.Println(err)
			continue
		}
		//fmt.Println(newOrder)
		orders = append(orders, newOrder)
	}
	wg.Done()
	wg.Done()
}
