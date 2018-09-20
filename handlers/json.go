package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Is0metry/listman-gcp/lists"

	"google.golang.org/appengine/datastore"
)

func jsonReturnError(w http.ResponseWriter, r *http.Request, err error) {
	enc := json.NewEncoder(w)
	w.WriteHeader(500)
	resp := ListResponse{GetSuccessful: false, ErrorText: err.Error()}
	enc.Encode(&resp)
}

func JSONViewHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	keytxt := r.URL.Path[len("/view/"):]
	key, err := datastore.DecodeKey(keytxt)
	var resp ListResponse
	if err != nil {
		jsonReturnError(w, r, err)
		return
	}
	lst, err := lists.GetList(ctx, key)
	if err != nil {
		jsonReturnError(w, r, err)

	} else {
		resp.GetSuccessful = true
		resp.Lst = lst
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	enc.Encode(resp)
}

func JSONAddHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	keytxt := r.URL.Path[len("/add/"):]
	key, err := datastore.DecodeKey(keytxt)
	if err != nil {
		jsonReturnError(w, r, err)
	}
	if r.Method == "POST" {
		dec := json.NewDecoder(r.Body)
		var req OperationRequest
		if err := dec.Decode(&req); err != nil {
			jsonReturnError(w, r, err)
			return
		}
		if req.ReqType == "add" {
			if err := lists.AddItem(ctx, req.Text, key); err != nil {
				jsonReturnError(w, r, err)
				return
			}
		}
	}
	r.URL.Path = "/view/" + keytxt
	JSONViewHandler(ctx, w, r)
}

func JSONDeleteHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		dec := json.NewDecoder(r.Body)
		var req OperationRequest
		if err := dec.Decode(&req); err != nil {
			jsonReturnError(w, r, err)
		}
		key, err := datastore.DecodeKey(req.Text)
		if err != nil {
			jsonReturnError(w, r, err)
			return
		}
		parKey, err := lists.DeleteItem(ctx, key)
		if err != nil {
			jsonReturnError(w, r, err)
			return
		}
		r.URL.Path = "/view/" + parKey
		JSONViewHandler(ctx, w, r)

	}
}
