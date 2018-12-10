package project

import (
	"local/xpo/entities"

	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

type Params struct {
	Name          string      `json:"name" validate:"required"`
	Description   null.String `json:"description"`
	RepositoryURL null.String `json:"repositoryUrl"`
}

type UpdatingParams struct {
	Params
	ID int64 `json:"id" validate:"required"`
}

type NameDurtyEvent struct {
	From string
	To   string
}

type Root struct {
	*entities.Project
	NameDurtyEvent *NameDurtyEvent
	Error          error
}

func NewWithEntity(p *entities.Project) *Root {
	return &Root{
		p,
		nil,
		nil,
	}
}

func (r *Root) SetAttributes(p Params) {

	if r.Name != p.Name {
		r.NameDurtyEvent = &NameDurtyEvent{From: r.Name, To: p.Name}
		r.Name = p.Name
	}

	if !p.Description.IsZero() {
		r.Description = p.Description.ValueOrZero()
	}

	if !p.RepositoryURL.IsZero() {
		r.RepositoryURL = p.RepositoryURL.ValueOrZero()
	}
}

func (r *Root) Validate() bool {
	_, err := Validate(r.Project)
	if err == nil {
		return true
	}

	r.Error = errors.Wrap(err, "validate failed")
	return false
}
