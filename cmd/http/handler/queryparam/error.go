package queryparam

import (
	"fmt"
	"reflect"
)

// ParseError occurs when it's impossible to convert the value for given type.
type ParseError struct {
	Name string
	Type reflect.Type
	Err  error
}

func newParseError(sf reflect.StructField, err error) error {
	return ParseError{sf.Name, sf.Type, err}
}

// Error returns the error message.
func (e ParseError) Error() string {
	return fmt.Sprintf(`parse error on field "%s" of type "%s": %v`, e.Name, e.Type, e.Err)
}

// NotStructPtrError occurs when pass something that is not a pointer to a Struct to Parse
type NotStructPtrError struct{}

// Error returns the error message.
func (e NotStructPtrError) Error() string {
	return "expected a pointer to a Struct"
}

// NoParserError occurs when there is no parser provided for given type
type NoParserError struct {
	Name string
	Type reflect.Type
}

func newNoParserError(sf reflect.StructField) error {
	return NoParserError{sf.Name, sf.Type}
}

// Error returns the error message.
func (e NoParserError) Error() string {
	return fmt.Sprintf(`no parser found for field "%s" of type "%s"`, e.Name, e.Type)
}

// NoSupportedTagOptionError occurs when the given tag is not supported
// In-built supported tags: "", "required", "unset", "notEmpty", "expand", "paramDefault"
type NoSupportedTagOptionError struct {
	Tag string
}

func newNoSupportedTagOptionError(tag string) error {
	return NoSupportedTagOptionError{tag}
}

// Error returns the error message.
func (e NoSupportedTagOptionError) Error() string {
	return fmt.Sprintf("tag option %q not supported", e.Tag)
}
