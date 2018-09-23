package testkit

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/favclip/testerator"
	"google.golang.org/appengine/aetest"
)

// BootstrapTest is a method for spinup a first testerator
func BootstrapTest(m *testing.M) {
	_, _, err := testerator.SpinUp()

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	status := m.Run()

	err = testerator.SpinDown()
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	os.Exit(status)
}

// StartTest is a methid for spinup new testerator. with panic handling
func StartTest(t *testing.T) (aetest.Instance, context.Context, func()) {
	i, c, err := testerator.SpinUp()
	if err != nil {
		t.Fatal(err)
	}

	done := func() {
		testerator.SpinDown()

		err := recover()
		if err != nil {
			pc, fn, line, _ := runtime.Caller(1)
			t.Fatalf("Panic!: %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
		}
	}

	return i, c, done
}