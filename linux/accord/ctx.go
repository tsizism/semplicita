package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

type result struct {
	url string
	err error
	latency time.Duration
}

// res := new(result)
// fmt.Printf("%v, %p\n", res, res)

var urls = []string{"https://www.amazon.ca/", "https://www.google.ca/", "https://www.nytimes.com/", "https://www.wsj.com/", "http://localhost:8080"}

// https://www.youtube.com/watch?v=0x_oUlxzw5A

func main3() {
	s := time.Now(); defer func(){t := time.Since(s).Round(time.Millisecond);fmt.Printf("Duration=%s", t)}()
	
	// allUrls()

	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Millisecond);	defer cancel()
	r, _ := firstUrl(ctx)

	if r.err != nil {
		log.Printf("%-20s %s %s\n", r.url, r.latency, r.err)
	} else {
		log.Printf("%-20s %s\n", r.url, r.latency)
	}

	log.Println("sleep 9 sec")
	time.Sleep(9 * time.Second)

	log.Println("quiting anyway ...", runtime.NumGoroutine(), "still running")
}

func firstUrl(ctx context.Context) (*result, error) {
	resultsCh := make(chan result, len(urls)) 
	ctx, cancel :=  context.WithCancel(ctx); defer cancel()

	for _, url := range urls {
		go get (ctx, url, resultsCh)
	}

	select {
	case <-ctx.Done(): // Timeout 
		fmt.Printf("select in action Done %v\n", ctx.Err())
		t, _ := ctx.Deadline()
		fmt.Printf("Deadline= %v\n", t)
		return &result{}, ctx.Err()

	case r := <- resultsCh:
		fmt.Printf("select in action %v\n", r)
		return &r, nil
	}


}

func allUrls()  {
	resultsCh := make(chan result) 

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second);	defer cancel()
	
	for _, url := range urls {
		go get (ctx, url, resultsCh)
	}

	for range urls {
		r := <-resultsCh // blocking 4 times

		if r.err != nil {
			log.Printf("%-20s %s %s\n", r.url, r.latency, r.err)
		} else {
			log.Printf("%-20s %s\n", r.url, r.latency)
		}
	}
}

func get(ctx context.Context,url string, woCh chan<- result) {
	tickerCh := time.NewTicker(1 * time.Second).C

	s := time.Now()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		resp.Body.Close()
	}

	t := time.Since(s).Round(time.Millisecond)		
	r := result{url, err, t}

	fmt.Println("Posting result on chWO for " + url)
	for {
		select {
			case woCh <- r: // bloking 
				fmt.Println("Posted result on chWO for " + url)
				return 
			case <- tickerCh: 
				log.Println("tick", r)
		}
	}
}


// func get(url string, chWO chan<- result) {
// 	s := time.Now(); 
   
//    resp, err := http.Get(url)
//    t := time.Since(s).Round(time.Millisecond)
//    chWO <- result{url, err, t}
//    resp.Body.Close()
//    if err == nil {
// 	   resp.Body.Close()
//    }
// }