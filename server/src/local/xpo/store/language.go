package store

import (
	"context"
	"local/gaekit"
	"local/xpo/entities"
)

type LanguageRepository struct {
	gaekit.DatastoreAccessObject
}

func NewLanguageRepository() *LanguageRepository {
	return &LanguageRepository{}
}

func (r *LanguageRepository) GetByName(c context.Context, name string) (l *entities.Language, err error) {
	l = &entities.Language{Name: name}
	err = r.Get(c, l)
	return
}
