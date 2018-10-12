package gaekit

import (
	"fmt"

	"github.com/mjibson/goon"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// AppEngineService is Basi Service of AppEngine
type AppEngineService struct {
}

func (s *AppEngineService) Goon(c context.Context) *goon.Goon {
	return goon.FromContext(c)
}

// KeyOf is a method for getting key from obj
func (s *AppEngineService) KeyOf(c context.Context, obj interface{}) *datastore.Key {
	return s.Goon(c).Key(obj)
}

// Get is a method for retriving object
func (s *AppEngineService) Get(c context.Context,obj interface{}) error {
	return s.Goon(c).Get(obj)
}

// Exists is a method for check object exists on DB
func (s *AppEngineService) Exists(c context.Context, obj interface{}) (bool, error) {
	err := s.Get(c, obj)
	if err == datastore.ErrNoSuchEntity {
		return false, nil
	}
	return true, err
}

func (s *AppEngineService) RunInTransaction(c context.Context, process func(context.Context) error) error {
	return s.RunInTransactionWithOption(c, process, nil)
}

func (s *AppEngineService) RunInXGTransaction(c context.Context, process func(context.Context) error) error {
  opt := &datastore.TransactionOptions{XG: true}
	return s.RunInTransactionWithOption(c, process, opt)
}

func (s *AppEngineService) RunInTransactionWithOption(c context.Context, process func(context.Context) error, opts *datastore.TransactionOptions) error {
	return datastore.RunInTransaction(c, process, opts)
}

// FindOrCreate is a method for find or create Object
func (s *AppEngineService) FindOrCreate(ctx context.Context, obj interface{}) (xret interface{}, err error) {
	err = s.RunInTransaction(ctx, func(c context.Context) error {
		if err := s.Get(c, obj); err != nil {
			if err != datastore.ErrNoSuchEntity {
				return err
			}

			log.Infof(c, "%v not found. create new one.", obj)
			return s.Put(c, obj)
		} else {
			log.Infof(c, "%v found!. get one.", obj)
		}
		return nil
	})
	xret = obj
	return
}

// Put is a method for saving obj
func (s *AppEngineService) Put(c context.Context, obj interface{}) (err error) {
	_, err = s.Goon(c).Put(obj)
	return
}

// Delete is a method for deleting obj
func (s *AppEngineService) Delete(c context.Context, obj interface{}) (err error) {
	return s.Goon(c).Delete(s.KeyOf(c, obj))
}

type UniqueIndex interface {
	Property() string
}

// CreateUnique can create unique index.
//
// type UniqueIndexOfString struct {
// 	value string `datastore:"-" goon:"id"`
// }
func (s *AppEngineService) CreateUnique(c context.Context, i UniqueIndex) error {
	return s.CreateUniqueWithProperty(c, i, i.Property())
}

func (s *AppEngineService) CreateUniqueWithProperty(c context.Context, i interface{}, property string) error {
	err := s.Get(c, i)
	if err == nil {
		log.Infof(c, "%v is not unique. %v", property, i)
		return &ValueNotUniqueError{Type: "valueNotUnique", Property: property}
	}

	if err != datastore.ErrNoSuchEntity {
		return err
	}

	log.Infof(c, "%v is free.", property)
	return s.Put(c, i)
}

func (s *AppEngineService) ChangeUniqueValueMustTr(c context.Context, i UniqueIndex, ni UniqueIndex) error {
	if i.Property() != ni.Property() {
		return fmt.Errorf("Property not match: %v and %v", i.Property(), ni.Property())
	}

	log.Infof(c, "CreateUnique")
	err := s.CreateUnique(c, ni)
	if err != nil {
		return err
	}

	log.Infof(c, "Get")
	err = s.Get(c, i)
	if err == nil {
		err = s.Delete(c, i)
		if err != nil {
			return err
		}
	} else if err != datastore.ErrNoSuchEntity {
		return err
	}

	log.Infof(c, "end ChangeUniqueValueMustTr")
	return nil
}

// ValueNotUniqueError is a struct for a error value not unique
type ValueNotUniqueError struct {
	Type     string `json:"type"`
	Property string `json:"property"`
}

// Error is a method for creating message
func (e *ValueNotUniqueError) Error() string {
	return "Not unique error"
}

// DuplicatedObjectError is a struc for a error when object already exist
type DuplicatedObjectError struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// Error is am method for creating message
func (e *DuplicatedObjectError) Error() string {
	return "Duplicated object error"
}
