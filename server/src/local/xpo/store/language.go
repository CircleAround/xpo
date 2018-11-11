package store

import (
	"context"
	"local/gaekit"
	"local/xpo/entities"

	"google.golang.org/appengine/datastore"
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

func (r *LanguageRepository) GetAll(c context.Context) (ls []*entities.Language, err error) {
	q := datastore.NewQuery("Language").Order("-ReportCount")
	ls = make([]*entities.Language, 0, 10)
	_, err = r.Goon(c).GetAll(q, &ls)
	return
}
