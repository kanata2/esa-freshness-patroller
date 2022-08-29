package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
)

type Outputter interface {
	Output(Result) error
}

type jsonOutputter struct {
	out io.Writer
}

func (o *jsonOutputter) Output(r Result) error {
	if err := json.NewEncoder(o.out).Encode(&r); err != nil {
		return fmt.Errorf("jsonOutputter.Output: %w", err)
	}
	return nil
}

type goTemplateOutputter struct {
	tmpl *template.Template
	out  io.Writer
}

func (o *goTemplateOutputter) Output(r Result) error {
	return o.tmpl.Execute(o.out, &r)
}
