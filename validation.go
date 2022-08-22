package val

import (
	"errors"
	"fmt"
	"strings"
)

const (
	field_sep = "."
	print_sep = ", "
)

var (
	// base error for all validations errors
	Err = errors.New("validation error")
)

type Errors struct {
	fields []Field
}

type Field struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (*Errors) Unwrap() error {
	return Err
}

func (e *Errors) Any() bool {
	return e != nil && len(e.fields) > 0
}

func (e Errors) GetFields() []Field {
	if !e.Any() {
		return nil
	}
	fields := make([]Field, len(e.fields))
	copy(fields, e.fields)
	return fields
}

func (e *Errors) Error() string {
	if e == nil {
		return "<nil> validation error"
	}

	sb := make([]string, 0, len(e.fields))
	for _, x := range e.fields {
		sb = append(sb, fmt.Sprintf("'%s' %s", x.Name, x.Description))
	}
	if len(e.fields) == 1 {
		return "validation error: " + strings.Join(sb, print_sep)
	}
	return "validation errors: " + strings.Join(sb, print_sep)
}

// new field error with name and description. The name field may be empty if it is part of a hierachy and specified using the NewChildren function
func New(name, format string, args ...interface{}) *Errors {
	return FromFields(Field{Name: name, Description: fmt.Sprintf(format, args...)})
}

func (e *Errors) Addf(field, format string, args ...interface{}) *Errors {
	return Concat(e, New(field, format, args...))
}

/*
new children validation errors with the parent prefix(es) to the field name

example usage:

	NewChildren("parent",
	  NewChildren("field1", x.Stuff.Validate())) // where x.Stuff.Validate() returns *Errors
*/
func NewChildren(parent string, xs ...*Errors) *Errors {
	fields := make([]Field, 0, len(xs))
	for _, x := range xs {
		if x.Any() {
			for _, f := range x.fields {
				if f.Name == "" {
					fields = append(fields, Field{Name: parent, Description: f.Description})
				} else {
					fields = append(fields, Field{Name: parent + field_sep + f.Name, Description: f.Description})
				}
			}
		}
	}
	return FromFields(fields...)
}

// concat errors. If all errors are empty/nil, nil is returned
func Concat(xs ...*Errors) *Errors {
	fields := make([]Field, 0)
	for _, x := range xs {
		fields = append(fields, x.fields...)
	}
	return FromFields(fields...)
}

func FromFields(xs ...Field) *Errors {
	if len(xs) == 0 {
		return nil
	}
	fields := make([]Field, len(xs))
	copy(fields, xs)
	return &Errors{fields: fields}
}
