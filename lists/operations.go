package lists

import (
	"context"
	"errors"
	"log"
	"reflect"
	"sort"
	"time"

	"google.golang.org/appengine/datastore"
)

//GetList loads a list from the datastore, filtering out all but the immediate children of the parent key.
func GetList(ctx context.Context, key *datastore.Key) (*List, error) {
	lst := new(List)
	var parent = new(Item)
	//Gets the parent node for its name and parent key.
	log.Println("Getting Parent")
	if err := datastore.Get(ctx, key, parent); err != nil {
		return nil, err
	}
	//if the key has a parent, it encodes the key. If not, it passes an empty string.
	if key.Parent() != nil {
		lst.Parent = key.Parent().Encode()
	} else {
		lst.Parent = ""
	}
	//then stores parent key in List struct
	lst.Name = parent.ItemText
	lst.Key = key.Encode()
	keys := make([]*datastore.Key, 0)
	//Queries the database for children of parent list. If they are a direct descendant, they are added to the list
	q := datastore.NewQuery("Item").Ancestor(key).KeysOnly()
	for i := q.Run(ctx); ; {
		k, err := i.Next(nil)
		if err == datastore.Done {
			break
		} else if err != nil {
			return nil, err
		}
		if reflect.DeepEqual(k.Parent(), key) {
			keys = append(keys, k)
		}
	}
	//then, initializes the List's items, and Gets all direct descendants using the keys list.
	lst.Items = make([]*Item, len(keys))
	for i, e := range keys {
		item := new(Item)
		if err := datastore.Get(ctx, e, item); err != nil {
			return nil, errors.New("error on get: " + err.Error())
		}
		item.Key = e.Encode()
		lst.Items[i] = item
	}
	//sorts final list by TimeAdded and returns the List.
	sort.Slice(lst.Items, func(i int, j int) bool {
		return lst.Items[i].TimeAdded.Before(lst.Items[j].TimeAdded)
	})
	return lst, nil
}

//GetRoot gets the root node of the lists. If the root doesn't exist, it instantiates it.
func GetRoot(ctx context.Context) (string, error) {
	key := datastore.NewKey(ctx, "Item", "root", 0, nil)
	root := new(Item)
	if err := datastore.Get(ctx, key, root); err == datastore.ErrNoSuchEntity {
		root.ItemText = ""
		root.TimeAdded = time.Now()
		key, err = datastore.Put(ctx, key, root)
		if err != nil {
			return "", err
		}
	}
	return key.Encode(), nil
}

//AddItem adds an item with context of text to datastore.
func AddItem(ctx context.Context, text string, parent *datastore.Key) error {
	key := datastore.NewIncompleteKey(ctx, "Item", parent)
	item := &Item{ItemText: text, TimeAdded: time.Now()}
	if _, err := datastore.Put(ctx, key, item); err != nil {
		return err
	}
	return nil
}

//DeleteItem deletes an item and all of its descendents
func DeleteItem(ctx context.Context, key *datastore.Key) (string, error) {
	parent := key.Parent()
	//queries datastore for item and ancestors.
	q := datastore.NewQuery("Item").Ancestor(key)
	for i := q.Run(ctx); ; {
		k, err := i.Next(nil)
		if err == datastore.Done {
			break
		} else if err != nil {
			return "", err
		}
		//deletes all but parent root
		if !reflect.DeepEqual(key, k) {
			if err := datastore.Delete(ctx, k); err != nil {
				return "", err
			}
		}

	}
	//deletes parent root and returns key of Parent.
	if err := datastore.Delete(ctx, key); err != nil {
		return "", err
	}
	return parent.Encode(), nil
}
