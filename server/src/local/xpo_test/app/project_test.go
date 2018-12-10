package app_test

import (
	"local/gaekit"
	"local/testkit"
	"local/xpo/app"
	"local/xpo/domain/project"
	"local/xpo/store"
	"local/xpo_test"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	null "gopkg.in/guregu/null.v3"
)

func TestProjectScenario(t *testing.T) {
	_, c, done := testkit.StartTest(t)
	defer done()

	t.Log("ProjectScenario")
	f := xpo.NewTestFactory()

	xu, err := f.CreateXUser(c)
	if err != nil {
		testkit.Fatal(t, err)
	}

	s := app.NewProjectService()
	pr := store.NewProjectRepository()

	{
		d := f.BuildProject()

		t.Log("Create")
		p, err := s.Create(c, &xu, &project.Params{
			Name:          d.Name,
			Description:   null.StringFrom(d.Description),
			RepositoryURL: null.StringFrom(d.RepositoryURL),
		})
		{
			if err != nil {
				testkit.Fatal(t, err)
			}
			if *p.OwnerKey != *pr.KeyOf(c, xu) {
				t.Errorf("It should have author's key: %v", p)
			}
			if p.Name != d.Name {
				t.Errorf("It should have specified Name: %v", p.Name)
			}
			if p.Description != d.Description {
				t.Errorf("It should have specified Description: %v", p.Description)
			}
			if p.RepositoryURL != d.RepositoryURL {
				t.Errorf("It should have specified RepositoryURL: %v", p.RepositoryURL)
			}

		}

		{
			t.Log("Create Duplicated")
			_, err := s.Create(c, &xu, &project.Params{Name: d.Name})
			{
				if err == nil {
					t.Fatal("It should have error")
				}

				err = errors.Cause(err)

				if reflect.TypeOf(err) != reflect.TypeOf(&gaekit.ValueNotUniqueError{}) {
					t.Fatalf("It should be ValueNotUniqueError: %v", reflect.TypeOf(err))
				}
			}
		}

		{
			t.Log("Update")
			d2 := f.BuildProject()
			p2, err := s.Update(c, &xu, &project.UpdatingParams{
				Params: project.Params{Name: d2.Name, Description: null.StringFrom(d2.Description)},
				ID:     p.ID,
			})
			if err != nil {
				testkit.Fatal(t, err)
			}
			if *p2.OwnerKey != *pr.KeyOf(c, xu) {
				t.Errorf("It should not change author's key: %v", p)
			}
			if p2.Name != d2.Name {
				t.Errorf("It should have specified Name: %v", p2.Name)
			}
			if p2.Description != d2.Description {
				t.Errorf("It should have specified Description: %v != %v", p2.Description, d2.Description)
			}
			if p2.RepositoryURL != d.RepositoryURL {
				t.Errorf("It should not have changed RepositoryURL: %v != %v", p2.RepositoryURL, d.RepositoryURL)
			}
		}

		{
			ps, err := s.SearchByOwnerID(c, xu.ID)
			if err != nil {
				testkit.Fatal(t, err)
			}

			if len(ps) != 1 {
				t.Errorf("It shoud have 1 project: %v", ps)
			}

			if ps[0].ID != p.ID {
				t.Errorf("It shoud have first object: %v", ps[0])
			}
		}
	}
}
