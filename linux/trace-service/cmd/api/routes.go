package main 


import (
	"net/http"

	"github.com/rs/cors"
)


func (appCtx *applicationContext) routes() http.Handler {
    mux := http.NewServeMux()

    // mux.HandleFunc("/", appCtx.handleRoot)
	mux.HandleFunc("/trace", appCtx.writeTrace)

	handler := cors.Default().Handler(mux)

	/* Error: 
	   Access to fetch at 'http://localhost:8080/handle' from origin 'http://localhost' has been blocked by CORS policy: 
	   Response to preflight request doesn't pass access control check: No 'Access-Control-Allow-Origin' header is present on the requested resource. 
	   If an opaque response serves your needs, set the request's mode to 'no-cors' to fetch the resource with CORS disabled.
	*/

	c := cors.New(cors.Options{
	 	//AllowedOrigins: []string{"https://*", "http://*"},   
		 AllowedOrigins: []string{"http://localhost", "http://localhost:80", "http://localhost:8888"},   // we allow to call "http://localhost:8080/handle" from http://localhost , (in adition to http://localhost:8080 - which would be the same origin as http://localhost/handle)
	})

	handler = c.Handler(handler)

    return handler
}
