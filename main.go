package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/kanata2/esa-freshness-patroller/internal/esa"
)

func main() {
	log.SetFlags(0)
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	ctx := context.Background()

	cfg, err := newConfigFrom(args)
	if err != nil {
		return err
	}

	if cfg.EsaApiKey == "" || cfg.Team == "" {
		return fmt.Errorf("esa API key and team must be set")
	}
	var (
		transformer Transformer = &jsonTransformer{}
		sender      Sender      = &stdoutSender{}
	)

	app := app{
		config:      cfg,
		debug:       cfg.Debug,
		client:      esa.NewClient(cfg.Team, cfg.EsaApiKey),
		checker:     &checker{},
		transformer: transformer,
		sender:      sender,
		logger:      log.Default(),
	}

	if cfg.Output == "go-template" {
		if cfg.Template == "" {
			return fmt.Errorf("must set template when specify go-template output type")
		}
		tmpl, err := template.ParseFiles(cfg.Template)
		if err != nil {
			return err
		}
		app.transformer = &goTemplateTransformer{tmpl: tmpl}
	}

	if cfg.Destination == "esa" {
		if cfg.Esa == nil || cfg.Esa.ReportPostNumber == 0 {
			return fmt.Errorf("must set a number of esa post for updating when specify esa destination type")
		}
		app.sender = &esaSender{
			client:           app.client,
			reportPostNumber: cfg.Esa.ReportPostNumber,
		}
	}

	return app.run(ctx, args)
}

type app struct {
	debug       bool
	config      *config
	client      *esa.Client
	checker     Checker
	template    *template.Template
	transformer Transformer
	sender      Sender
	logger      *log.Logger
}

func (app app) Debugf(format string, v ...interface{}) {
	if app.debug {
		app.logger.Printf("[debug] "+format, v...)
	}
}

func (app app) Infof(format string, v ...interface{}) {
	app.logger.Printf("[info] "+format, v...)
}

func (app app) Warnf(format string, v ...interface{}) {
	app.logger.Printf("[warn] "+format, v...)
}

func (app app) run(ctx context.Context, args []string) error {
	page := 1
	var result Result
	for {
		resp, err := app.client.ListPosts(
			ctx,
			esa.WithListPostsOptionQuery(app.config.Query),
			esa.WithListPostsOptionPage(page),
			esa.WithListPostsOptionPerPage(esa.MaxElementsPerPage),
		)
		if err != nil {
			return err
		}
		app.Debugf("Hit %d posts", len(resp.Posts))
		for _, p := range resp.Posts {
			mo, err := app.checker.Check(p.BodyMarkdown)
			if err != nil {
				result.Warnings = append(result.Warnings, &Warning{
					Title:  p.Name,
					URL:    p.URL,
					Reason: err.Error(),
				})
				continue
			}
			if mo == nil {
				continue
			}
			mo.Title = p.Name
			mo.URL = p.URL
			result.Items = append(result.Items, mo)
		}
		if resp.NextPage == nil {
			break
		}
		page = *resp.NextPage
	}
	r, err := app.transformer.Transform(result)
	if err != nil {
		return err
	}
	return app.sender.Send(r)
}

type Result struct {
	Items    []*MaybeOutdated `json:"items"`
	Warnings []*Warning       `json:"warnings"`
}

type MaybeOutdated struct {
	Title         string    `json:"title"`
	URL           string    `json:"url"`
	LastCheckedAt time.Time `json:"last_checked_at"`
	Owner         string    `json:"owner"`
}

type Warning struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Reason string `json:"reason"`
}
