package app

import (
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
	_, err := project.Validate(p)
	if err != nil {
		return nil, errors.Wrap(err, "params validate failed")
	}

	prj := s.ps.Build(c, xu)
	err = prj.UpdateAttributes(*p)
	if err != nil {
		return nil, errors.Wrap(err, "UpdateAttributes failed")
	}

	err = s.ps.RunInXGTransaction(c, func(c context.Context) error {
		err = s.ps.Store(c, prj)
		if err != nil {
			return errors.Wrap(err, "Store failed")
		}
		return nil
	})
	return prj.Project, err
}

func (s *ProjectService) Update(c context.Context, xu *entities.XUser, p *project.UpdatingParams) (*entities.Project, error) {
	_, err := project.Validate(p)
	if err != nil {
		return nil, errors.Wrap(err, "params validate failed")
	}

	var prj project.Root
	err = s.ps.RunInXGTransaction(c, func(c context.Context) error {
		prj, err := s.ps.Find(c, xu, p.ID)
		if err != nil {
			return errors.Wrap(err, "Retrive failed")
		}

		err = prj.UpdateAttributes(p.Params)
		if err != nil {
			return errors.Wrap(err, "UpdateAttributes failed")
		}

		err = s.ps.Store(c, prj)
		if err != nil {
			return errors.Wrap(err, "Store failed")
		}
		return nil
	})

	return prj.Project, err
}
