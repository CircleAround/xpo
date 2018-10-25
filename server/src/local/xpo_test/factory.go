package xpo

import (
	"fmt"
	"local/gaekit"
	"local/xpo/entities"

	"golang.org/x/net/context"
)

type TestFactory struct {
	gaekit.DatastoreAccessObject
	XUserCounter  int64
	ReportCounter int64
}

var singleton *TestFactory

func NewTestFactory() *TestFactory {
	if singleton != nil {
		return singleton
	}

	f := new(TestFactory)
	f.XUserCounter = 0

	singleton = f
	return f
}

func (f *TestFactory) CreateXUser(c context.Context) (entities.XUser, error) {
	xu := f.BuildXUser()
	err := f.Put(c, &xu)
	return xu, err
}

func (f *TestFactory) BuildXUser() entities.XUser {
	f.XUserCounter++
	return entities.XUser{
		Email:    fmt.Sprintf("test%v,@example.com", f.XUserCounter),
		ID:       fmt.Sprintf("%v", f.XUserCounter),
		Name:     fmt.Sprintf("name%v", f.XUserCounter),
		Nickname: fmt.Sprintf("ニックネーム%v", f.XUserCounter),
	}
}

func (f *TestFactory) BuildReport() entities.Report {
	f.ReportCounter++
	return entities.Report{
		Content:     fmt.Sprintf("This is Content %v", f.ReportCounter),
		ContentType: "text/x-markdown",
		ID:          f.ReportCounter,
	}
}

func (f *TestFactory) BuildReportWithAuthor(c context.Context, xu *entities.XUser) entities.Report {
	report := f.BuildReport()
	report.AuthorKey = f.KeyOf(c, xu)
	report.AuthorID = xu.ID
	report.Author = xu.Name
	report.AuthorNickname = xu.Nickname
	return report
}
