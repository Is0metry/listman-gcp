package handlers

import (
	"context"
	"net/http"

	"github.com/Is0metry/listman-gcp/lists"
	"google.golang.org/appengine"
)

//ContextHandler takes a handler which requires a context, generates the context and runs the function with the provided handler.
type ContextHandler struct {
	Real func(ctx context.Context, w http.ResponseWriter, r *http.Request)
}

//ServeHTTP generates a new context and runs the Real handler to serve the Request.
func (f ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	f.Real(ctx, w, r)
}

//BackgroundHandler is a ContextHandler which presents the function with a background context.
type BackgroundHandler struct {
	Real func(ctx context.Context, w http.ResponseWriter, r *http.Request)
}

//ServeHTTP generates a new background context and runs the Real http handler.
func (f BackgroundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	f.Real(ctx, w, r)
}

//RootHandler returns the root of the
func RootHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	key, err := lists.GetRoot(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+key, http.StatusFound)
}
