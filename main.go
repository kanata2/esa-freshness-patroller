package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kanata2/esa-freshness-patroller/internal/esa"
	"github.com/slack-go/slack"
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
	var notifier Notifier = &defaultNotifier{out: os.Stdout}
	if cfg.NotificationType == "slack" {
		if cfg.Slack.Token == "" || cfg.Slack.Channel == "" {
			return fmt.Errorf("slack API key and chaannel must be set")
		}
		notifier = &slackNotifier{
			client:  slack.New(cfg.Slack.Token),
			channel: cfg.Slack.Channel,
		}
	}
	app := app{
		config:   cfg,
		debug:    cfg.Debug,
		client:   esa.NewClient(cfg.EsaApiKey),
		checker:  &checker{},
		notifier: notifier,
		logger:   log.Default(),
	}

	return app.run(ctx, args)
}

type app struct {
	debug    bool
	config   *config
	client   *esa.Client
	checker  Checker
	notifier Notifier
	logger   *log.Logger
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
	outdateCandidates := []*MaybeOutdated{}
	for {
		resp, err := app.client.ListPosts(
			ctx,
			app.config.Team,
			esa.WithListPostsOptionQuery(app.config.Query),
			esa.WithListPostsOptionPage(page),
			esa.WithListPostsOptionPerPage(100),
		)
		if err != nil {
			return err
		}
		app.Debugf("Hit %d posts", len(resp.Posts))
		for _, p := range resp.Posts {
			mo, err := app.checker.Check(p.BodyMarkdown)
			if err != nil {
				app.Warnf("failed to check whether outdated or not. %s(URL: %s) reason: %s", p.Name, p.URL, err)
				continue
			}
			if mo == nil {
				continue
			}
			mo.Title = p.Name
			mo.URL = p.URL
			outdateCandidates = append(outdateCandidates, mo)
		}
		if resp.NextPage == nil {
			break
		}
		page = *resp.NextPage
	}
	if err := app.notifier.Notify(outdateCandidates); err != nil {
		return err
	}
	return nil
}

type MaybeOutdated struct {
	Title         string
	URL           string
	LastCheckedAt time.Time
	Owner         string
}
