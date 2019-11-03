package template

import "strings"

type Template struct {
	variables map[string]string
}

func NewTemplate() *Template {
	return &Template{
		variables: make(map[string]string),
	}
}

func (t *Template) WithVariables(variables map[string]string) *Template {
	t.variables = variables
	return t
}

func (t *Template) AddVariable(key, value string) {
	if t.variables == nil {
		t.variables = make(map[string]string)
	}
	t.variables[key] = value
}

func (t *Template) Replace(s string) string {
	for k, v := range t.variables {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}
