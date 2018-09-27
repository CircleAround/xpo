package xpo

import (
	"net/http"
	"os"
	"reflect"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/iancoleman/strcase"
	validator "gopkg.in/go-playground/validator.v9"

	"local/apikit"
	"local/gaekit"
)

func allowOrigin(w http.ResponseWriter, origin string) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}

func allowClient(w http.ResponseWriter) {
	allowOrigin(w, os.Getenv("ALLOW_ORIGIN"))
}

func safeFilter(w http.ResponseWriter, r *http.Request, err error) {
	c := appengine.NewContext(r)

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
