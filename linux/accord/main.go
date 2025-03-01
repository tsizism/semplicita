package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)

var rawOrders = []string{
	`{"productCode":1111, "quantity": -5, 	"status": 0}`,
	`{"productCode":2222, "quantity": 42.3, "status": 0}`,
	`{"productCode":3333, "quantity": 19, 	"status": 0}`,
	`{"productCode":4444, "quantity": 8, 	"status": 0}`,
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

	var bla [2]int
	bla[0] = 1
	log.Printf("%v", bla)

	baz("2")
	orderMgmt()
}

func baz(x interface{}) {
	switch x.(type) {
	case int:
		println("int")
	
	case string:
		println("str")
	}
}

func bufChannel() {
	ch := make(chan string, 1)
	ch <- "message "
	res := <-ch
	fmt.Println(res)
}

func orderMgmt() {
	var orderMgmtWaitGroup sync.WaitGroup
	orderMgmtWaitGroup.Add(1) // drain invalidOrderCh

	// workflow
	// producer
	receivedOrdercCh := receiveOrders() // ro <-chan order

	// consumer-producer
	validateOrdercCh, invalidOrdercCh := validateOrders(receivedOrdercCh) // ro chan oreder, ro chan invalidOrder	
	
	// multi-consumer-producer
	reserveInventoryCh := reserveInventory(validateOrdercCh) // ro <-chan order
	
	// consumer
	fillOrders(reserveInventoryCh, &orderMgmtWaitGroup)
	// fillOrdersCh := fillOrders(reserveInventoryCh)
	

	// drain invalidOrderCh
	go func(invalidOrderCh <-chan invalidOrder) {
		for order := range invalidOrderCh{
		 	fmt.Printf("Invalid Order Received: %v, err=%v\n", order.Order, order.Err)
		}
		orderMgmtWaitGroup.Done()
	}(invalidOrdercCh) 


	// go func() { // eavesdropping channels 
	// 	defer fmt.Printf("eavesdropping dead\n")
	// 	select {
	// 		case o := <- invalidOrdercCh:	
	// 			fmt.Printf("eavesdropping invalidOrdercCh o=%v\n", o)
	// 		// case o := <- receivedOrdercCh:
	// 		// 	fmt.Printf("eavesdropping receivedOrdercCh o=%v\n", o)
	// 		// case o := <- validateOrdercCh:	
	// 		// 	fmt.Printf("eavesdropping validateOrdercCh o=%v\n", o)
	// 		// case o := <- reserveInventoryCh:	
	// 		// 	fmt.Printf("eavesdropping reserveInventoryCh o=%v\n", o)
	// 	}
	// 	orderMgmtWaitGroup.Done()
	// }()


	// drain the last workflow fillOrdersCh, pass ro channel
	// const workers = 3
	// wg.Add(workers)
	// for i := 0; i < workers; i++ {
	// 	go func(fillOrdersCh <-chan order) {
	// 		fmt.Printf("Start drain worker %d\n", i)
	// 		for order := range fillOrdersCh {
	// 			fmt.Printf("Filled order %v dw=%d\n", order, i)
	// 		}
	// 		defer fmt.Printf("Stop drain worker %d\n", i)
	// 		wg.Done()
	// 	}(fillOrdersCh)
	// }

	orderMgmtWaitGroup.Wait()
	fmt.Println("Done!")

}

// mutiple consumers
func fillOrders(fillOrdersInCh <-chan order, wg* sync.WaitGroup) { //<-chan order{
	// fillOrdersOutCh := make(chan order)

	const workers = 3
	wg.Add(workers)
	for i:=0; i < workers; i++ {
		go func() { // spawn anonymus goroutine to change order status
			for order := range fillOrdersInCh {
				order.Status = filled
				fmt.Printf("Filled order %v workerId=%d\n", order, i)
				// fillOrdersOutCh <- order
			}
			// No more downstream workers dependant on this channel
			//close(fillOrdersOutCh)
			// fmt.Printf("------------------->fillOrders done\n")
			wg.Done()
		}()
	}

	// return fillOrdersOutCh 
}

// multiple producers - sends to 3 workers
func reserveInventory(reserveInventoryInCh <-chan order) <-chan order {
	reserveInventoryOutCh := make(chan order)

	const workers = 3
	var wg sync.WaitGroup
	wg.Add(workers)
	
	for i:= 0; i < workers; i++ {
		go func() { // spawn anonymus goroutine to change order status
			defer wg.Done()
			for o := range reserveInventoryInCh { // consume messages
				o.Status = reserved
				reserveInventoryOutCh <- o		 // produce 1 message
			}
			
			// close(reserveInventoryOutCh)
		}()
	}
	
	go func() { // spawn anonymus goroutine to close out ch
		wg.Wait()
		close(reserveInventoryOutCh)
	}()

	return reserveInventoryOutCh
}

// mutiple producer
func validateOrders(validateOrdersInCh <-chan order)  (<-chan order, <-chan invalidOrder) {
	validateOrdersOutCh := make(chan order)
	errCh := make(chan invalidOrder,  1)

	const workers = 3
	var closeChWg sync.WaitGroup
	closeChWg.Add(workers)

	for i := 0; i < workers; i ++ {
		go func () {
			defer closeChWg.Done() // worker is done 
			for order := range validateOrdersInCh {
				if order.Quantity <= 0 {
					errCh <- invalidOrder{Order: order, Err: errors.New("invalid negavtive order quantity")}
					// errCh <- invalidOrder { Order: order, Err: fmt.Errorf("invalid order quantity %f", order.Quantity) }
				} else {
					// fmt.Printf("Valid order  %v\n", order)
					validateOrdersOutCh <- order // outCh snd to ch
				}
			}
		}()
	}
	
	go func() { // spawn anonymus goroutine to close out channels
		closeChWg.Wait()
		close(validateOrdersOutCh) 	// close only rw or wo
		close(errCh) 				// close only rw or wo
	}()
	
	return validateOrdersOutCh, errCh // as ro channel

}


func receiveOrders() <-chan order {  // return ro chan
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
