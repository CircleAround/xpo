package xpo_test

import (
	"local/apikit"
	"local/gaekit"
	"local/testkit"
	"local/xpo/app"
	"reflect"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/appengine/user"
)

func TestXUserScenario(t *testing.T) {
	_, c, done := testkit.StartTest(t)
	defer done()

	f := NewTestFactory()
	s := app.NewXUserService()

	d := f.BuildXUser()
	var u user.User
	u.Email = d.Email
	u.ID = d.ID

	{
		t.Log("Scenario")

		xu, err := s.Create(c, u, app.XUserProfileParams{Name: d.Name, Nickname: d.Nickname})
		if err != nil {
			t.Fatal(err)
		}

		{
			t.Logf("Standard")
			ret := app.XUser{ID: xu.ID}
			if err = s.Get(c, &ret); err != nil {
				t.Fatal(err)
			}

			checkXUser(t, c, s, f, u, ret, d)
		}
	}

	{
		t.Logf("Update")
		ud := f.BuildXUser()

		ret, err := s.Update(c, u, app.XUserProfileParams{
			Name:     ud.Name,
			Nickname: ud.Nickname,
		})
		if err != nil {
			t.Fatal(err)
		}
		if ret == nil {
			t.Fatal("updated ret is nil")
		}

		ud.Email = u.Email // Email not changed
		checkXUser(t, c, s, f, u, *ret, ud)

		used, err := s.IsUsedName(c, d.Name)
		if err != nil {
			t.Fatal(err)
		}
		if used {
			t.Fatalf("It should release before name!: %v", d.Name)
		}
	}
}

func TestValidation(t *testing.T) {
	_, c, done := testkit.StartTest(t)
	defer done()

	f := NewTestFactory()
	s := app.NewXUserService()
	t.Logf("Validation")

	d := f.BuildXUser()

	var u user.User
	u.Email = d.Email
	u.ID = d.ID

	{
		t.Logf("Name")
		{
			_, err := s.Create(c, u, app.XUserProfileParams{Nickname: d.Nickname})
			apikit.ShouldHaveRequiredError(t, err, "Name")
		}

		{
			_, err := s.Create(c, u, app.XUserProfileParams{Name: "", Nickname: d.Nickname})
			apikit.ShouldHaveRequiredError(t, err, "Name")
		}

		{
			_, err := s.Create(c, u, app.XUserProfileParams{Name: "a_&", Nickname: d.Nickname})
			apikit.ShouldHaveInvalidFormatError(t, err, "Name", "username_format")
		}

		{
			_, err := s.Create(c, u, app.XUserProfileParams{Name: "admin", Nickname: d.Nickname})
			if reflect.TypeOf(err) != reflect.TypeOf(&apikit.InvalidParameterError{}) {
				t.Fatalf("It should be apikit.InvalidParameterError: %v, %v", reflect.TypeOf(err), err)
			}
		}
	}

	{
		t.Logf("Nickname")
		{
			_, err := s.Create(c, u, app.XUserProfileParams{Name: d.Name})
			apikit.ShouldHaveRequiredError(t, err, "Nickname")
		}

		{
			_, err := s.Create(c, u, app.XUserProfileParams{Name: d.Name, Nickname: ""})
			apikit.ShouldHaveRequiredError(t, err, "Nickname")
		}

		{
			_, err := s.Create(c, u, app.XUserProfileParams{Name: d.Name, Nickname: "<nynickname"})
			apikit.ShouldHaveInvalidFormatError(t, err, "Nickname", "usernickname_format")
		}

		{
			_, err := s.Create(c, u, app.XUserProfileParams{Name: d.Name, Nickname: "reports"}) // reports is blocked
			if reflect.TypeOf(err) != reflect.TypeOf(&apikit.InvalidParameterError{}) {
				t.Fatalf("It should be apikit.InvalidParameterError: %v, %v", reflect.TypeOf(err), err)
			}
		}

	}
}

func checkXUser(t *testing.T, c context.Context, s *app.XUserService, f *TestFactory, u user.User, ret app.XUser, d app.XUser) {
	t.Logf("Update!")

	if ret.Email != d.Email {
		t.Fatalf("It should get saved email!: %v", ret.Email)
	}
	if ret.Name != d.Name {
		t.Fatalf("It should get saved Name!: %v", ret.Name)
	}
	if ret.Nickname != d.Nickname {
		t.Fatalf("It should get saved Nickname!: %v", ret.Nickname)
	}
	t.Logf("Update!!")

	used, err := s.IsUsedName(c, d.Name)
	if err != nil {
		t.Fatal(err)
	}
	if !used {
		t.Fatalf("It should Name is Used!: %v", d.Name)
	}
	t.Logf("Update!!!")

	{
		t.Logf("Duplicaed")
		_, err := s.Create(c, u, app.XUserProfileParams{Name: d.Name, Nickname: d.Nickname})
		if err == nil {
			t.Fatal("It should error on creating duplicated user")
		}

		if reflect.TypeOf(err) != reflect.TypeOf(&gaekit.DuplicatedObjectError{}) {
			t.Fatalf("It should be DuplicatedObjectError: %v", reflect.TypeOf(err))
		}
	}

	{
		t.Logf("Uniqueness")

		d2 := f.BuildXUser()
		var u user.User
		u.Email = d2.Email
		u.ID = d2.ID

		// duplicated name
		_, err := s.Create(c, u, app.XUserProfileParams{Name: d.Name, Nickname: d.Nickname})
		if err == nil {
			t.Fatal("It should error on creating duplicated name user")
		}

		if reflect.TypeOf(err) != reflect.TypeOf(&gaekit.ValueNotUniqueError{}) {
			t.Fatalf("It should be ValueNotUniqueError: %v", reflect.TypeOf(err))
		}
	}
}
