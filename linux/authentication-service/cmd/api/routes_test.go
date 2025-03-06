package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
)

// go test -v .
func Test_routes_exis(t *testing.T) {
	appCtx := &applicationContext{
		cfg: config{},
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}

	serveMux := appCtx.routes().(*http.ServeMux)

	routes := []string {"/authenticate", "/bla"}

	for i, routePattern := range routes {
		res := routeExists(serveMux, routePattern)	

		if i % 2 == 0 && !res{
			t.Errorf("Route %s is not found", routePattern)
			continue
		}

		if i % 2 == 1 && res {
			t.Errorf("Route %s is found", routePattern)
			continue
		}
	}
}

func routeExists(mux *http.ServeMux, routePattern string) bool {
	request, _ := http.NewRequest("GET", routePattern, nil)

	h, pattern := mux.Handler(request)

	fmt.Printf("pattern=%s,h=%v\n", pattern, h)

	return pattern == routePattern
}