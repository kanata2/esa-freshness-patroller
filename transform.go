package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
)

type Transformer interface {
	Transform(Result) (io.Reader, error)
}

type jsonTransformer struct{}

func (t *jsonTransformer) Transform(r Result) (io.Reader, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(&r); err != nil {
		return nil, fmt.Errorf("jsonTransformer.Transform: %w", err)
	}
	return buf, nil
}

type goTemplateTransformer struct {
	tmpl *template.Template
}

func (t *goTemplateTransformer) Transform(r Result) (io.Reader, error) {
	buf := new(bytes.Buffer)
	if err := t.tmpl.Execute(buf, &r); err != nil {
		return nil, fmt.Errorf("goTemplateTransformer.Transform: %w", err)
	}
	return buf, nil
}
