package store

import (
	"context"
	"local/gaekit"
	"local/xpo/entities"

	"github.com/pkg/errors"
)

type IdentityNameUniqueIndexRepository struct {
	gaekit.DatastoreAccessObject
}

func NewIdentityNameUniqueIndexRepository() *IdentityNameUniqueIndexRepository {
	return new(IdentityNameUniqueIndexRepository)
}

// IsUsedName is method for checking UserName already taken.
func (s *IdentityNameUniqueIndexRepository) IsUsedName(c context.Context, name string) (bool, error) {
	i := entities.IdentityNameUniqueIndex{Value: name}
	return s.Exists(c, &i)
}

func (s *IdentityNameUniqueIndexRepository) ChangeMustTr(c context.Context, from, to string) error {
	i := &entities.IdentityNameUniqueIndex{Value: from}
	ni := &entities.IdentityNameUniqueIndex{Value: to}
	return s.ChangeUniqueValueMustTr(c, i, ni)
}

func NewIdentityNamedEntityCreator(e interface{}, n string) *gaekit.EntityCreator {
	r := gaekit.NewEntityCreator()

	r.OnStartCreate = func(c context.Context) error {
		i := &entities.IdentityNameUniqueIndex{Value: n}
		err := r.CreateUnique(c, i)
		if err != nil {
			return errors.Wrap(err, "CreateUnique failed")
		}
		return nil
	}

	return r
}
