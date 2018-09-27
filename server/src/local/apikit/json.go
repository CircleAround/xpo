package apikit

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// ParseJSONBody is parse from requested body formatted JSON
func ParseJSONBody(r *http.Request, jsonBody interface{}) error {
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	body := make([]byte, length)
	length, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		return err
	}

	err = json.Unmarshal(body[:length], jsonBody)
	if err != nil {
		return err
	}

	return nil
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
		return err
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
