package web_test

import (
	"local/testkit"
	"local/xpo/app"
	"local/xpo/web"
	"net/http/httptest"
	"strings"

	"testing"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/user"
)

func TestWebScenario(t *testing.T) {
	i, _, done := testkit.StartTest(t)
	defer done()

	var u user.User
	u.Email = "test@example.com"
	u.ID = "1"

	{
		t.Log("GET /user/me")
		{
			t.Log("Google not loggedin")

			req, _ := i.NewRequest("GET", "/users/me", nil)

			rr := httptest.NewRecorder()
			web.Router().ServeHTTP(rr, req)

			if rr.Code != 401 {
				t.Errorf("It should return 401")
			}
		}

		{
			t.Log("Google loggedin but XUser not found")

			req, _ := i.NewRequest("GET", "/users/me", nil)
			aetest.Login(&u, req)

			rr := httptest.NewRecorder()
			web.Router().ServeHTTP(rr, req)

			if rr.Code != 200 {
				t.Errorf("It should return 200")
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

			data := app.XUserProfileParams{Name: "aaaa", Nickname: "てすと"}
			req, err := testkit.NewRequestWithBody(i, "POST", "/users/me", data)
			if err != nil {
				t.Fatal(err)
			}

			aetest.Login(&u, req)

			rr := httptest.NewRecorder()
			web.Router().ServeHTTP(rr, req)

			if rr.Code != 200 {
				t.Fatal("It should return 200")
			}

			if "application/json" != rr.HeaderMap.Get("Content-Type") {
				t.Errorf("It should have Content-Type application/json : %v", rr.HeaderMap.Get("Content-Type"))
			}
		}

		{
			t.Log("Get Success")

			req, _ := i.NewRequest("GET", "/users/me", nil)
			aetest.Login(&u, req)

			rr := httptest.NewRecorder()
			web.Router().ServeHTTP(rr, req)

			if rr.Code != 200 {
				t.Errorf("It should return 200")
			}

			if "application/json" != rr.HeaderMap.Get("Content-Type") {
				t.Errorf("It should have Content-Type application/json : %v", rr.HeaderMap.Get("Content-Type"))
			}

			data := web.XUserResponse{}
			if err := testkit.UnmarshalJSONBody(rr, &data); err != nil {
				t.Errorf("Couldnot parse as JSON :%v, %v", rr.Body.String(), err)
			}

			if data.ID != "1" {
				t.Errorf("It should have specified value :%v", data.ID)
			}
			if data.Name != "aaaa" {
				t.Errorf("It should have specified value :%v", data.Name)
			}
			if data.Email != "test@example.com" {
				t.Errorf("It should have specified value :%v", data.Email)
			}
			if !strings.HasPrefix(data.LogoutURL, "http:///_ah/logout?continue=http%3A//") {
				t.Errorf("It should have specified value :%v", data.LogoutURL)
			}
		}
	}

}
