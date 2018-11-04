package entities

import (
	"time"

	"google.golang.org/appengine/datastore"
)

// Report struct
type Report struct {
	ID             int64          `json:"id" datastore:"-" goon:"id"`
	AuthorKey      *datastore.Key `json:"-" datastore:"-" goon:"parent" validate:"required"`
	AuthorID       string         `json:"authorId" validate:"required"`
	Author         string         `json:"author" validate:"required"`
	AuthorNickname string         `json:"authorNickname" validate:"required"`
	Content        string         `json:"content" validate:"required,max=20000" datastore:"Content,noindex"`
	ContentType    string         `json:"contentType" validate:"required"`
	Languages      []string       `json:"languages" validate:"languages"`
	ReportedAt     time.Time      `json:"reportedAt"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}
