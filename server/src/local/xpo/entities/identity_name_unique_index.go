package entities

// IdentityNameUniqueIndex is unique index of XUser and Project and Organization's Name
type IdentityNameUniqueIndex struct {
	Value string `datastore:"-" goon:"id"`
}

func (i *IdentityNameUniqueIndex) Property() string {
	return "name"
}
