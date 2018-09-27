package xpo_test

import (
	"fmt"
	"local/gaekit"
	"local/testkit"
	"local/xpo"
	"testing"

	"golang.org/x/net/context"
)

type TestFactory struct {
	gaekit.AppEngineService
	XUserCounter  int64
	ReportCounter int64
}

var singleton *TestFactory

func NewTestFactory(c context.Context) *TestFactory {
	if singleton != nil {
		return singleton
	}

	f := new(TestFactory)
	f.InitAppEngineService(c)
	f.XUserCounter = 0

	singleton = f
	return f
}

func (f *TestFactory) CreateXUser() (xpo.XUser, error) {
	xu := f.BuildXUser()
	err := f.Put(&xu)
	return xu, err
}

func (f *TestFactory) BuildXUser() xpo.XUser {
	f.XUserCounter++
	return xpo.XUser{
		Email:    fmt.Sprintf("test%v,@example.com", f.XUserCounter),
		ID:       fmt.Sprintf("%v", f.XUserCounter),
		Name:     fmt.Sprintf("name%v", f.XUserCounter),
		Nickname: fmt.Sprintf("ニックネーム%v", f.XUserCounter),
	}
}

func (f *TestFactory) BuildReport() xpo.Report {
	f.ReportCounter++
	return xpo.Report{
		Content:     fmt.Sprintf("This is Content %v", f.ReportCounter),
		ContentType: "text/x-markdown",
		ID:          f.ReportCounter,
	}
}

func TestMain(m *testing.M) {
	testkit.BootstrapTest(m)
}
