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

func (r *ProjectRepository) Retrive(c context.Context, xu *entities.XUser, i int64) (*project.Root, error) {
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

func (r *ProjectRepository) Create(c context.Context, pr *entities.Project) (err error) {
	return NewIdentityNamedEntityCreator(pr, pr.Name).Execute(c, pr)
}

func (s *ProjectRepository) Update(c context.Context, pr *project.Root) (err error) {
	return s.RunInXGTransaction(c, func(c context.Context) error {
		if pr.NameDurtyEvent != nil {
			err = s.ChangeMustTr(c, pr.NameDurtyEvent.From, pr.NameDurtyEvent.To)
			if err != nil {
				return errors.Wrap(err, "updateUniqueIndex failed")
			}
		}
		return s.Put(c, pr.Project)
	})

}
