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

	{
		t.Log("Scenario")

		d := f.BuildXUser()

		var u user.User
		u.Email = d.Email
		u.ID = d.ID

		xu, err := s.Create(&u, &xpo.XUserCreationParams{Name: d.Name, Nickname: d.Nickname})
		if err != nil {
			t.Fatal(err)
		}

		{
			t.Logf("Standard")
			ret := &xpo.XUser{ID: xu.ID}
			if err = s.Get(ret); err != nil {
				t.Fatal(err)
			}

			if ret.Email != d.Email {
				t.Fatalf("It should get saved email!: %v", ret.Email)
			}
			if ret.Name != d.Name {
				t.Fatalf("It should get saved Name!: %v", ret.Name)
			}
			if ret.Nickname != d.Nickname {
				t.Fatalf("It should get saved Nickname!: %v", ret.Nickname)
			}

			{
				t.Logf("Duplicaed")
				xu, err = s.Create(&u, &xpo.XUserCreationParams{Name: d.Name, Nickname: d.Nickname})
				if err == nil {
					t.Fatal("It should error on creating duplicated user")
				}

				if reflect.TypeOf(err) != reflect.TypeOf(&gaekit.DuplicatedObjectError{}) {
					t.Fatalf("It should be DuplicatedObjectError: %v", reflect.TypeOf(err))
				}
			}

			{
				t.Logf("Unique email")

				d2 := f.BuildXUser()
				var u user.User
				u.Email = d2.Email
				u.ID = d2.ID

				xu, err = s.Create(&u, &xpo.XUserCreationParams{Name: d.Name, Nickname: d.Nickname})
				if err == nil {
					t.Fatal("It should error on creating duplicated name user")
				}

				if reflect.TypeOf(err) != reflect.TypeOf(&gaekit.ValueNotUniqueError{}) {
					t.Fatalf("It should be ValueNotUniqueError: %v", reflect.TypeOf(err))
				}
			}
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
			_, err := s.Create(&u, &xpo.XUserCreationParams{Nickname: d.Nickname})
			apikit.ShouldHaveRequiredError(t, err, "Name")
		}

		{
			_, err := s.Create(&u, &xpo.XUserCreationParams{Name: "", Nickname: d.Nickname})
			apikit.ShouldHaveRequiredError(t, err, "Name")
		}

		{
			_, err := s.Create(&u, &xpo.XUserCreationParams{Name: "a_&", Nickname: d.Nickname})
			apikit.ShouldHaveInvalidFormatError(t, err, "Name", "username_format")
		}
	}

	{
		t.Logf("Nickname")
		{
			_, err := s.Create(&u, &xpo.XUserCreationParams{Name: d.Name})
			apikit.ShouldHaveRequiredError(t, err, "Nickname")
		}

		{
			_, err := s.Create(&u, &xpo.XUserCreationParams{Name: d.Name, Nickname: ""})
			apikit.ShouldHaveRequiredError(t, err, "Nickname")
		}

		{
			_, err := s.Create(&u, &xpo.XUserCreationParams{Name: d.Name, Nickname: "<nynickname"})
			apikit.ShouldHaveInvalidFormatError(t, err, "Nickname", "usernickname_format")
		}
	}
}
