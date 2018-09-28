package xpo

import (
	"bufio"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"local/apikit"
	"local/gaekit"
	"local/stdkit"
	"local/validatekit"
	"local/xpo/assets"
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

// XUserProfileParams is parameter's basic
type XUserProfileParams struct {
	Name     string `json:"name" validate:"required,max=15,username_format"`
	Nickname string `json:"nickname" validate:"required,max=128,usernickname_format"`
}

// NewXUserService is function for construction
func NewXUserService() *XUserService {
	s := new(XUserService)
	return s
}

// Create is method for creation XUser
func (s *XUserService) Create(c context.Context, u user.User, params XUserProfileParams) (xu *XUser, err error) {
	log.Infof(c, "Create: %v", params)
	v, err := validate(params)
	if err != nil {
		return nil, err
	}

	xu = &XUser{ID: u.ID, Name: params.Name, Email: u.Email, Nickname: params.Nickname}
	err = v.Struct(xu)
	if err != nil {
		return nil, err
	}

	if err = s.Get(c, xu); err == nil {
		log.Infof(c, "%v found!. duplicated.", xu)
		return nil, &gaekit.DuplicatedObjectError{Type: "DuplicatedObjectError"}
	}

	if err != datastore.ErrNoSuchEntity {
		log.Infof(c, "%v error.", err)
		return nil, err
	}

	err = s.RunInTransaction(c, func() error {
		// for idempotent. if already create success, return process.
		if err = s.Get(c, xu); err == nil {
			return nil
		}

		if err != datastore.ErrNoSuchEntity {
			log.Infof(c, "%v error.", err)
			return err
		}

		log.Infof(c, "%v not found.", xu)

		i := &_XUserNameUniqueIndex{value: xu.Name}
		err = s.CreateUnique(c, i)
		if err != nil {
			return err
		}

		log.Infof(c, "%v not found. create new one.", xu)
		return s.Put(c, xu)
	})
	return
}

// Update is method for updating XUser
func (s *XUserService) Update(c context.Context, u user.User, params XUserProfileParams) (xu *XUser, err error) {
	log.Infof(c, "Update: %v", params)
	_, err = validate(params)
	if err != nil {
		return nil, err
	}

	err = s.RunInTransaction(c, func() error {
		xu, err = s.GetByUser(c, u)
		if err != nil {
			return err
		}

		if xu.Name == params.Name {
			if xu.Nickname == params.Nickname {
				// Maybe already succeed.
				return nil
			}
		} else {
			err = s.updateUniqueIndex(c, *xu, params)
			if err != nil {
				return err
			}

			xu.Name = params.Name
		}

		xu.Nickname = params.Nickname
		return s.Put(c, xu)
	})
	return
}

// GetByUser is method for getting XUser by user.User
func (s *XUserService) GetByUser(c context.Context, u user.User) (xu *XUser, err error) {
	xu = &XUser{ID: u.ID}
	err = s.Get(c, xu)
	return
}

// IsUsedName is method for checking UserName already taken.
func (s *XUserService) IsUsedName(c context.Context, name string) (bool, error) {
	i := _XUserNameUniqueIndex{value: name}
	return s.Exists(c, &i)
}

func validate(params XUserProfileParams) (*validatekit.Validate, error) {
	v := newValidator()
	err := v.Struct(params)
	if err != nil {
		return nil, err
	}

	var prop string
	ng, err := checkBlockedWord(func(word string) bool {
		if params.Name == word {
			prop = "name"
			return true
		}
		if params.Nickname == word {
			prop = "nickname"
			return true
		}
		return false
	})

	if err != nil {
		return nil, err
	}

	if ng {
		return nil, apikit.NewInvalidParameterError(prop)
	}

	return v, nil
}

func (s *XUserService) updateUniqueIndex(c context.Context, xu XUser, params XUserProfileParams) error {
	i := &_XUserNameUniqueIndex{value: xu.Name}
	ni := &_XUserNameUniqueIndex{value: params.Name}
	return s.ChangeUniqueValueMustTr(c, i, ni)
}

func newValidator() *validatekit.Validate {
	v := validatekit.NewValidate()
	v.RegisterRegexValidation("username_format", `^[0-9a-z_]+$`)
	v.RegisterRegexValidation("usernickname_format", `^[0-9a-zA-Z_ぁ-んァ-ヶー一-龠]+$`)
	return v
}

func checkBlockedWord(callback func(line string) bool) (bool, error) {
	f, err := assets.Assets.Open("/assets/reserved_username_list")
	if err != nil {
		return false, err
	}

	defer f.Close()

	reader := bufio.NewReaderSize(f, 128)
	hit, err := stdkit.FindLine(reader, callback)
	return hit, err
}
