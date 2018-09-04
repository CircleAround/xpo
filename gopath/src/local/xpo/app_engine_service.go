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
