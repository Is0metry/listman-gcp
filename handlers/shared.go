package handlers

import (
	"context"
	"net/http"

	"github.com/Is0metry/listman-gcp/lists"
	"google.golang.org/appengine"
)

type ContextHandler struct {
	Real func(ctx context.Context, w http.ResponseWriter, r *http.Request)
}

func (f ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	f.Real(ctx, w, r)
}
func RootHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	key, err := lists.GetRoot(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+key, http.StatusFound)
}
