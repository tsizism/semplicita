package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func currentTime() string {
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("now=%s\n", now)
	return now
}

// $go build -v -a -o frontApp ./cmd/web
// ./frontApp -port 8888
func main() {
	port2 := flag.Int("port", 80, "Web Port")
	flag.Parse()
	port := *port2

	fmt.Printf("Starting front-end service on port %d\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.html")
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	// var err error
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, page string) {

	partials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	partials = append(partials, fmt.Sprintf("./cmd/web/templates/%s", page))

	fmt.Printf("partials=%v\n", partials)

	// var templateSlice []string
	// templateSlice = append(templateSlice, fmt.Sprintf("./templates/%s", page))
	// fmt.Printf("templateSlice=%v\n", templateSlice)
	// for _, x := range partials {
	// 	templateSlice = append(templateSlice, x)
	// }
	// fmt.Printf("templateSlice=%v\n", templateSlice)

	tmpl := template.New(page)

	tmpl.Funcs(template.FuncMap{
		"currentTime": currentTime, // Register the custom time function
	})

	fmt.Printf("Current time: %s\n", currentTime())

	var err error
	// tmpl, err = tmpl.ParseFiles(templateSlice...)
	tmpl, err = tmpl.ParseFiles(partials...)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
