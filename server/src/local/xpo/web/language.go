package web

import (
	"local/xpo/domain"
	"net/http"
)

func GetLanguages(w http.ResponseWriter, r *http.Request) error {
	return NewResponder(w, r).RenderObjectOrError(domain.Languages.Src(), nil)
}
