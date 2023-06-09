package gotext

import (
	"bytes"
	"go/format"
	"text/template"
)

type Template interface {
	Generate(state interface{}) ([]byte, error)
}

// MustParse panics if unable to parse
func MustParse(name, code string) Template {
	template, err := Parse(name, code)
	if err != nil {
		panic(err)
	}
	return template
}

// Parse parses Go code
func Parse(name, code string) (Template, error) {
	tpl, err := template.New(name).Parse(code)
	if err != nil {
		return nil, err
	}
	return &gotext{name, tpl}, nil
}

// Template struct
type gotext struct {
	name string
	tpl  *template.Template
}

// Name returns the name of the template
func (t *gotext) Name() string {
	return t.name
}

// Generate the code
func (t *gotext) Generate(state interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := t.tpl.Execute(buf, state); err != nil {
		return nil, err
	}
	if val, err := format.Source(buf.Bytes()); err == nil {
		return val, nil
	}
	return buf.Bytes(), nil
}
