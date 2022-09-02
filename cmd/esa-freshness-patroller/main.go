package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"

	patroller "github.com/kanata2/esa-freshness-patroller"
	"github.com/kanata2/esa-freshness-patroller/internal/config"
	"github.com/kanata2/esa-freshness-patroller/internal/esa"
	"github.com/kanata2/esa-freshness-patroller/internal/output"
	"github.com/kanata2/esa-freshness-patroller/internal/transform"
)

func main() {
	log.SetFlags(0)
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(argv []string) error {
	ctx := context.Background()
	cfg, err := config.New(argv)
	if err != nil {
		return err
	}

	if cfg.EsaApiKey == "" || cfg.Team == "" {
		return fmt.Errorf("esa API key and team must be specified")
	}

	var transformer transform.Transformer = transform.NewJSONTransformer()
	if cfg.Output == "go-template" {
		if cfg.Template == "" {
			return fmt.Errorf("must set template when specify go-template output type")
		}
		tmpl, err := template.ParseFiles(cfg.Template)
		if err != nil {
			return err
		}
		transformer = transform.NewGoTemplateTransformer(tmpl)
	}
	var outputter output.Outputter = output.NewStdoutOutputter()
	if cfg.Destination == "esa" {
		if cfg.Esa == nil || cfg.Esa.ReportPostNumber == 0 {
			return fmt.Errorf("must set a number of esa post for updating when specify esa destination type")
		}
		c := esa.NewClient(cfg.Team, cfg.EsaApiKey)
		outputter = output.NewEsaOutputter(c, cfg.Esa.ReportPostNumber)
	}

	opts := []patroller.PatrollerOptionFn{}

	if cfg.Debug {
		opts = append(opts, patroller.WithDebug())
	}

	p := patroller.New(cfg.EsaApiKey, cfg.Team, cfg.Query, opts...)

	result, err := p.Patrol(ctx)
	if err != nil {
		return err
	}
	r, err := transformer.Transform(result)
	if err != nil {
		return err
	}
	if err := outputter.Output(r); err != nil {
		return err
	}
	return nil
}
