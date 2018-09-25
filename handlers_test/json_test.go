package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Is0metry/listman-gcp/handlers"
	"github.com/Is0metry/listman-gcp/lists"
	"github.com/Is0metry/listman-gcp/testhelp"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

func TestJSONViewHandler(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Error("Error getting context from aetest")
		t.FailNow()
	}
	defer done()
	rootKey, err := testhelp.InitializeTestDatastore(ctx)
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
		t.Error("Error getting context from aetest")
		t.FailNow()
	}
	defer done()
	rootKey, _ := testhelp.InitializeTestDatastore(ctx)
	reqBody, err := json.Marshal(handlers.OperationRequest{ReqType: "add", Text: "This is a test."})
	t.Log(string(reqBody))
	bodyReader := bytes.NewReader(reqBody)
	r := httptest.NewRequest("POST", "/add/"+rootKey.Encode(), bodyReader)
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
	if len(resp.Lst.Items) != 1 {
		t.Errorf("Lst is wrong length! expected 1, got %d", len(resp.Lst.Items))
		t.FailNow()
	}
	if resp.Lst.Items[0].ItemText != "This is a test." {
		t.Errorf("Text wrong, expected \"this is a test.\" got \"%s\" instead.", resp.Lst.Items[0].ItemText)

	}
}
func TestJSONDeleteHandler(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Errorf("Unexpected error with aetest: %s", err.Error())
	}
	defer done()
	rootKey, err := testhelp.InitializeTestDatastore(ctx)
	if err != nil {
		t.Errorf("Unexpected error initializing datastore: %s", err.Error())
	}
	key := datastore.NewIncompleteKey(ctx, "Item", rootKey)
	key, err = datastore.Put(ctx, key, &lists.Item{ItemText: "to be deleted.", TimeAdded: time.Now()})
	if err != nil {
		t.Errorf("Unexpected error adding item: %s", err.Error())
	}
	req := &handlers.OperationRequest{ReqType: "delete", Text: key.Encode()}
	reqString, _ := json.Marshal(req)
	reqBody := bytes.NewReader(reqString)
	r := httptest.NewRequest("POST", "/delete", reqBody)
	w := httptest.NewRecorder()
	handlers.JSONDeleteHandler(ctx, w, r)
	if w.Code != 200 {
		t.Errorf("Wrong code! expected 200 got %d", w.Code)
		t.FailNow()
	}
	dec := json.NewDecoder(w.Body)
	var resp handlers.ListResponse
	if err := dec.Decode(&resp); err != nil {
		t.Errorf("error decoding json: %s", err.Error())
		t.FailNow()
	}
	if !resp.GetSuccessful {
		t.Errorf("Get unsuccessful: %s", resp.ErrorText)
	}
	if resp.Lst.Key != rootKey.Encode() {
		t.Errorf("Expected root key back, got %s instead", resp.Lst.Name)
	}
}
