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
	value string `datastore:"-" goon:"id"`
}

func (i *_XUserNameUniqueIndex) Property() string {
	return "name"
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

type XUserUpdatingParams XUserBasicParams

// NewXUserService is function for construction
func NewXUserService(c context.Context) *XUserService {
	s := new(XUserService)
	s.InitAppEngineService(c)
	return s
}

// Create is method for creation XUser
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

	if err = s.Get(xu); err == nil {
		log.Infof(s.Context, "%v found!. duplicated.", xu)
		return nil, &gaekit.DuplicatedObjectError{Type: "DuplicatedObjectError"}
	}

	err = s.RunInTransaction(func() error {
		// for idempotent. if already create success, return process.
		if err = s.Get(xu); err == nil {
			return nil
		}

		if err != datastore.ErrNoSuchEntity {
			log.Infof(s.Context, "%v error.", err)
			return err
		}

		log.Infof(s.Context, "%v not found.", xu)

		i := &_XUserNameUniqueIndex{value: xu.Name}
		err = s.CreateUnique(i)
		if err != nil {
			return err
		}

		log.Infof(s.Context, "%v not found. create new one.", xu)
		return s.Put(xu)
	})
	return
}

// Update is method for updating XUser
func (s *XUserService) Update(u *user.User, params *XUserUpdatingParams) (xu *XUser, err error) {
	log.Infof(s.Context, "Update: %v", params)
	v := newValidator()
	err = v.Struct(params)
	if err != nil {
		return nil, err
	}

	err = s.RunInTransaction(func() error {
		xu, err = s.GetByUser(u)
		if err != nil {
			return err
		}

		if xu.Name == params.Name {
			if xu.Nickname == params.Nickname {
				// Maybe already succeed.
				return nil
			}
		} else {
			err = s.updateUniqueIndex(*xu, params)
			if err != nil {
				return err
			}

			xu.Name = params.Name
		}

		xu.Nickname = params.Nickname
		return s.Put(xu)
	})
	return
}

// GetByUser is method for getting XUser by user.User
func (s *XUserService) GetByUser(u *user.User) (xu *XUser, err error) {
	xu = &XUser{ID: u.ID}
	err = s.Get(xu)
	return
}

// IsUsedName is method for checking UserName already taken.
func (s *XUserService) IsUsedName(name string) (bool, error) {
	i := _XUserNameUniqueIndex{value: name}
	return s.Exists(&i)
}

func (s *XUserService) updateUniqueIndex(xu XUser, params *XUserUpdatingParams) error {
	i := &_XUserNameUniqueIndex{value: xu.Name}
	ni := &_XUserNameUniqueIndex{value: params.Name}
	return s.ChangeUniqueValueMustTr(i, ni)
}

func newValidator() *validatekit.Validate {
	v := validatekit.NewValidate()
	v.RegisterRegexValidation("username_format", `^[0-9a-z_]+$`)
	v.RegisterRegexValidation("usernickname_format", `^[0-9a-zA-Z_ぁ-んァ-ヶー一-龠]+$`)
	return v
}
