package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Is0metry/listman-gcp/lists"

	"google.golang.org/appengine/datastore"
)

func jsonReturnError(w http.ResponseWriter, r *http.Request, err error) {
	enc := json.NewEncoder(w)
	w.WriteHeader(500)
	resp := OperationResponse{Success: false, ErrorText: err.Error()}
	enc.Encode(resp)
}
func jsonReturnSuccess(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	resp := OperationResponse{Success: true, ErrorText: ""}
	enc.Encode(resp)
}

//JSONViewHandler returns the list specified in the request body.
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

//JSONAddHandler adds a new list item from the OperationRequest in the request body to the list specified by the key in the URL.
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
		if err := lists.AddItem(ctx, req.Text, key); err != nil {
			jsonReturnError(w, r, err)
			return
		}
		jsonReturnSuccess(w, r)
	} else {
		jsonReturnError(w, r, errors.New("request not of type POST"))
	}
}

//JSONDeleteHandler deletes an item specified by the key in the OperationRequest in the request body and returns an OperationResponse.
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
		_, err = lists.DeleteItem(ctx, key)
		if err != nil {
			jsonReturnError(w, r, err)
			return
		}
		jsonReturnSuccess(w, r)

	}
}
