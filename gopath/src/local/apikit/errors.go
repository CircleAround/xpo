package apikit

import (
	"fmt"
	"strings"
)

const (
	// Required is specify required error
	Required string = "required"

	InvalidFormat string = "invalid_format"

	// TooLong is specify too long contents
	TooLong string = "toolong"
)

// ValidationError is for validation error
type ValidationError struct {
	Type  string                          `json:"type"`
	Items map[string]*ValidationErrorItem `json:"items"`
}

func (err ValidationError) Error() string {
	str := ""
	for _, v := range err.Items {
		str += fmt.Sprintf("%v|", v)
	}
	return strings.TrimRight(str, "|")
}

// HasItem is a checker for err has ValidationErrorItem
func (err ValidationError) HasItem() bool {
	return len(err.Items) != 0
}

// PushItem is a function for push new ValidationErrorItem to err
func (err *ValidationError) PushItem(property string, reasons []string) {
	item := err.Items[property]
	if item == nil {
		err.Items[property] = &ValidationErrorItem{Property: property, Reasons: reasons}
	} else {
		item.Reasons = append(item.Reasons, reasons...)
	}
}

// PushOne is a function for push new ValidationErrorItem
func (err *ValidationError) PushOne(property string, reasons string) {
	err.PushItem(property, []string{reasons})
}

// NewValidationError if a function for creates ValidatioNError
func NewValidationError() *ValidationError {
	err := new(ValidationError)
	err.Type = "ValidationError"
	err.Items = map[string]*ValidationErrorItem{}
	return err
}

// ValidationErrorItem is inner struct for ValidationError
type ValidationErrorItem struct {
	Property string   `json:"property"`
	Reasons  []string `json:"reasons"`
}

func (eri ValidationErrorItem) Error() string {
	return fmt.Sprintf("%s %s", eri.Property, eri.JoinReasons())
}

// JoinReasons is a function for getting joinded reasons to string
func (eri *ValidationErrorItem) JoinReasons() string {
	return strings.Join(eri.Reasons, ",")
}
