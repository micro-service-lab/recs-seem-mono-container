package queryparam

import (
	"fmt"
	"reflect"
)

// The error occurs when it's impossible to convert the value for given type.
type ParseError struct {
	Name string
	Type reflect.Type
	Err  error
}

func newParseError(sf reflect.StructField, err error) error {
	return ParseError{sf.Name, sf.Type, err}
}

func (e ParseError) Error() string {
	return fmt.Sprintf(`parse error on field "%s" of type "%s": %v`, e.Name, e.Type, e.Err)
}

// The error occurs when pass something that is not a pointer to a Struct to Parse
type NotStructPtrError struct{}

func (e NotStructPtrError) Error() string {
	return "expected a pointer to a Struct"
}

// This error occurs when there is no parser provided for given type
type NoParserError struct {
	Name string
	Type reflect.Type
}

func newNoParserError(sf reflect.StructField) error {
	return NoParserError{sf.Name, sf.Type}
}

func (e NoParserError) Error() string {
	return fmt.Sprintf(`no parser found for field "%s" of type "%s"`, e.Name, e.Type)
}

// This error occurs when the given tag is not supported
// In-built supported tags: "", "required", "unset", "notEmpty", "expand", "paramDefault"
type NoSupportedTagOptionError struct {
	Tag string
}

func newNoSupportedTagOptionError(tag string) error {
	return NoSupportedTagOptionError{tag}
}

func (e NoSupportedTagOptionError) Error() string {
	return fmt.Sprintf("tag option %q not supported", e.Tag)
}

// This error occurs when the required variable is not set
type QueryParamIsNotSetError struct {
	Key string
}

func newQueryParamIsNotSet(key string) error {
	return QueryParamIsNotSetError{key}
}

func (e QueryParamIsNotSetError) Error() string {
	return fmt.Sprintf(`required query parameter %q is not set`, e.Key)
}

// This error occurs when the variable which must be not empty is existing but has an empty value
type EmptyQueryParamError struct {
	Key string
}

func newEmptyQueryParamError(key string) error {
	return EmptyQueryParamError{key}
}

func (e EmptyQueryParamError) Error() string {
	return fmt.Sprintf("query param %q should not be empty", e.Key)
}

// This error occurs when it's impossible to convert value using given parser
type ParseValueError struct {
	Msg string
	Err error
}

func newParseValueError(message string, err error) error {
	return ParseValueError{message, err}
}

func (e ParseValueError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}
