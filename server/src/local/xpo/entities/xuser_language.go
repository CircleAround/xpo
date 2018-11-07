package entities

import "google.golang.org/appengine/datastore"

type XUserLanguage struct {
	Name        string         `datastore:"-" goon:"id" json:"name"`
	XUserKey    *datastore.Key `datastore:"-" goon:"parent" json:"xuserKey"`
	ReportCount int64          `json:"reportCount"`
}
