package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/is0metry/listman-gcp/handlers"
	"github.com/is0metry/listman-gcp/lists"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

func initializeTestDatastore(ctx context.Context) (*datastore.Key, error) {
	rootKey := datastore.NewKey(ctx, "Item", "root", 0, nil)
	root := &lists.Item{ItemText: "root", TimeAdded: time.Now()}
	rootKey, err := datastore.Put(ctx, rootKey, root)
	if err != nil {
		return nil, err
	}
	return rootKey, nil
}

func TestJSONViewHandler(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Error("Error getting context from aetest")
		t.FailNow()
	}
	defer done()
	rootKey, err := initializeTestDatastore(ctx)
	if err := lists.AddItem(ctx, "this is a test item", rootKey); err != nil {
		t.Error("Error adding item to database")
		t.FailNow()
	}
	r := httptest.NewRequest("GET", "/view/"+rootKey.Encode(), nil)
	w := httptest.NewRecorder()
	handlers.JSONViewHandler(ctx, w, r)
	if w.Code != 200 {
		t.Errorf("Unexpected response code, expected 200, got %d instead", w.Code)
	}
	var lst handlers.ListResponse
	json.Unmarshal(w.Body.Bytes(), &lst)
	if len(lst.Lst.Items) != 1 {
		t.Errorf("lst.Items length wrong! expected 1 got %d instead", len(lst.Lst.Items))
		t.Log(w.Body.String())
	} else if lst.Lst.Items[0].ItemText != "this is a test item" {
		t.Errorf("item wrong! expected \"this is a test item\" got \"%s\" instead", lst.Lst.Items[0].ItemText)
	}
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/view/9c4685ae56853575bbace", nil)
	handlers.JSONViewHandler(ctx, w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("wrong code! expected %d found %d instead", http.StatusInternalServerError, w.Code)
	}
}
func TestJSONAddHandler(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Error("Error getting context from aetest")8
		t.FailNow()
	}
	defer done()
	rootKey, _ := initializeTestDatastore(ctx)
	reqBody, err := json.Marshal(handlers.OperationRequest{ReqType: "add", Text: "This is a test."})
	t.Log(reqBody)
	bodyReader := bytes.NewReader(reqBody)
	r := httptest.NewRequest("GET", "/add/"+rootKey.Encode(), bodyReader)
	w := httptest.NewRecorder()
	handlers.JSONAddHandler(ctx, w, r)
	var resp handlers.ListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("Error unmarshalling json with error %s.", err.Error())
		t.Log(w.Body.String())
		t.FailNow()
	}
	if !resp.GetSuccessful {
		t.Errorf("Response unsuccessful with message %s.", resp.ErrorText)
		t.Log(w.Body.String())

	}
}
