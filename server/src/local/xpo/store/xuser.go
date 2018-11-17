package store

import (
	"context"
	"errors"
	"local/gaekit"
	"local/xpo/entities"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type XUserRepository struct {
	gaekit.DatastoreAccessObject
}

func NewXUserRepository() *XUserRepository {
	return new(XUserRepository)
}

// _XUserNameUniqueIndex is unique index of XUser's Name
type _XUserNameUniqueIndex struct {
	value string `datastore:"-" goon:"id"`
}

func (i *_XUserNameUniqueIndex) Property() string {
	return "name"
}

func (r *XUserRepository) Create(c context.Context, xu *entities.XUser) (err error) {
	if err = r.Get(c, xu); err == nil {
		log.Infof(c, "%v found!. duplicated.", xu)
		return &gaekit.DuplicatedObjectError{Type: "duplicatedObject"}
	}

	if err != datastore.ErrNoSuchEntity {
		log.Infof(c, "%v error.", err)
		return err
	}

	return r.RunInXGTransaction(c, func(ctx context.Context) error {
		// for idempotent. if already create success, return process.
		if err = r.Get(ctx, xu); err == nil {
			return nil
		}

		if err != datastore.ErrNoSuchEntity {
			log.Infof(ctx, "%v error.", err)
			return err
		}

		log.Infof(ctx, "%v not found.", xu)

		i := &entities.IdentityNameUniqueIndex{Value: xu.Name}
		err = r.CreateUnique(ctx, i)
		if err != nil {
			return err
		}

		log.Infof(ctx, "%v not found. create new one.", xu)
		return r.Put(ctx, xu)
	})
}

func (s *XUserRepository) Update(c context.Context, xu *entities.XUser, params entities.XUserProfileParams) (err error) {
	return s.RunInXGTransaction(c, func(ctx context.Context) error {
		if xu.Name == params.Name {
			if xu.Nickname == params.Nickname {
				// Maybe already succeed. return nil for idempotent.
				return nil
			}
		} else {
			err = s.updateUniqueIndex(ctx, *xu, params)
			if err != nil {
				return err
			}

			xu.Name = params.Name
		}

		xu.Nickname = params.Nickname
		return s.Put(ctx, xu)
	})

}

func (s *XUserRepository) GetByName(c context.Context, name string) (*entities.XUser, error) {
	q := datastore.NewQuery("XUser").Filter("Name=", name).Limit(1)
	var xus []entities.XUser
	_, err := s.Goon(c).GetAll(q, &xus)
	if err != nil {
		return nil, err
	}
	if len(xus) == 0 {
		return nil, datastore.ErrNoSuchEntity
	}
	return &xus[0], nil
}

func (s *XUserRepository) updateUniqueIndex(c context.Context, xu entities.XUser, params entities.XUserProfileParams) error {
	i := &entities.IdentityNameUniqueIndex{Value: xu.Name}
	ni := &entities.IdentityNameUniqueIndex{Value: params.Name}
	return s.ChangeUniqueValueMustTr(c, i, ni)
}

// IsUsedName is method for checking UserName already taken.
func (s *XUserRepository) IsUsedName(c context.Context, name string) (bool, error) {
	i := entities.IdentityNameUniqueIndex{Value: name}
	return s.Exists(c, &i)
}

func (s *XUserRepository) MigrateUniqueIndex(c context.Context) (err error) {
	q := datastore.NewQuery("_XUserNameUniqueIndex")
	var uis []entities.XUser
	keys, err := q.GetAll(c, &uis)
	if err != nil {
		return err
	}

	for _, key := range keys {
		log.Infof(c, "key: %v", key)
		value := key.StringID()
		log.Infof(c, "string name: %v", value)

		if value == "" {
			log.Errorf(c, "値が取れない: %v", key)
			return errors.New("値が取れない")
		}

		ii := &entities.IdentityNameUniqueIndex{Value: value}
		err = s.Get(c, ii)
		if err == nil {
			continue
		}

		if err != datastore.ErrNoSuchEntity {
			return err
		}

		err = s.Put(c, ii)
		if err != nil {
			return err
		}
	}
	return nil
}
