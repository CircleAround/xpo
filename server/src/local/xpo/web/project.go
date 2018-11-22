package web

import (
	"local/apikit/exchi"
	"local/xpo/app"
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
	p := app.ProjectCreationParams{}
	if err := parseJSONBody(r, &p); err != nil {
		return err
	}

	return NewResponder(w, r).RenderObjectOrError(Services.Project().Create(c, xu, &p))
}
