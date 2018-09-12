package testhelp

import (
	"context"
	"time"

	"github.com/is0metry/listman-gcp/lists"
	"google.golang.org/appengine/datastore"
)

func InitializeTestDatastore(ctx context.Context) (*datastore.Key, error) {
	rootKey := datastore.NewKey(ctx, "Item", "root", 0, nil)
	root := &lists.Item{ItemText: "root", TimeAdded: time.Now()}
	rootKey, err := datastore.Put(ctx, rootKey, root)
	if err != nil {
		return nil, err
	}
	return rootKey, nil
}
