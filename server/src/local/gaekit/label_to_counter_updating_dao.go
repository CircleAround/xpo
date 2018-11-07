package gaekit

import (
	funk "github.com/thoas/go-funk"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// b and a' s array diff and return create applying Entities
// usualy Entitiy{Name, Counter}
type LabelToCounterUpdatingDao struct {
	DatastoreAccessObject

	NewFunc       func(string) interface{}
	IncrementFunc func()
	DecrementFunc func()
}

// b and a' s array diff and return create applying Entities
// usualy Entitiy{Name, Counter}
func (d *LabelToCounterUpdatingDao) BuildEntities(c context.Context, b []string, a []string) ([]interface{}, error) {
	es := []interface{}{}
	for _, value := range a {
		if funk.Contains(b, value) {
			continue
		}

		l := d.NewFunc(value)
		err := d.Get(c, l)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return nil, err
		}
		d.IncrementFunc()
		es = append(es, l)
	}

	for _, value := range b {
		if funk.Contains(a, value) {
			continue
		}

		l := d.NewFunc(value)
		err := d.Get(c, l)
		if err != nil {
			return nil, err
		}
		d.DecrementFunc()
		es = append(es, l)
	}

	return es, nil
}
