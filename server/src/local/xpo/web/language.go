package web

import (
	"local/xpo/domain"
	"net/http"
)

func GetAllLanguageNames(w http.ResponseWriter, r *http.Request) error {
	return NewResponder(w, r).RenderObjectOrError(domain.Languages.Src(), nil)
}

func GetAllLanguages(w http.ResponseWriter, r *http.Request) error {
	c := Context(r)
	return NewResponder(w, r).RenderObjectOrError(Services.Language().GetAll(c))
}
