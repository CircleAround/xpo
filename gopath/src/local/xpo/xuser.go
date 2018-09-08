package xpo

import (
	"local/apikit"
	"regexp"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

// UserNameRegex is for validation
const UserNameRegex = `^[0-9a-zA-Z_]{1,15}$`

// XUser struct
type XUser struct {
	ID       string `datastore:"-" goon:"id" json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	NickName string `json:"nickname"`
}

// _XUserNameUniqueIndex is unique index of XUser's Name
type _XUserNameUniqueIndex struct {
	Value string `datastore:"-" goon:"id"`
}

// XUserService is Service for XUser
type XUserService struct {
	AppEngineService
}

// NewXUserService is function for construction
func NewXUserService(c context.Context) *XUserService {
	s := new(XUserService)
	s.InitAppEngineService(c)
	return s
}

// Create is method for create XUser
func (s *XUserService) Create(u *user.User, name string, nickname string) (xu *XUser, err error) {
	xu = &XUser{ID: u.ID, Name: name, Email: u.Email, NickName: nickname}
	err = datastore.RunInTransaction(s.Context, func(ctx context.Context) error {
		if err = s.Get(xu); err == nil {
			log.Infof(s.Context, "%v found!. duplicated.", xu)
			return &DuplicatedObjectError{Type: "DuplicatedObjectError"}
		}
		if err != datastore.ErrNoSuchEntity {
			log.Infof(s.Context, "%v error.", err)
			return err
		}

		log.Infof(s.Context, "%v not found.", xu)

		log.Infof(s.Context, "validation start.")
		verr := apikit.NewValidationError()
		unr := regexp.MustCompile(UserNameRegex)
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
			return verr
		}

		log.Infof(s.Context, "validation end.")

		i := &_XUserNameUniqueIndex{Value: name}
		err = s.CreateUnique(i, "name")
		if err != nil {
			return err
		}

		log.Infof(s.Context, "keep name of %v.", name)

		log.Infof(s.Context, "%v not found. create new one.", xu)
		_, err := s.Goon.Put(xu)
		if err != nil {
			return err
		}

		return nil
	}, nil)
	return
}

// GetByUser is method for getting XUser by user.User
func (s *XUserService) GetByUser(u *user.User) (xu *XUser, err error) {
	xu = &XUser{ID: u.ID}
	err = s.Get(xu)
	return
}