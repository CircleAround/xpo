package xpo

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"local/gaekit"
	"local/validatekit"
)

// XUser struct
type XUser struct {
	ID       string `datastore:"-" goon:"id" json:"id"`
	Name     string `json:"name" validate:"required,max=15,username_format"`
	Email    string `json:"email" validate:"required"`
	Nickname string `json:"nickname" validate:"required,max=128,usernickname_format"`
}

// _XUserNameUniqueIndex is unique index of XUser's Name
type _XUserNameUniqueIndex struct {
	Value string `datastore:"-" goon:"id"`
}

// XUserService is Service for XUser
type XUserService struct {
	gaekit.AppEngineService
}

// XUserBasicParams is parameter's basic
type XUserBasicParams struct {
	Name     string `json:"name" validate:"required,max=15,username_format"`
	Nickname string `json:"nickname" validate:"required,max=128,usernickname_format"`
}


// XUserCreationParams is parameter of Create
type XUserCreationParams XUserBasicParams

type XUserUpdatingParams struct{
	AuthorID string `json:"author_id" validate:"required`
	ID string `json:"id" validate:"required"`
	XUserCreationParams
}

// NewXUserService is function for construction
func NewXUserService(c context.Context) *XUserService {
	s := new(XUserService)
	s.InitAppEngineService(c)
	return s
}

// Create is method for create XUser
func (s *XUserService) Create(u *user.User, params *XUserCreationParams) (xu *XUser, err error) {
	log.Infof(s.Context, "Create: %v", params)
	v := newValidator()
	err = v.Struct(params)
	if err != nil {
		return nil, err
	}

	xu = &XUser{ID: u.ID, Name: params.Name, Email: u.Email, Nickname: params.Nickname}
	err = v.Struct(xu)
	if err != nil {
		return nil, err
	}

	err = datastore.RunInTransaction(s.Context, func(ctx context.Context) error {
		if err = s.Get(xu); err == nil {
			log.Infof(s.Context, "%v found!. duplicated.", xu)
			return &gaekit.DuplicatedObjectError{Type: "DuplicatedObjectError"}
		}
		if err != datastore.ErrNoSuchEntity {
			log.Infof(s.Context, "%v error.", err)
			return err
		}

		log.Infof(s.Context, "%v not found.", xu)

		i := &_XUserNameUniqueIndex{Value: xu.Name}
		err = s.CreateUnique(i, "name")
		if err != nil {
			return err
		}

		log.Infof(s.Context, "keep name of %v.", xu.Name)

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

func newValidator() *validatekit.Validate {
	v := validatekit.NewValidate()
	v.RegisterRegexValidation("username_format", `^[0-9a-z_]+$`)
	v.RegisterRegexValidation("usernickname_format", `^[0-9a-zA-Z_ぁ-んァ-ヶー一-龠]+$`)
	return v
}
