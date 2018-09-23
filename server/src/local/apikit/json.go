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

func ResponseJSON(w http.ResponseWriter, obj interface{}) {
	w.WriteHeader(http.StatusOK)
	writeJSON(w, obj)
}

func ResponseOk(w http.ResponseWriter) {
	ResponseJSON(w, NewSuccess())
}

func ResponseFailure(w http.ResponseWriter, r *http.Request, err interface{}, code int) {
	DoResponseFailure(w, r, NewFailure(err), code)
}

func DoResponseFailure(w http.ResponseWriter, r *http.Request, failure Failure, code int) {
	w.WriteHeader(code)

	writeJSON(w, failure)
}

func writeJSON(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(obj)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(res)
}
