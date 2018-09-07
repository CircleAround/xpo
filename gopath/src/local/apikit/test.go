package apikit

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/favclip/testerator"
	"google.golang.org/appengine/aetest"
)

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

func StartTest(t *testing.T) (aetest.Instance, context.Context, func()) {
	i, c, err := testerator.SpinUp()
	if err != nil {
		t.Fatal(err)
	}

	done := func() {
		testerator.SpinDown()

		err := recover()
		if err != nil {
			t.Fatalf("Panic!:", err)
		}
	}

	return i, c, done
}

// ShouldHaveValidationError is matcher for property with reason
func ShouldHaveValidationError(t *testing.T, err error, property string, reason string) {
	t.Logf("for %v", property)
	if err == nil {
		t.Fatal("It should have error")
	}

	if reflect.TypeOf(err) != reflect.TypeOf(&ValidationError{}) {
		t.Fatalf("It should be apikit.ValidationError: %v", reflect.TypeOf(err))
	}

	ei := err.(*ValidationError).Items[property]
	if ei == nil {
		t.Fatalf("It should not be nil. should have property: %v", property)
	}

	if !ei.HasReason(reason) {
		t.Fatalf("It should be includes %v in %v", reason, ei.Reasons)
	}
}

// ShouldHaveRequiredError is matcher for property with required
func ShouldHaveRequiredError(t *testing.T, err error, property string) {
	ShouldHaveValidationError(t, err, property, Required)
}

// ShouldHaveInvalidFormatError is matcher for property with invalid format
func ShouldHaveInvalidFormatError(t *testing.T, err error, property string) {
	ShouldHaveValidationError(t, err, property, InvalidFormat)
}

// ShouldHaveTooLongError is matcher for property with too long
func ShouldHaveTooLongError(t *testing.T, err error, property string) {
	ShouldHaveValidationError(t, err, property, TooLong)
}
