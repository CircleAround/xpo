package xpo

import (
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

func (s *AppEngineService) InitAppEngineService(c context.Context) {
	s.Context = c
	s.Goon = goon.FromContext(c)
}

func (s *AppEngineService) KeyOf(obj interface{}) *datastore.Key {
	return s.Goon.Key(obj)
}

func (s *AppEngineService) Get(obj interface{}) error {
	return s.Goon.Get(obj)
}

func (s *AppEngineService) FindOrCreate(xu interface{}) (xret interface{}, err error) {
	err = datastore.RunInTransaction(s.Context, func(ctx context.Context) error {
		if err := s.Goon.Get(xu); err != nil {
			if err != datastore.ErrNoSuchEntity {
				return err
			}

			log.Infof(s.Context, "%v not found. create new one.", xu)
			_, ierr := s.Goon.Put(xu)
			if ierr != nil {
				return ierr
			}
		} else {
			log.Infof(s.Context, "%v found!. get one.", xu)
		}
		return nil
	}, nil)
	xret = xu
	return
}

func (s *AppEngineService) Put(obj interface{}) (err error) {
	_, err = s.Goon.Put(obj)
	return
}

// CreateUnique can create unique index.
//
// type UniqueIndexOfString struct {
// 	Value string `datastore:"-" goon:"id"`
// }
func (s *AppEngineService) CreateUnique(i interface{}, property string) error {
	err := s.Get(i)
	if err == nil {
		log.Infof(s.Context, "%v is not unique. %v", property, i)
		return &ValueNotUniqueError{Type: "ValueNotUniqueError", Property: property}
	}

	if err != datastore.ErrNoSuchEntity {
		return err
	}

	log.Infof(s.Context, "%v is free.", property)

	_, err = s.Goon.Put(i)
	if err != nil {
		return err
	}

	return nil
}

type ValueNotUniqueError struct {
	Type     string `json:"type"`
	Property string `json:"property"`
}

func (e *ValueNotUniqueError) Error() string {
	return "Not unique error"
}

type DuplicatedObjectError struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func (e *DuplicatedObjectError) Error() string {
	return "Duplicated object error"
}
