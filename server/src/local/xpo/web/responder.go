package web

import (
	"local/apikit"
	"local/gaekit"
	"local/xpo/entities"
	"net/http"
	"os"
)

type XUserResponse struct {
	*entities.XUser
	LoginURL  string `json:"loginUrl"`
	LogoutURL string `json:"logoutUrl"`
	Version   string `json:"version"`
}

type Responder struct {
	w  http.ResponseWriter
	r  *http.Request
	jr apikit.JSONRenderer
}

func NewResponder(w http.ResponseWriter, r *http.Request) *Responder {
	return &Responder{w: w, r: r, jr: apikit.NewJSONRenderer(w)}
}

func (r *Responder) RenderMeOrError(xu *entities.XUser, err error) error {
	res := XUserResponse{
		XUser:     xu,
		LoginURL:  gaekit.LoginFullURL(r.r, loggedInPath),
		LogoutURL: gaekit.LogoutFullURL(r.r),
		Version:   os.Getenv("CURRENT_VERSION_ID"),
	}
	return r.jr.RenderOrError(res, err)
}

func (r *Responder) RenderObjectOrError(obj interface{}, err error) error {
	return r.jr.RenderOrError(obj, err)
}
