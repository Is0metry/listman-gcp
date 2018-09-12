package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/is0metry/listman-gcp/lists"

	"google.golang.org/appengine/datastore"
)

func JSONViewHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	keytxt := r.URL.Path[len("/view/"):]
	key, err := datastore.DecodeKey(keytxt)
	var resp ListResponse
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.GetSuccessful = false
		resp.ErrorText = err.Error()
	} else {
		lst, err := lists.GetList(ctx, key)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp.GetSuccessful = false
			resp.ErrorText = err.Error()
		} else {
			resp.GetSuccessful = true
			resp.Lst = lst
		}
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	enc.Encode(resp)
}
func JSONAddHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	keytxt := r.URL.Path[len("/add/"):]
	key, err := datastore.DecodeKey(keytxt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if r.Method == "POST" {
		dec := json.NewDecoder(r.Body)
		var req OperationRequest
		if err := dec.Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if req.ReqType == "add" {
			if err := lists.AddItem(ctx, req.Text, key); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
	JSONViewHandler(ctx, w, r)
}
