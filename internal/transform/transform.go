package transform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"

	patroller "github.com/kanata2/esa-freshness-patroller"
)

type Transformer interface {
	Transform(*patroller.Result) (io.Reader, error)
}

type jsonTransformer struct{}

func NewJSONTransformer() *jsonTransformer {
	return &jsonTransformer{}
}

func (t *jsonTransformer) Transform(r *patroller.Result) (io.Reader, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(&r); err != nil {
		return nil, fmt.Errorf("jsonTransformer.Transform: %w", err)
	}
	return buf, nil
}

type goTemplateTransformer struct {
	tmpl *template.Template
}

func NewGoTemplateTransformer(t *template.Template) *goTemplateTransformer {
	return &goTemplateTransformer{tmpl: t}
}

func (t *goTemplateTransformer) Transform(r *patroller.Result) (io.Reader, error) {
	buf := new(bytes.Buffer)
	if err := t.tmpl.Execute(buf, &r); err != nil {
		return nil, fmt.Errorf("goTemplateTransformer.Transform: %w", err)
	}
	return buf, nil
}
