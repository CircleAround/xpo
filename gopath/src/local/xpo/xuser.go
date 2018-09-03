package xpo

import (
	"net/http"

	"google.golang.org/appengine/user"
)

// XUser struct
type XUser struct {
	ID    string `datastore:"-" goon:"id" json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type XUserService struct {
	AppEngineService
}

func NewXUserService(r *http.Request) *XUserService {
	s := new(XUserService)
	s.InitAppEngineService(r)
	return s
}

func (s *XUserService) GetOrCreate(u *user.User, name string) (xu *XUser, err error) {
	xu = &XUser{ID: u.ID, Name: name, Email: u.Email}
	xret, err := s.FindOrCreate(xu)
	xu = xret.(*XUser)
	return
}
