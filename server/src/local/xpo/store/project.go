package store

import (
	"context"
	"local/gaekit"
	"local/xpo/domain/project"
	"local/xpo/entities"

	"github.com/pkg/errors"
	"google.golang.org/appengine/datastore"
)

type ProjectSearchParams struct {
	OwnerID string
}

type ProjectRepository struct {
	gaekit.DatastoreAccessObject
	IdentityNameUniqueIndexRepository
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

func (s *ProjectRepository) Build(c context.Context, xu *entities.XUser) *project.Root {
	pr := &entities.Project{
		OwnerID:  xu.ID,
		OwnerKey: s.KeyOf(c, xu),
	}
	return project.NewWithEntity(pr)
}

func (r *ProjectRepository) Find(c context.Context, xu *entities.XUser, i int64) (*project.Root, error) {
	pr := &entities.Project{
		OwnerID:  xu.ID,
		OwnerKey: r.KeyOf(c, xu),
		ID:       i,
	}

	err := r.Get(c, pr)
	if err != nil {
		return nil, errors.Wrap(err, "Get failed")
	}

	return project.NewWithEntity(pr), nil
}

func (s *ProjectRepository) Store(c context.Context, pr *project.Root) error {
	if pr.HasNameChanged() {
		err := s.ChangeMustTr(c, pr.NameChangedEvent.From, pr.NameChangedEvent.To)
		if err != nil {
			return errors.Wrap(err, "updateUniqueIndex failed")
		}
	}
	return s.Put(c, pr.Project)
}
