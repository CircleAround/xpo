package xpo

import (
	"local/apikit"
	"net/http"
	"regexp"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

const USER_NAME_REGEX = `^[0-9a-zA-Z_]{1,15}$`

// XUser struct
type XUser struct {
	ID       string `datastore:"-" goon:"id" json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	NickName string `json:"nickname"`
}

// _XUserNameUniqueIndex is unique index of XUser's Name
type _XUserNameUniqueIndex struct {
	Name string `datastore:"-" goon:"id"`
}

type XUserService struct {
	AppEngineService
}

type AlreadyKeptNameError struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func (e *AlreadyKeptNameError) Error() string {
	return "Already kept name"
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

func (s *XUserService) Create(u *user.User, name string, nickname string) (xu *XUser, err error) {
	verr := apikit.NewValidationError()

	unr := regexp.MustCompile(USER_NAME_REGEX)
	if name == "" {
		verr.PushOne("name", apikit.Required)
	} else if !unr.MatchString(name) {
		verr.PushOne("name", apikit.InvalidFormat)
	}

	if nickname == "" {
		verr.PushOne("nickname", apikit.Required)
	}
	// TODO: screenname や　user_name の禁則文字対応

	if verr.HasItem() {
		return nil, verr
	}

	xu = &XUser{ID: u.ID, Name: name, Email: u.Email}
	err = datastore.RunInTransaction(s.Context, func(ctx context.Context) error {
		if err := s.Get(xu); err != nil {
			if err != datastore.ErrNoSuchEntity {
				return err
			}

			log.Infof(s.Context, "%v not found.", xu)

			nameUniqueIndex := &_XUserNameUniqueIndex{Name: name}
			err = s.Goon.Get(nameUniqueIndex)
			if err == nil {
				log.Infof(s.Context, "%v alreay kept.", name)
				return &AlreadyKeptNameError{Type: "AlreadyKeptNameError", Name: name}
			}

			if err != datastore.ErrNoSuchEntity {
				return err
			}

			log.Infof(s.Context, "%v is free name.", name)

			_, err = s.Goon.Put(nameUniqueIndex)
			if err != nil {
				return err
			}

			log.Infof(s.Context, "keep name of %v.", name)

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
	return
}

func (s *XUserService) GetByUser(u *user.User) (xu *XUser, err error) {
	xu = &XUser{ID: u.ID}
	err = s.Goon.Get(xu)
	return
}
