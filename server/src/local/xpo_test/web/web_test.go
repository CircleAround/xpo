package web_test

import (
	"fmt"
	"io"
	"local/testkit"
	"local/xpo/entities"
	"local/xpo/web"
	"local/xpo_test"

	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"testing"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/user"
)

const origin = "http://example.com"

func TestWebXUserScenario(t *testing.T) {
	os.Setenv("ALLOW_ORIGIN", origin)

	i, _, done := testkit.StartTest(t)
	defer done()

	f := xpo.NewTestFactory()
	xu := f.BuildXUser()

	var u user.User
	u.Email = xu.Email
	u.ID = xu.ID

	{
		t.Log("GET /user/me")
		{
			t.Log("no headers")
			req, err := i.NewRequest("GET", "/users/me", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := ServeHTTP(req)
			if rr.Code != 403 {
				t.Errorf("It should return 403")
			}
		}

		{
			t.Log("without X-Requested-With")

			req, err := i.NewRequest("GET", "/users/me", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Origin", origin)

			rr := ServeHTTP(req)

			if rr.Code != 403 {
				t.Errorf("It should return 403")
			}
		}

		{
			t.Log("without Origin")

			req, _ := i.NewRequest("GET", "/users/me", nil)
			req.Header.Set("X-Requested-With", "XmlHttpRequest")

			rr := ServeHTTP(req)

			if rr.Code != 403 {
				t.Errorf("It should return 403")
			}
		}

		{
			t.Log("Google not loggedin")

			req, err := XHGet(i, "/users/me")
			if err != nil {
				t.Fatal(err)
			}

			rr := ServeHTTP(req)

			if rr.Code != 401 {
				t.Errorf("It should return 401")
			}
		}

		{
			t.Log("Google loggedin but XUser not found")

			req, err := XHGet(i, "/users/me")
			if err != nil {
				t.Fatal(err)
			}
			aetest.Login(&u, req)

			rr := ServeHTTP(req)

			if rr.Code != 200 {
				t.Errorf("It should return 200: %v", rr.Code)
			}

			if "application/json" != rr.HeaderMap.Get("Content-Type") {
				t.Errorf("It should have Content-Type application/json : %v", rr.HeaderMap.Get("Content-Type"))
			}

			if rr.Body.String() != `"BE_SIGN_UP"` {
				t.Errorf("It should have response body : %v", rr.Body.String())
			}
		}

		{
			t.Log("Create Success")

			data := entities.XUserProfileParams{Name: xu.Name, Nickname: xu.Nickname}
			req, err := XHPost(i, "/users/me", data)
			if err != nil {
				t.Fatal(err)
			}

			aetest.Login(&u, req)

			rr := ServeHTTP(req)

			if rr.Code != 200 {
				t.Errorf("It should return 200: %v", rr.Code)
			}

			if "application/json" != rr.HeaderMap.Get("Content-Type") {
				t.Errorf("It should have Content-Type application/json : %v", rr.HeaderMap.Get("Content-Type"))
			}
		}

		{
			t.Log("Get Success")

			req, err := XHGet(i, "/users/me")
			if err != nil {
				t.Fatal(err)
			}
			aetest.Login(&u, req)

			rr := ServeHTTP(req)

			if rr.Code != 200 {
				t.Errorf("It should return 200: %v", rr.Code)
			}

			if "application/json" != rr.HeaderMap.Get("Content-Type") {
				t.Errorf("It should have Content-Type application/json : %v", rr.HeaderMap.Get("Content-Type"))
			}

			data := web.XUserResponse{}
			if err := testkit.UnmarshalJSONBody(rr, &data); err != nil {
				t.Errorf("Couldnot parse as JSON :%v, %v", rr.Body.String(), err)
			}

			if data.ID != u.ID {
				t.Errorf("It should have specified value :%v", data.ID)
			}
			if data.Name != xu.Name {
				t.Errorf("It should have specified value :%v", data.Name)
			}
			if data.Nickname != xu.Nickname {
				t.Errorf("It should have specified value :%v", data.Nickname)
			}
			if data.Email != u.Email {
				t.Errorf("It should have specified value :%v", data.Email)
			}
			if !strings.HasPrefix(data.LogoutURL, "http:///_ah/logout?continue=http%3A//") {
				t.Errorf("It should have specified value :%v", data.LogoutURL)
			}
		}
	}

	{
		t.Log("GET /users/:id")
		{
			t.Log("Get Success")
			p := fmt.Sprintf("/users/%v", xu.ID)
			t.Logf("path: %v", p)
			req, err := XHGet(i, p)
			if err != nil {
				t.Fatal(err)
			}

			rr := ServeHTTP(req)

			if rr.Code != 200 {
				t.Errorf("It should return 200: %v", rr.Code)
			}

			if "application/json" != rr.HeaderMap.Get("Content-Type") {
				t.Errorf("It should have Content-Type application/json : %v", rr.HeaderMap.Get("Content-Type"))
			}
		}
	}
}

// Utils

func XHRequest(i aetest.Instance, method string, path string, body io.Reader) (*http.Request, error) {
	req, err := i.NewRequest(method, path, body)
	if err != nil {
		return req, err
	}

	setHeaders(req)
	return req, nil
}

func setHeaders(req *http.Request) {
	req.Header.Set("Origin", origin)
	req.Header.Set("X-Requested-With", "XmlHttpRequest")
}

func XHGet(i aetest.Instance, path string) (*http.Request, error) {
	return XHRequest(i, "GET", path, nil)
}

func XHPost(i aetest.Instance, path string, data interface{}) (*http.Request, error) {
	req, err := testkit.NewPostRequest(i, path, data)
	if err != nil {
		return req, err
	}

	setHeaders(req)
	return req, nil
}

func ServeHTTP(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	web.Router().ServeHTTP(rr, req)
	return rr
}
