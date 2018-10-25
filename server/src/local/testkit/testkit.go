package testkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strconv"
	"testing"

	"context"

	"github.com/favclip/testerator"
	"google.golang.org/appengine/aetest"
)

// BootstrapTest is a method for spinup a first testerator
func BootstrapTest(m *testing.M) {
	_, _, err := testerator.SpinUp()

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	status := m.Run()

	err = testerator.SpinDown()
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	os.Exit(status)
}

// StartTest is a methid for spinup new testerator. with panic handling
func StartTest(t *testing.T) (aetest.Instance, context.Context, func()) {
	i, c, err := testerator.SpinUp()
	if err != nil {
		t.Fatal(err)
	}

	done := func() {
		testerator.SpinDown()

		if err := recover(); err != nil {
			t.Logf("[ERROR] %s\n", err)
			for depth := 0; ; depth++ {
				_, file, line, ok := runtime.Caller(depth)
				if !ok {
					break
				}
				t.Logf("======> %d: %v:%d", depth, file, line)
			}
		}
	}

	return i, c, done
}

func UnmarshalJSONBody(rr *httptest.ResponseRecorder, data interface{}) error {
	return json.Unmarshal(([]byte)(rr.Body.String()), data)
}

func NewRequestWithBody(i aetest.Instance, method string, path string, data interface{}) (*http.Request, error) {
	res, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	req, err := i.NewRequest(method, path, bytes.NewBuffer(res))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Length", strconv.Itoa(len(res)))
	return req, nil
}

func NewPostRequest(i aetest.Instance, path string, data interface{}) (*http.Request, error) {
	return NewRequestWithBody(i, "POST", path, data)
}
