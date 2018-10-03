package web

import (
	"local/apikit"
	"local/gaekit"
	"local/xpo/app"
	"net/http"
)

var Services *app.ServiceFactory = app.Factory()

const loggedInPath = "/loggedin"

func parseJSONBody(r *http.Request, p interface{}) error {
	err := apikit.ParseJSONBody(r, &p)
	if err != nil {
		return apikit.NewInvalidParameterErrorWithMessage("json", err.Error())
	}
	return nil
}

func responseUnauthorized(w http.ResponseWriter, r *http.Request) {
	apikit.RespondFailure(w, gaekit.LoginFullURL(r, loggedInPath), http.StatusUnauthorized)
}
