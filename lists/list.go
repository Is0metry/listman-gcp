package lists

import "time"

type Item struct {
	Key       string `datastore:"-"`
	ItemText  string
	TimeAdded time.Time
}
type List struct {
	Name   string
	Key    string `datastore:"-"`
	Parent string `datastore:"-"`
	Items  []*Item
}
