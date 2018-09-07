package xpo

import (
	"local/apikit"
	"reflect"
	"testing"

	"google.golang.org/appengine/user"
)

func TestMain(m *testing.M) {
	apikit.BootstrapTest(m)
}

func TestDummy(t *testing.T) {
	t.Logf("Start")
	actual := 30
	expected := 30
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestScenario(t *testing.T) {
	_, c, done := apikit.StartTest(t)
	defer done()

	s := NewXUserService(c)

	{
		t.Log("Scenario")

		var u user.User
		u.Email = "test@example.com"
		u.ID = "1"

		xu, err := s.Create(&u, "myname", "nynickname")
		if err != nil {
			t.Fatal(err)
		}

		{
			t.Logf("Standard")
			ret := &XUser{ID: xu.ID}
			if err = s.Get(ret); err != nil {
				t.Fatal(err)
			}

			if ret.Email != "test@example.com" {
				t.Fatalf("It should get saved email!: %v", ret.Email)
			}
			if ret.Name != "myname" {
				t.Fatalf("It should get saved Name!: %v", ret.Name)
			}
			if ret.NickName != "nynickname" {
				t.Fatalf("It should get saved Nickname!: %v", ret.NickName)
			}

			{
				t.Logf("Duplicaed")
				xu, err = s.Create(&u, "myname", "nynickname")
				if err == nil {
					t.Fatal("It should error on creating duplicated user")
				}

				if reflect.TypeOf(err) != reflect.TypeOf(&DuplicatedObjectError{}) {
					t.Fatalf("It should be DuplicatedObjectError: %v", reflect.TypeOf(err))
				}
			}

			{
				t.Logf("Unique email")

				var u user.User
				u.Email = "test3@example.com"
				u.ID = "3"

				xu, err = s.Create(&u, "myname", "nynickname3")
				if err == nil {
					t.Fatal("It should error on creating duplicated name user")
				}

				if reflect.TypeOf(err) != reflect.TypeOf(&ValueNotUniqueError{}) {
					t.Fatalf("It should be ValueNotUniqueError: %v", reflect.TypeOf(err))
				}
			}
		}
	}
}

func TestValidation(t *testing.T) {
	_, c, done := apikit.StartTest(t)
	defer done()

	s := NewXUserService(c)
	t.Logf("Validation")

	var u user.User
	u.Email = "test2@example.com"
	u.ID = "2"

	{
		_, err := s.Create(&u, "", "mynickname")
		apikit.ShouldHaveRequiredError(t, err, "name")
	}

	{
		_, err := s.Create(&u, "a_&", "mynickname")
		apikit.ShouldHaveInvalidFormatError(t, err, "name")
	}

	{
		_, err := s.Create(&u, "myname", "")
		apikit.ShouldHaveRequiredError(t, err, "nickname")
	}
}
