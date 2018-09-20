package lists_test

import (
	"testing"
	"time"

	"github.com/Is0metry/listman-gcp/lists"
	"github.com/Is0metry/listman-gcp/testhelp"
	"github.com/icrowley/fake"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

func TestGetList(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Error("Error getting context from aetest")
	}
	defer done()
	rootKey, err := testhelp.InitializeTestDatastore(ctx)
	if err != nil {
		t.Error("error initializing datastore")
	}
	lst, err := lists.GetList(ctx, rootKey)
	if len(lst.Items) > 0 {
		t.Error("Returned non-empty list when should have been empty")
	}
	refList := &lists.List{Parent: "", Key: rootKey.Encode()}
	refList.Items = make([]*lists.Item, 10)
	for i := 0; i < 10; i++ {
		k := datastore.NewIncompleteKey(ctx, "Item", rootKey)
		item := &lists.Item{ItemText: fake.Sentence(), TimeAdded: time.Now()}
		k, err = datastore.Put(ctx, k, item)
		if err != nil {
			t.Error("Error inserting into datastore, " + err.Error())
		}
		item.Key = k.Encode()
		refList.Items[i] = item
	}
	lst, err = lists.GetList(ctx, rootKey)
	if err != nil {
		t.Error("Error on GetList: " + err.Error())
	}
	if len(refList.Items) != len(lst.Items) {
		t.Errorf("Lengths are not the same! refList:%d, lst:%d", len(refList.Items), len(lst.Items))
	}
	for i, itm := range refList.Items {
		if itm.Key != lst.Items[i].Key {
			t.Errorf("Item Keys not the same! %s expected, %s found", itm.Key, lst.Items[i].Key)
		}
		if itm.ItemText != lst.Items[i].ItemText {
			t.Errorf("ItemText not the same! %s expected, %s found", itm.ItemText, lst.Items[i].ItemText)
		}
		if itm.TimeAdded.UTC().Round(time.Millisecond) != lst.Items[i].TimeAdded.Round(time.Millisecond) {
			t.Errorf("Item TimeAdded not the same! %s expected, %s found", itm.TimeAdded.UTC().Round(time.Millisecond), lst.Items[i].TimeAdded)
		}
	}
	sublistParentKey, err := datastore.DecodeKey(refList.Items[3].Key)
	if err != nil {
		t.Error("IDKWTF but decoding the key gave an error")
	}
	subKey := datastore.NewIncompleteKey(ctx, "Item", sublistParentKey)
	sublistItem := &lists.Item{ItemText: "This is a sublist item", TimeAdded: time.Now()}
	refSublist := &lists.List{Key: subKey.Encode(), Parent: rootKey.Encode()}
	refSublist.Items = []*lists.Item{sublistItem}
	subKey, err = datastore.Put(ctx, subKey, sublistItem)
	if err != nil {
		t.Error("error inserting sublist into datastore: " + err.Error())
	}
	sublist, err := lists.GetList(ctx, sublistParentKey)
	if err != nil {
		t.Error("error getting sublist from datastore: " + err.Error())
	}
	if len(sublist.Items) != len(refSublist.Items) {
		t.Errorf("sublist length wrong! %d expected, %d found", len(refSublist.Items), len(sublist.Items))
	}
	if sublist.Items[0].ItemText != "This is a sublist item" {
		t.Errorf("sublist ItemText wrong! \"This is a sublist item\" expected \"%s\" found.", sublist.Items[0].ItemText)
	}
	if err = datastore.Delete(ctx, sublistParentKey); err != nil {
		t.Error("error deleting sublistParent for nonexistent list test")
	}
	if _, err = lists.GetList(ctx, sublistParentKey); err == nil {
		t.Error("nonexistent sublist found. Expected ErrNoSuchEntity got nil instead.")
	}
	if err != datastore.ErrNoSuchEntity {
		t.Errorf("Unexpected error. Expected %s got %s instead", datastore.ErrNoSuchEntity.Error(), err.Error())
	}
}
func TestAddList(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Error("Error getting aetest context")
	}
	defer done()

	rootKey, err := testhelp.InitializeTestDatastore(ctx)
	if err != nil {
		t.Error("error initializing datastore")
	}
	if err = lists.AddItem(ctx, "", rootKey); err == nil {
		t.Error("Inserted empty item into list. Expected error, got nil instead!")
	} else if err != lists.ErrEmptyItem {
		t.Errorf("Unexpected error! Expected %s got %s instead", lists.ErrEmptyItem.Error(), err.Error())
	}
	if err = lists.AddItem(ctx, "		    ", rootKey); err == nil {
		t.Error("Inserted empty item into list! expected error, got nil instead")
	} else if err != lists.ErrEmptyItem {
		t.Errorf("Unexpected error! Expected %s, got %s instead", lists.ErrEmptyItem.Error(), err.Error())
	}
	if err = lists.AddItem(ctx, "This item should work", rootKey); err != nil {
		t.Errorf("Unexpected error! Expected nil, got %s instead", err.Error())
		t.FailNow()
	}
	lst, err := lists.GetList(ctx, rootKey)
	if err != nil {
		t.Errorf("Unexpected error getting list! expected nil, got %s instead", err.Error())
		t.FailNow()
	} else if len(lst.Items) != 1 {
		t.Errorf("lst.Items length wrong! expected 1, got %d instead.", len(lst.Items))
		t.FailNow()
	} else if lst.Items[0].ItemText != "This item should work" {
		t.Errorf("List item incorrect, expected \"This item should work,\" got \"%s\" instead", lst.Items[0].ItemText)
	}
	if err = lists.AddItem(ctx, "Hi I'm booby and I'm an asshole tryna break this shit", rootKey); err != nil {
		t.Errorf("Unexpected error! expected nil got %s instead", err.Error())
	}
	lst, err = lists.GetList(ctx, rootKey)
	if err != nil {
		t.Errorf("Unexpected error! expected nil got %s instead", err.Error())
		t.FailNow()
	} else if len(lst.Items) != 2 {
		t.Errorf("lst.Items length wrong! expected 2, got %d instead.", len(lst.Items))
		t.FailNow()
	} else if lst.Items[1].ItemText != "Hi I'm booby and I'm an asshole tryna break this shit" {
		t.Errorf("List item incorrect, expected \"Hi I'm booby and I'm an asshole tryna break this shit,\" got \"%s\" instead", lst.Items[1].ItemText)
	}
}
