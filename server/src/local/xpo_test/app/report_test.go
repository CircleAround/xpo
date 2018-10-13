package app_test

import (
	"local/testkit"
	"local/the_time"
	"local/xpo/app"
	"local/xpo_test"
	"testing"
	"time"

	"google.golang.org/appengine/datastore"
)

func TestReportScenario(t *testing.T) {
	_, c, done := testkit.StartTest(t)
	defer done()

	t.Log("ReportScenario")
	f := xpo.NewTestFactory()

	xu, err := f.CreateXUser(c)
	if err != nil {
		t.Fatal(err)
	}

	d := f.BuildReport()
	tp := the_time.Machine()

	tm, err := time.Parse("2006-01-02 15:04:05 MST", "2014-12-31 12:31:24 JST")
	if err != nil {
		t.Fatal(err)
	}
	oneHourBefore := tm.Add(-1 * time.Hour)

	tp.TravelTo(oneHourBefore)
	s := app.NewReportServiceWithTheTime(tp)
	{
		t.Log("Create")
		r, err := s.Create(c, xu, app.ReportCreationParams{Content: d.Content, ContentType: d.ContentType})
		if err != nil {
			t.Fatal(err)
		}
		if *r.AuthorKey != *s.KeyOf(c, xu) {
			t.Errorf("It should have author's key: %v", r)
		}
		if r.Content != d.Content {
			t.Errorf("It should have specified Content: %v", r.Content)
		}
		if r.ContentType != d.ContentType {
			t.Errorf("It should have specified ContentType: %v", r.ContentType)
		}
		if !r.CreatedAt.Equal(oneHourBefore) {
			t.Errorf("It should have specified CreateddAt: %v", r.CreatedAt)
		}
		if !r.UpdatedAt.Equal(oneHourBefore) {
			t.Errorf("It should have specified UpdatedAt: %v", r.UpdatedAt)
		}
		if !r.ReportedAt.Equal(oneHourBefore) {
			t.Errorf("It should have specified ReportedAt: %v", r.ReportedAt)
		}

		now := tp.TravelTo(tm)
		t.Logf("Update: %v", r)
		{
			before := *r
			ud := f.BuildReport()
			p := app.ReportUpdatingParams{
				ReportCreationParams: app.ReportCreationParams{Content: ud.Content, ContentType: ud.ContentType},
				ID:                   r.ID,
			}

			{
				r, err := s.Update(c, xu, p)
				if err != nil {
					t.Fatal(err)
				}
				if r.ID != before.ID {
					t.Errorf("It should not change ID %v ,%v", before.ID, ud.ID)
				}
				if *r.AuthorKey != *before.AuthorKey {
					t.Errorf("It should not change AuthorKey %v ,%v", *before.AuthorKey, *ud.AuthorKey)
				}
				if r.Content == before.Content {
					t.Errorf("It should change Content from %v to %v", before.Content, ud.Content)
				}
				if !r.UpdatedAt.Equal(now) {
					t.Errorf("It should change UpdatedAt to current Time: %v, %v", r.UpdatedAt, now)
				}
				if r.UpdatedAt.Equal(before.UpdatedAt) {
					t.Errorf("It should change UpdatedAt")
				}
			}

			{
				other := f.BuildXUser()
				_, err := s.Update(c, other, p)
				if err == nil {
					t.Error("It should block update by other user, must throw error")
				}
				if err != datastore.ErrNoSuchEntity {
					t.Errorf("It should block update by other user: %v", err)
				}
			}

		}

		{
			t.Logf("Find: %v", r)
			hit, err := s.Find(c, xu.ID, r.ID)
			if err != nil {
				t.Fatal(err)
			}

			if *hit.AuthorKey != *s.KeyOf(c, xu) {
				t.Error("It should be equal AuthorKey and Key of xu")
			}
			if hit.ID != r.ID {
				t.Error("It should be equal ID")
			}
		}
	}

	now := tp.TravelTo(tm)
	{
		t.Log("Create with ReportedAt")
		r, err := s.Create(c, xu, app.ReportCreationParams{Content: d.Content, ContentType: d.ContentType, ReportedAt: oneHourBefore})
		if err != nil {
			t.Fatal(err)
		}
		if *r.AuthorKey != *s.KeyOf(c, xu) {
			t.Errorf("It should have author's key: %v", r)
		}
		if r.Content != d.Content {
			t.Errorf("It should have specified Content: %v", r.Content)
		}
		if r.ContentType != d.ContentType {
			t.Errorf("It should have specified ContentType: %v", r.ContentType)
		}
		if !r.CreatedAt.Equal(now) {
			t.Errorf("It should have specified CreateddAt: %v", r.CreatedAt)
		}
		if !r.UpdatedAt.Equal(now) {
			t.Errorf("It should have specified UpdatedAt: %v", r.UpdatedAt)
		}
		if !r.ReportedAt.Equal(oneHourBefore) {
			t.Errorf("It should have specified ReportedAt: %v", r.ReportedAt)
		}
	}

	{
		// Create one by another user
		oxu, err := f.CreateXUser(c)
		if err != nil {
			t.Fatal(err)
		}
	
		d := f.BuildReport()
		_, err = s.Create(c, oxu, app.ReportCreationParams{Content: d.Content, ContentType: d.ContentType, ReportedAt: oneHourBefore})
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		t.Log("Search By")
		rs, err := s.SearchBy(c, xu.ID, now.Year(), int(now.Month()), now.Day())
		if err != nil {
			t.Fatal(err)
		}

		if len(rs) != 2 {
			t.Errorf("It should have 2 results: %v", rs)
		}
	}

	{
		t.Log("Search By Author")
		rs, err := s.SearchByAuthor(c, xu.ID)
		if err != nil {
			t.Fatal(err)
		}

		if len(rs) != 2 {
			t.Errorf("It should have 2 results: %v", rs)
		}
	}
}
