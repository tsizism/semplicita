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
	defaultPort := 8888
	port2 := flag.Int("port", defaultPort, "Web Port")
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

// this will be used to store the path to the rendered page to the read only file system
//go:embed templates
var templateFS embed.FS

func render(w http.ResponseWriter, pageFileName string) {

	templatesBase := "" // "./cmd/web/"

	partials := []string{
		templatesBase + "templates/base.layout.gohtml",
		templatesBase + "templates/header.partial.gohtml",
		templatesBase + "templates/footer.partial.gohtml",
	}

	renderedPagePath =: fmt.Sprintf(templatesBase + "templates/%s", pageFileName)
	partials = append(partials, renderedPagePath)

	fmt.Printf("partials=%v\n", partials)

	// var templateSlice []string
	// templateSlice = append(templateSlice, fmt.Sprintf("./templates/%s", page))
	// fmt.Printf("templateSlice=%v\n", templateSlice)
	// for _, x := range partials {
	// 	templateSlice = append(templateSlice, x)
	// }
	// fmt.Printf("templateSlice=%v\n", templateSlice)

	tmpl := template.New(pageFileName)

	tmpl.Funcs(template.FuncMap{
		"currentTime": currentTime, // Register the custom time function
	})

	fmt.Printf("Current time: %s\n", currentTime())

	var err error
	// tmpl, err = tmpl.ParseFiles(templateSlice...)
	// tmpl, err = tmpl.ParseFiles(partials...)

	tmpl, err := tmpl.ParseFS(templateFS, partials...)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
