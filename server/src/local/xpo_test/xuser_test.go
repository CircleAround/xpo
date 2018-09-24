package xpo_test

import (
	"local/apikit"
	"local/gaekit"
	"local/testkit"
	"local/xpo"
	"reflect"
	"testing"

	"google.golang.org/appengine/user"
)

func TestXUserScenario(t *testing.T) {
	_, c, done := testkit.StartTest(t)
	defer done()

	f := NewTestFactory(c)
	s := xpo.NewXUserService(c)

	d := f.BuildXUser()
	var u user.User
	u.Email = d.Email
	u.ID = d.ID

	{
		t.Log("Scenario")

		xu, err := s.Create(&u, &xpo.XUserProfileParams{Name: d.Name, Nickname: d.Nickname})
		if err != nil {
			t.Fatal(err)
		}

		{
			t.Logf("Standard")
			ret := &xpo.XUser{ID: xu.ID}
			if err = s.Get(ret); err != nil {
				t.Fatal(err)
			}

			checkXUser(t, s, f, u, *ret, d)
		}
	}

	{
		t.Logf("Update")
		ud := f.BuildXUser()

		ret, err := s.Update(&u, &xpo.XUserProfileParams{
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
		checkXUser(t, s, f, u, *ret, ud)

		used, err := s.IsUsedName(d.Name)
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

	f := NewTestFactory(c)
	s := xpo.NewXUserService(c)
	t.Logf("Validation")

	d := f.BuildXUser()

	var u user.User
	u.Email = d.Email
	u.ID = d.ID

	{
		t.Logf("Name")
		{
			_, err := s.Create(&u, &xpo.XUserProfileParams{Nickname: d.Nickname})
			apikit.ShouldHaveRequiredError(t, err, "Name")
		}

		{
			_, err := s.Create(&u, &xpo.XUserProfileParams{Name: "", Nickname: d.Nickname})
			apikit.ShouldHaveRequiredError(t, err, "Name")
		}

		{
			_, err := s.Create(&u, &xpo.XUserProfileParams{Name: "a_&", Nickname: d.Nickname})
			apikit.ShouldHaveInvalidFormatError(t, err, "Name", "username_format")
		}

		{
			_, err := s.Create(&u, &xpo.XUserProfileParams{Name: "admin", Nickname: d.Nickname})
			if reflect.TypeOf(err) != reflect.TypeOf(&gaekit.DuplicatedObjectError{}) {
				t.Fatalf("It should be gaekit.DuplicatedObjectError: %v, %v", reflect.TypeOf(err), err)
			}
		}
	}

	{
		t.Logf("Nickname")
		{
			_, err := s.Create(&u, &xpo.XUserProfileParams{Name: d.Name})
			apikit.ShouldHaveRequiredError(t, err, "Nickname")
		}

		{
			_, err := s.Create(&u, &xpo.XUserProfileParams{Name: d.Name, Nickname: ""})
			apikit.ShouldHaveRequiredError(t, err, "Nickname")
		}

		{
			_, err := s.Create(&u, &xpo.XUserProfileParams{Name: d.Name, Nickname: "<nynickname"})
			apikit.ShouldHaveInvalidFormatError(t, err, "Nickname", "usernickname_format")
		}

		{
			_, err := s.Create(&u, &xpo.XUserProfileParams{Name: d.Name, Nickname: "reports"}) // reports is blocked
			if reflect.TypeOf(err) != reflect.TypeOf(&gaekit.DuplicatedObjectError{}) {
				t.Fatalf("It should be gaekit.DuplicatedObjectError: %v, %v", reflect.TypeOf(err), err)
			}
		}

	}
}

func checkXUser(t *testing.T, s *xpo.XUserService, f *TestFactory, u user.User, ret xpo.XUser, d xpo.XUser) {
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

	used, err := s.IsUsedName(d.Name)
	if err != nil {
		t.Fatal(err)
	}
	if !used {
		t.Fatalf("It should Name is Used!: %v", d.Name)
	}
	t.Logf("Update!!!")

	{
		t.Logf("Duplicaed")
		_, err := s.Create(&u, &xpo.XUserProfileParams{Name: d.Name, Nickname: d.Nickname})
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
		_, err := s.Create(&u, &xpo.XUserProfileParams{Name: d.Name, Nickname: d.Nickname})
		if err == nil {
			t.Fatal("It should error on creating duplicated name user")
		}

		if reflect.TypeOf(err) != reflect.TypeOf(&gaekit.ValueNotUniqueError{}) {
			t.Fatalf("It should be ValueNotUniqueError: %v", reflect.TypeOf(err))
		}
	}
}
