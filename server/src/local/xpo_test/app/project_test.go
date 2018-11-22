package app_test

import (
	"local/testkit"
	"local/xpo/app"
	"local/xpo/store"
	"local/xpo_test"
	"testing"
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
		p, err := s.Create(c, &xu, &app.ProjectCreationParams{Name: d.Name})
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
			t.Log("Update")
			d2 := f.BuildProject()
			p2, err := s.Update(c, &xu, &app.ProjectUpdatingParams{
				ProjectCreationParams: app.ProjectCreationParams{Name: d2.Name},
				ID:                    p.ID,
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
				t.Errorf("It should have specified Description: %v", p2.Description)
			}
			if p2.RepositoryURL != d2.RepositoryURL {
				t.Errorf("It should have specified RepositoryURL: %v", p2.RepositoryURL)
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
