package entities

import (
	"time"

	"google.golang.org/appengine/datastore"
)

// Project struct
type Project struct {
	ID            int64          `json:"id" datastore:"-" goon:"id"`
	OwnerKey      *datastore.Key `json:"-" datastore:"ownerKey" validate:"required"`
	OwnerID       string         `json:"ownerId" validate:"required"`
	Name          string         `json:"name" validate:"required"`
	Description   string         `json:"description"`
	RepositoryURL string         `json:"repositoryUrl" validate:"url"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}
