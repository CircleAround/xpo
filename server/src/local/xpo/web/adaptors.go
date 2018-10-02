package web

import (
	"local/apikit"
	"local/gaekit"
	"local/xpo/app"
	"net/http"
	"reflect"

	"github.com/iancoleman/strcase"
	"github.com/mjibson/goon"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
	validator "gopkg.in/go-playground/validator.v9"
)

func Auth(next func(context.Context, http.ResponseWriter, *http.Request, *app.XUser) error) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		c := Context(r)
		xu := xUserOrResponse(w, r)
		if xu == nil {
			log.Warningf(c, "xu==nil. response 401")
			return nil
		}

		return next(c, w, r, xu)
	}
}

func xUserOrResponse(w http.ResponseWriter, r *http.Request) *app.XUser {
	c := Context(r)
	u := user.Current(c)
	g := goon.NewGoon(r)

	xu := &app.XUser{ID: u.ID}
	if err := g.Get(xu); err != nil {
		log.Warningf(c, "Oops! has not user!")
		responseUnauthorized(w, r)
		return nil
	}
	return xu
}

func Catch(handler func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		safeFilter(w, r, handler(w, r))
	}
}

func safeFilter(w http.ResponseWriter, r *http.Request, err error) {
	c := Context(r)

	if err != nil {
		log.Infof(c, "Handle Error: %v", err)

		switch err.(type) {
		case *gaekit.ValueNotUniqueError:
			apikit.RespondFailure(w, err, http.StatusUnprocessableEntity)
			return

		case *gaekit.DuplicatedObjectError:
			apikit.RespondFailure(w, err, http.StatusUnprocessableEntity)
			return

		case *validator.InvalidValidationError:
			apikit.RespondFailure(w, err, http.StatusUnprocessableEntity)
			return

		case *apikit.InvalidParameterError:
			apikit.RespondFailure(w, err, http.StatusUnprocessableEntity)
			return

		case validator.ValidationErrors:
			ve := apikit.NewValidationError()
			for _, err := range err.(validator.ValidationErrors) {
				ve.PushOne(strcase.ToSnake(err.Field()), err.Tag())
			}
			apikit.RespondFailure(w, ve, http.StatusUnprocessableEntity)
			return

		case *apikit.IllegalAccessError:
			apikit.RespondFailure(w, err, http.StatusForbidden)
			return

		default:
			log.Warningf(c, "err: %v, %v\n", err.Error(), reflect.TypeOf(err))
			apikit.RespondFailure(w, err, http.StatusInternalServerError)
			return
		}
	}
}
