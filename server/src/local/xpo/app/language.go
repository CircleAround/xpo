package app

import (
	"local/xpo/entities"
	"local/xpo/store"

	"golang.org/x/net/context"
)

// Language is Service for Language
type LanguageService struct {
	lrep *store.LanguageRepository
}

// NewLanguageService is function for construction
func NewLanguageService() *LanguageService {
	s := new(LanguageService)
	s.lrep = store.NewLanguageRepository()
	return s
}

func (s *LanguageService) GetAll(c context.Context) ([]*entities.Language, error) {
	return s.lrep.GetAll(c)
}
