package gaekit

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type EntityCreator struct {
	DatastoreAccessObject
	OnBeforeCreate func(c context.Context) error
	OnStartCreate  func(c context.Context) error
	OnCreate       func(c context.Context, e interface{}) error
}

func NewEntityCreator() *EntityCreator {
	ec := new(EntityCreator)
	ec.OnBeforeCreate = func(c context.Context) error { return nil }
	ec.OnStartCreate = func(c context.Context) error { return nil }
	ec.OnCreate = func(c context.Context, e interface{}) error { return ec.Put(c, e) }
	return ec
}

func (r *EntityCreator) Execute(c context.Context, e interface{}) (err error) {
	err = r.OnBeforeCreate(c)
	if err != nil {
		return errors.Wrap(err, "failed OnBeforeCreate")
	}

	return r.RunInXGTransaction(c, func(ctx context.Context) error {
		if r.OnStartCreate != nil {
			err = r.OnStartCreate(ctx)
			if err != nil {
				return errors.Wrap(err, "failed OnStartCreate")
			}
		}

		if r.OnCreate != nil {
			err = r.OnCreate(ctx, e)
			if err != nil {
				return errors.Wrap(err, "failed OnCreate")
			}
		} else {
			err = r.Put(ctx, e)
			if err != nil {
				return errors.Wrap(err, "failed Put")
			}
		}
		return nil
	})
}
