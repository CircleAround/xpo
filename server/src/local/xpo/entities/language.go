package entities

// Language struct
type Language struct {
	Name          string `datastore:"-" goon:"id" json:"name"`
	ReportCount int64  `json:"reportCount"`
}
