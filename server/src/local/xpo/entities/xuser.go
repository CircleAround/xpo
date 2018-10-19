package entities

// XUser struct
type XUser struct {
	ID          string `datastore:"-" goon:"id" json:"id"`
	Name        string `json:"name" validate:"required,min=3,max=15,username_format"`
	Email       string `json:"email" validate:"required"`
	Nickname    string `json:"nickname" validate:"required,min=3,max=24,usernickname_format"`
	ReportCount int64  `json:"reportCount"`
}
