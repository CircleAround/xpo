package app

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"local/xpo/domain/project"
	"local/xpo/entities"
	"local/xpo/store"
)

// ProjectService is Service for XUser
type ProjectService struct {
	ps *store.ProjectRepository
}

// NewProjectService is function for construction
func NewProjectService() *ProjectService {
	return &ProjectService{ps: store.NewProjectRepository()}
}

func (s *ProjectService) SearchByOwnerID(c context.Context, i string) ([]*entities.Project, error) {
	return s.ps.Search(c, store.ProjectSearchParams{OwnerID: i}, 0)
}

func (s *ProjectService) Create(c context.Context, xu *entities.XUser, p *project.Params) (*entities.Project, error) {
	now := time.Now()

	pr := &entities.Project{
		OwnerID:       xu.ID,
		OwnerKey:      s.ps.KeyOf(c, xu),
		Description:   p.Description.ValueOrZero(),
		RepositoryURL: p.RepositoryURL.ValueOrZero(),
		Name:          p.Name,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	_, err := project.Validate(pr)
	if err != nil {
		return nil, errors.Wrap(err, "validate failed")
	}

	err = s.ps.Create(c, pr)
	if err != nil {
		return nil, errors.Wrap(err, "Put failed")
	}
	return pr, nil
}

func (s *ProjectService) Update(c context.Context, xu *entities.XUser, p *project.UpdatingParams) (*entities.Project, error) {
	_, err := project.Validate(p)
	if err != nil {
		return nil, errors.Wrap(err, "params validate failed")
	}

	prj, err := s.ps.Retrive(c, xu, p.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Retrive failed")
	}

	prj.SetAttributes(p.Params)
	now := time.Now()
	prj.UpdatedAt = now

	if !prj.Validate() {
		return nil, errors.Wrap(prj.Error, "validate failed")
	}

	err = s.ps.Update(c, prj)
	if err != nil {
		return nil, errors.Wrap(err, "Put failed")
	}
	return prj.Project, nil
}
