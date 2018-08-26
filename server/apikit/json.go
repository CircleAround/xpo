package apikit

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// ParseJSONBody is parse from requested body formatted JSON
func ParseJSONBody(r *http.Request) (map[string]interface{}, error) {
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		return nil, err
	}

	body := make([]byte, length)
	length, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		return nil, err
	}

	var jsonBody map[string]interface{}
	err = json.Unmarshal(body[:length], &jsonBody)
	if err != nil {
		return nil, err
	}

	return jsonBody, nil
}
