package web

import (
	"local/apikit"
	"local/gaekit"
	"local/xpo/app"
	"net/http"
	"reflect"

	"google.golang.org/appengine/datastore"

	"github.com/iancoleman/strcase"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
	validator "gopkg.in/go-playground/validator.v9"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func Auth(next func(context.Context, http.ResponseWriter, *http.Request, *app.XUser) error) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		c := Context(r)
		u := user.Current(c)

		xu, err := Services.XUser().GetByUser(c, *u)
		if err == datastore.ErrNoSuchEntity {
			apikit.RespondJSON(w, "BE_SIGN_UP")
			return nil
		}

		return next(c, w, r, xu)
	}
}

func Handler(next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		safeFilter(w, r, next(w, r))
	}
}

func safeFilter(w http.ResponseWriter, r *http.Request, err error) {
	if err == apikit.UnauthorizedError {
		responseUnauthorized(w, r)
		return
	}

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
