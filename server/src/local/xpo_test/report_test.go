package xpo_test

import (
	"local/testkit"
	"local/timekit"
	"local/xpo"
	"testing"
	"time"

	"google.golang.org/appengine/datastore"
)

func TestReportScenario(t *testing.T) {
	_, c, done := testkit.StartTest(t)
	defer done()

	t.Log("ReportScenario")
	f := NewTestFactory(c)

	xu, err := f.CreateXUser()
	if err != nil {
		t.Fatal(err)
	}

	d := f.BuildReport()
	tp := new(timekit.TestTimeProvider)

	tm, err := time.Parse("2006-01-02 15:04:05 MST", "2014-12-31 12:31:24 JST")
	if err != nil {
		t.Fatal(err)
	}
	oneHourBefore := tm.Add(-1 * time.Hour)

	tp.StopAt(oneHourBefore)
	s := xpo.NewReportServiceWithTimeProvider(c, tp)
	{
		t.Log("Create")
		r, err := s.Create(xu, xpo.ReportCreationParams{Content: d.Content})
		if err != nil {
			t.Fatal(err)
		}
		if *r.AuthorKey != *s.KeyOf(xu) {
			t.Errorf("It should have author's key: %v", r)
		}
		if r.Content != d.Content {
			t.Errorf("It should have specified Content: %v", r.Content)
		}
		if r.Year != 2014 {
			t.Errorf("It should have specified Year: %v", r.Year)
		}
		if r.Month != 12 {
			t.Errorf("It should have specified Month: %v", r.Month)
		}
		if r.Day != 31 {
			t.Errorf("It should have specified Day: %v", r.Day)
		}
		if r.YearDay != 365 {
			t.Errorf("It should have specified YearDay: %v", r.YearDay)
		}

		now := tp.StopAt(tm)
		t.Logf("Update: %v", r)
		{
			before := *r
			ud := f.BuildReport()
			p := xpo.ReportUpdatingParams{
				ReportCreationParams: xpo.ReportCreationParams{ud.Content},
				ID:                   r.ID,
			}

			{
				r, err := s.Update(xu, p)
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
				if r.UpdatedAt == before.UpdatedAt {
					t.Errorf("It should change UpdatedAt")
				}
			}

			{
				other := f.BuildXUser()
				_, err := s.Update(other, p)
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
			hit, err := s.Find(xu.ID, r.ID)
			if err != nil {
				t.Fatal(err)
			}

			if *hit.AuthorKey != *s.KeyOf(xu) {
				t.Error("It should be equal AuthorKey and Key of xu")
			}
			if hit.ID != r.ID {
				t.Error("It should be equal ID")
			}
		}
	}
}
