package app

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"local/xpo/domain"
	"local/xpo/entities"
	"local/xpo/store"
)

type ProjectCreationParams struct {
	Description   string `json:"description"`
	Name          string `json:"name" validate:"required"`
	RepositoryURL string `json:"repositoryUrl"`
}

type ProjectUpdatingParams struct {
	ProjectCreationParams
	ID int64 `json:"id" validate:"required"`
}

// ProjectService is Service for XUser
type ProjectService struct {
	ps *store.ProjectRepository
}

// NewProjectService is function for construction
func NewProjectService() *ProjectService {
	return &ProjectService{}
}

func (s *ProjectService) SearchByOwnerID(c context.Context, i string) ([]*entities.Project, error) {
	return s.ps.Search(c, store.ProjectSearchParams{OwnerID: i}, 0)
}

func (s *ProjectService) Create(c context.Context, xu *entities.XUser, p *ProjectCreationParams) (*entities.Project, error) {
	now := time.Now()
	pr := &entities.Project{
		OwnerID:     xu.ID,
		OwnerKey:    s.ps.KeyOf(c, xu),
		Description: p.Description,
		Name:        p.Name,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err := domain.ValidateProject(pr)
	if err != nil {
		return nil, errors.Wrap(err, "validate failed")
	}

	err = s.ps.Put(c, pr)
	if err != nil {
		return nil, errors.Wrap(err, "Put failed")
	}
	return pr, nil
}

func (s *ProjectService) Update(c context.Context, xu *entities.XUser, p *ProjectUpdatingParams) (*entities.Project, error) {
	pr := &entities.Project{
		OwnerID:  xu.ID,
		OwnerKey: s.ps.KeyOf(c, xu),
		ID:       p.ID,
	}

	err := s.ps.Get(c, pr)
	if err != nil {
		return nil, errors.Wrap(err, "Get failed")
	}

	pr.Name = p.Name
	pr.RepositoryURL = p.RepositoryURL

	now := time.Now()
	pr.UpdatedAt = now

	_, err = domain.ValidateProject(pr)
	if err != nil {
		return nil, errors.Wrap(err, "validate failed")
	}

	err = s.ps.Put(c, pr)
	if err != nil {
		return nil, errors.Wrap(err, "Put failed")
	}
	return pr, nil
}
