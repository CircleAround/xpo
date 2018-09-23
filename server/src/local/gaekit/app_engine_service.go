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
	Goon    *goon.Goon
	Context context.Context
}

// InitAppEngineService is initialzer for Extended Services
func (s *AppEngineService) InitAppEngineService(c context.Context) {
	s.Context = c
	s.Goon = goon.FromContext(c)
}

// KeyOf is a method for getting key from obj
func (s *AppEngineService) KeyOf(obj interface{}) *datastore.Key {
	return s.Goon.Key(obj)
}

// Get is a method for retriving object
func (s *AppEngineService) Get(obj interface{}) error {
	return s.Goon.Get(obj)
}

// Exists is a method for check object exists on DB
func (s *AppEngineService) Exists(obj interface{}) (bool, error) {
	err := s.Get(obj)
	if err == datastore.ErrNoSuchEntity {
		return false, nil
	}
	return true, err
}

func (s *AppEngineService) RunInTransaction(process func() error) error {
	return s.RunInTransactionWithOption(process, nil)
}

func (s *AppEngineService) RunInTransactionWithOption(process func() error, opts *datastore.TransactionOptions) error {
	return datastore.RunInTransaction(s.Context, func(ctx context.Context) error {
		return process()
	}, opts)
}

// FindOrCreate is a method for find or create Object
func (s *AppEngineService) FindOrCreate(obj interface{}) (xret interface{}, err error) {
	err = s.RunInTransaction(func() error {
		if err := s.Goon.Get(obj); err != nil {
			if err != datastore.ErrNoSuchEntity {
				return err
			}

			log.Infof(s.Context, "%v not found. create new one.", obj)
			return s.Put(obj)
		} else {
			log.Infof(s.Context, "%v found!. get one.", obj)
		}
		return nil
	})
	xret = obj
	return
}

// Put is a method for saving obj
func (s *AppEngineService) Put(obj interface{}) (err error) {
	_, err = s.Goon.Put(obj)
	return
}

// Delete is a method for deleting obj
func (s *AppEngineService) Delete(obj interface{}) (err error) {
	return s.Goon.Delete(s.KeyOf(obj))
}

type UniqueIndex interface {
	Property() string
}

// CreateUnique can create unique index.
//
// type UniqueIndexOfString struct {
// 	value string `datastore:"-" goon:"id"`
// }
func (s *AppEngineService) CreateUnique(i UniqueIndex) error {
	return s.CreateUniqueWithProperty(i, i.Property())
}

func (s *AppEngineService) CreateUniqueWithProperty(i interface{}, property string) error {
	err := s.Get(i)
	if err == nil {
		log.Infof(s.Context, "%v is not unique. %v", property, i)
		return &ValueNotUniqueError{Type: "ValueNotUniqueError", Property: property}
	}

	if err != datastore.ErrNoSuchEntity {
		return err
	}

	log.Infof(s.Context, "%v is free.", property)
	return s.Put(i)
}

func (s *AppEngineService) ChangeUniqueValueMustTr(i UniqueIndex, ni UniqueIndex) error {
	if i.Property() != ni.Property() {
		return fmt.Errorf("Property not match: %v and %v", i.Property(), ni.Property())
	}

	log.Infof(s.Context, "CreateUnique")
	err := s.CreateUnique(ni)
	if err != nil {
		return err
	}

	log.Infof(s.Context, "Get")
	err = s.Get(i)
	if err == nil {
		err = s.Delete(i)
		if err != nil {
			return err
		}
	} else if err != datastore.ErrNoSuchEntity {
		return err
	}

	log.Infof(s.Context, "end ChangeUniqueValueMustTr")
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
