package main

import (
	"net/http"

	"github.com/Is0metry/listman-gcp/handlers"

	"google.golang.org/appengine"
)

func main() {
	http.Handle("/", handlers.ContextHandler{handlers.RootHandler})
	http.Handle("/view/", handlers.ContextHandler{handlers.JSONViewHandler})
	http.Handle("/add/", handlers.BackgroundHandler{handlers.JSONAddHandler})
	http.Handle("/delete", handlers.BackgroundHandler{handlers.JSONDeleteHandler})
	appengine.Main()
}
