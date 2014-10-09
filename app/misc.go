package main

import (
	"appengine"
	"appengine/datastore"
)

const (
	OP_LIMIT = 500
)

// PartPutMulti performs a number of small PutMulti operations over a large
// set of input data.
func PartPutMulti(c appengine.Context, keys []*datastore.Key, ents []interface{}) error {
	size := OP_LIMIT
	for len(keys) > 0 {
		if size > len(keys) {
			size = len(keys)
		}
		if _, err := datastore.PutMulti(c, keys[:size], ents[:size]); err != nil {
			return err
		}
		keys = keys[size:]
		ents = ents[size:]
	}
	return nil
}

// PartDeleteMulti performs a number of small DeleteMulti operations over a
// large set of keys to delete.
func PartDeleteMulti(c appengine.Context, keys []*datastore.Key) error {
	size := OP_LIMIT
	for len(keys) > 0 {
		if size > len(keys) {
			size = len(keys)
		}
		if err := datastore.DeleteMulti(c, keys[:size]); err != nil {
			return err
		}
		keys = keys[size:]
	}
	return nil
}
