package apikit

import (
	"context"

	"google.golang.org/appengine/datastore"
)

// XUser struct
type XUser struct {
	ID    string `datastore:"-" goon:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func xUserKey(c context.Context, ID string) *datastore.Key {
	return datastore.NewKey(c, "XUser", ID, 0, nil)
}
