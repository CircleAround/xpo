package project

import (
	"local/xpo/entities"
	"time"

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

type NameChangedEvent struct {
	From string
	To   string
}

type Root struct {
	*entities.Project
	NameChangedEvent *NameChangedEvent
	Error            error
}

func NewWithEntity(p *entities.Project) *Root {
	return &Root{
		p,
		nil,
		nil,
	}
}

func (r *Root) UpdateAttributes(p Params) error {

	if r.Name != p.Name {
		r.NameChangedEvent = &NameChangedEvent{From: r.Name, To: p.Name}
		r.Name = p.Name
	}

	if !p.Description.IsZero() {
		r.Description = p.Description.ValueOrZero()
	}

	if !p.RepositoryURL.IsZero() {
		r.RepositoryURL = p.RepositoryURL.ValueOrZero()
	}

	now := time.Now()
	r.UpdatedAt = now
	if r.CreatedAt.IsZero() {
		r.CreatedAt = now
	}

	_, err := Validate(r.Project)
	if err == nil {
		return nil
	}

	return errors.Wrap(err, "validate failed")
}

func (r *Root) HasNameChanged() bool {
	return r.NameChangedEvent != nil
}
