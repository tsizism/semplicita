package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)

var rawOrders = []string{
	`{"productCode":1111, "quantity": 5, 	"status": 0}`,
	`{"productCode":2222, "quantity": 42.3, "status": 1}`,
	`{"productCode":3333, "quantity": 19, 	"status": 2}`,
	`{"productCode":4444, "quantity": 8, 	"status": 3}`,
}

// var orders = []order{}

type blablabla struct {
	bla int
	blabla    float64
}

func main() {
	// var newOrder *blablabla
	// newOrder = blablabla{}
	// fmt.Printf("New Order : %v\n", newOrder)

	orderMgmt()
}

func bufChannel() {
	ch := make(chan string, 1)
	ch <- "message "
	res := <-ch
	fmt.Println(res)
}

func orderMgmt() {
	var wg sync.WaitGroup
	// var receivedOrdercCh = make(chan order)
	// var validateOrdercCh = make(chan order)
	// var invalidOrdercCh = make(chan invalidOrder)

	wg.Add(1)
	// go receiveOrders(receivedOrdercCh)
	receivedOrdercCh := receiveOrders() // ro <-chan order
	validateOrdercCh, invalidOrdercCh := validateOrders(receivedOrdercCh) // ro chan oreder, ro chan invalidOrder

	// go validateOrders(receivedOrdercCh, validateOrdercCh, invalidOrdercCh)

	go func(validOrderCh <-chan order, invalidOrderCh <-chan invalidOrder) { // inCh rcv from channel
		loop : for{
		select {
		case order, ok := <-validOrderCh:
			if ok {
				fmt.Printf("Valid Order Received: %v\n", order)
			} else {
				break loop
			}

		case order, ok := <-invalidOrderCh:
			if ok { 
				fmt.Printf("Invalid Order Received: %v, err=%v\n", order.Order, order.Err)
			} else {
				break loop
			}
		// default:
		// 	println(`default`)
		
		}
	}
	wg.Done()

	}(validateOrdercCh, invalidOrdercCh)

	// go func(inCh <-chan order) { // rcv from ch
	// 	// order := <- validateOrdercCh
	// 	order := <- inCh
	// 	fmt.Printf("Valid Order Received: %v\n", order)
	// 	wg.Done()
	// }(validateOrdercCh)

	// go func(inCh <-chan invalidOrder) {
	// 	invalidOrder := <- inCh // rcv from ch
	// 	// invalidOrder := <- invalidOrdercCh
	// 	fmt.Printf("Invalid Order Received: %v, err=%v\n", invalidOrder.Order, invalidOrder.Err)
	// 	wg.Done()
	// }(invalidOrdercCh)

	wg.Wait()
	fmt.Println("Done!")

}


func validateOrders(inCh <-chan order)  (<-chan order, <-chan invalidOrder) {
	// order := <- inCh // inCh rcv from ch

	outCh := make(chan order)
	errCh := make(chan invalidOrder, 1)

	go func () {
		for order := range inCh {
			if order.Quantity <= 0 {
				errCh <- invalidOrder{Order: order, Err: errors.New("invalid negavtive order quantity")}
				// errCh <- invalidOrder { Order: order, Err: fmt.Errorf("invalid order quantity %f", order.Quantity) }
			} else {
				// fmt.Printf("Valid order  %v\n", order)
				outCh <- order // outCh snd to ch
			}
		}
		close(outCh)
		close(errCh)
	}()

	return outCh, errCh // as ro channel

}



// // <-chan order ro
// func validateOrders(inCh <-chan order, outCh chan<- order, errCh chan<- invalidOrder) {
// 	// order := <- inCh // inCh rcv from ch

// 	for order := range inCh {
// 		if order.Quantity <= 0 {
// 			errCh <- invalidOrder{Order: order, Err: errors.New("invalid negavtive order quantity")}
// 			// errCh <- invalidOrder { Order: order, Err: fmt.Errorf("invalid order quantity %f", order.Quantity) }
// 		} else {
// 			// fmt.Printf("Valid order  %v\n", order)
// 			outCh <- order // outCh snd to ch
// 		}
// 	}
// 	close(outCh)
// 	close(errCh)
// }


func receiveOrders() <-chan order { 
	outCh := make(chan order) // rw chan

	go func() {
		for _, rawOrder := range rawOrders {
			var newOrder order

			err := json.Unmarshal([]byte(rawOrder), &newOrder)

			if err != nil {
				log.Println(err)
				continue
			}
			// fmt.Printf("Receiving new order %v\n", newOrder)
			//orders = append(orders, newOrder)
			outCh <- newOrder
		}
		close(outCh)
	}()

	return outCh // returm as ro
}


// func receiveOrders(outCh chan<- order) {
// 	for _, rawOrder := range rawOrders {
// 		var newOrder order

// 		err := json.Unmarshal([]byte(rawOrder), &newOrder)

// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 		// fmt.Printf("Receiving new order %v\n", newOrder)
// 		//orders = append(orders, newOrder)
// 		outCh <- newOrder
// 	}
// 	close(outCh)
// }
