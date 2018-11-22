package store

import (
	"context"
	"local/gaekit"
	"local/xpo/entities"

	"google.golang.org/appengine/datastore"
)

type ProjectSearchParams struct {
	OwnerID string
}

type ProjectRepository struct {
	gaekit.DatastoreAccessObject
}

func NewProjectRepository() *ProjectRepository {
	return new(ProjectRepository)
}

func (s *ProjectRepository) Search(c context.Context, p ProjectSearchParams, limit int) (ps []*entities.Project, err error) {

	ps = make([]*entities.Project, 0, 10)
	err = s.FilterSearch(c, "Project", &ps, func(q *datastore.Query) *datastore.Query {
		if limit != 0 {
			q = q.Limit(limit)
		}

		if p.OwnerID != "" {
			q.Filter("OwnerID=", p.OwnerID)
		}

		return q
	})

	return
}
