package web

import (
	"local/apikit/exchi"
	"local/xpo/domain/project"
	"local/xpo/entities"
	"net/http"

	"golang.org/x/net/context"
)

func GetUserProjects(w http.ResponseWriter, r *http.Request) error {
	c := Context(r)
	p := exchi.URLParams(r)
	i := p.Get("ownerId")
	return NewResponder(w, r).RenderObjectOrError(Services.Project().SearchByOwnerID(c, i))
}

func PostProject(c context.Context, w http.ResponseWriter, r *http.Request, xu *entities.XUser) error {
	p := project.Params{}
	if err := parseJSONBody(r, &p); err != nil {
		return err
	}

	return NewResponder(w, r).RenderObjectOrError(Services.Project().Create(c, xu, &p))
}
