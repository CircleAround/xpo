package apikit

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// ParseJSONBody is parse from requested body formatted JSON
func ParseJSONBody(r *http.Request, jsonBody interface{}) error {
	limitedReader := &io.LimitedReader{R: r.Body, N: 10000000}
	b, err := ioutil.ReadAll(limitedReader)
	defer r.Body.Close()
	if err != nil {
		return errors.Wrap(err, "ReadAll failed")
	}

	return json.Unmarshal(b, jsonBody)
}

func RespondJSON(w http.ResponseWriter, obj interface{}) error {
	w.WriteHeader(http.StatusOK)
	return writeJSON(w, obj)
}

func RespondOk(w http.ResponseWriter) error {
	return RespondJSON(w, NewSuccess())
}

func RespondFailure(w http.ResponseWriter, err interface{}, code int) error {
	return DoRespondFailure(w, NewFailure(err), code)
}

func DoRespondFailure(w http.ResponseWriter, failure Failure, code int) error {
	w.WriteHeader(code)
	return writeJSON(w, failure)
}

func writeJSON(w http.ResponseWriter, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(obj)

	if err != nil {
		return errors.Wrap(err, "Marshal failed")
	}

	_, err = w.Write(res)
	return err
}

type JSONRenderer struct {
	writer http.ResponseWriter
}

func NewJSONRenderer(w http.ResponseWriter) JSONRenderer {
	return JSONRenderer{writer: w}
}

func (r *JSONRenderer) NG(err interface{}, code int) error {
	return RespondFailure(r.writer, err, code)
}

func (r *JSONRenderer) OK(err interface{}, code int) error {
	return RespondOk(r.writer)
}

func (r *JSONRenderer) Render(obj interface{}) error {
	return RespondJSON(r.writer, obj)
}

func (r *JSONRenderer) RenderOrError(obj interface{}, err error) error {
	if err != nil {
		return err
	}
	return RespondJSON(r.writer, obj)
}
