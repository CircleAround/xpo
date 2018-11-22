package app

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"local/validatekit"
	"local/xpo/domain"
	"local/xpo/entities"
	"local/xpo/store"
)

// XUserService is Service for XUser
type XUserService struct {
	urep  *store.XUserRepository
	ulrep *store.XUserLanguageRepository
}

// NewXUserService is function for construction
func NewXUserService() *XUserService {
	s := new(XUserService)
	s.urep = store.NewXUserRepository()
	s.ulrep = store.NewXUserLanguageRepository()
	return s
}

// Create is method for creation XUser
func (s *XUserService) Create(c context.Context, u user.User, params entities.XUserProfileParams) (xu *entities.XUser, err error) {
	log.Infof(c, "Create: %v", params)
	v, err := validate(params)
	if err != nil {
		return nil, errors.Wrap(err, "params Validation failed")
	}

	xu = &entities.XUser{ID: u.ID, Name: params.Name, Email: u.Email, Nickname: params.Nickname}
	err = v.Struct(xu)
	if err != nil {
		return nil, errors.Wrap(err, "xu Validation failed")
	}

	err = s.urep.Create(c, xu)
	if err != nil {
		return nil, errors.Wrap(err, "Create failed")
	}
	return
}

// Update is method for updating entities.XUser
func (s *XUserService) Update(c context.Context, xu *entities.XUser, params entities.XUserProfileParams) (*entities.XUser, error) {
	log.Infof(c, "Update: %v", params)
	_, err := validate(params)
	if err != nil {
		return nil, errors.Wrap(err, "params Validation failed")
	}
	err = s.urep.Update(c, xu, params)
	if err != nil {
		return nil, errors.Wrap(err, "Update failed")
	}
	return xu, err
}

// GetByUser is method for getting entities.XUser by user.User
func (s *XUserService) GetByUser(c context.Context, u user.User) (xu *entities.XUser, err error) {
	return s.GetByID(c, u.ID)
}

// GetByUser is method for getting entities.XUser by user.User
func (s *XUserService) GetByID(c context.Context, i string) (xu *entities.XUser, err error) {
	xu = &entities.XUser{ID: i}
	err = s.urep.Get(c, xu)
	return
}

func (s *XUserService) GetByName(c context.Context, name string) (*entities.XUser, error) {
	return s.urep.GetByName(c, name)
}

func (s *XUserService) IsUsedName(c context.Context, name string) (bool, error) {
	return s.urep.IsUsedName(c, name)
}

func (s *XUserService) GetLanguages(c context.Context, i string) (l []*entities.XUserLanguage, err error) {
	xu, err := s.GetByID(c, i)
	if err != nil {
		return nil, errors.Wrap(err, "GetByID failed")
	}

	l, err = s.ulrep.GetByXUser(c, xu)
	if err != nil {
		return nil, errors.Wrap(err, "GetByXUser failed")
	}
	return
}

func (s *XUserService) MigrateUniqueIndex(c context.Context) (err error) {
	return s.urep.MigrateUniqueIndex(c)
}

func validate(params entities.XUserProfileParams) (*validatekit.Validate, error) {
	return domain.ValidateXUserProfileParams(params)
}
