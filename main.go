package main

import (
	"net/http"

	"github.com/Is0metry/listman-gcp/handlers"

	"google.golang.org/appengine"
)

func main() {
	http.Handle("/", handlers.ContextHandler{handlers.RootHandler})
	http.Handle("/view/", handlers.ContextHandler{handlers.JSONViewHandler})
	http.Handle("/add/", handlers.ContextHandler{handlers.JSONAddHandler})
	http.Handle("/delete", handlers.ContextHandler{handlers.JSONDeleteHandler})
	appengine.Main()
}
