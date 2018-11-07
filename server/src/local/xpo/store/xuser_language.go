package store

import (
	"context"
	"local/gaekit"
	"local/xpo/entities"

	"google.golang.org/appengine/datastore"
)

type XUserLanguageRepository struct {
	gaekit.DatastoreAccessObject
}

func NewXUserLanguageRepository() *XUserLanguageRepository {
	return &XUserLanguageRepository{}
}

func (r *XUserLanguageRepository) GetByXUserAndName(c context.Context, xu *entities.XUser, n string) (ul *entities.XUserLanguage, err error) {
	ul = &entities.XUserLanguage{Name: n, XUserKey: r.KeyOf(c, xu)}
	err = r.Get(c, ul)
	return
}

func (r *XUserLanguageRepository) GetByXUser(c context.Context, xu *entities.XUser) (uls []*entities.XUserLanguage, err error) {
	q := datastore.NewQuery("XUserLanguage").Ancestor(r.KeyOf(c, xu)).Order("Name")
	uls = make([]*entities.XUserLanguage, 0)
	_, err = r.Goon(c).GetAll(q, &uls)
	return
}
