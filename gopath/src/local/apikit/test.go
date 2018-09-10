package apikit

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/favclip/testerator"
	"google.golang.org/appengine/aetest"

	validator "gopkg.in/go-playground/validator.v9"
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

// ShouldHaveValidationError is matcher for property with tag
func ShouldHaveValidationError(t *testing.T, err error, property string, tag string) {
	t.Logf("for %v", property)
	if err == nil {
		t.Fatal("It should have error")
	}

	if reflect.TypeOf(err) != reflect.TypeOf(validator.ValidationErrors{}) {
		t.Fatalf("It should be validator.ValidationErrors: %v, %v", reflect.TypeOf(err), err)
	}

	ei := err.(validator.ValidationErrors)[0]
	t.Logf("kind: %v, type: %v, value: %v, param: %v", ei.Kind(), ei.Type(), ei.Value(), ei.Param())
	if ei.Field() != property {
		t.Fatalf("It should not be nil. should have property: %v, %v", ei.Field(), property)
	}

	if ei.Tag() != tag {
		t.Fatalf("It should be tag %v in %v", tag, ei.Tag())
	}
}

// ShouldHaveRequiredError is matcher for property with required
func ShouldHaveRequiredError(t *testing.T, err error, property string) {
	ShouldHaveValidationError(t, err, property, "required")
}

// ShouldHaveInvalidFormatError is matcher for property with invalid format
func ShouldHaveInvalidFormatError(t *testing.T, err error, property string, tag string) {
	ShouldHaveValidationError(t, err, property, tag)
}

// ShouldHaveTooLongError is matcher for property with too long
func ShouldHaveTooLongError(t *testing.T, err error, property string) {
	ShouldHaveValidationError(t, err, property, "max")
}
