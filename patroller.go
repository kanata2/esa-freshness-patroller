package patroller

import (
	"context"
	"html/template"
	"log"
	"time"

	"github.com/kanata2/esa-freshness-patroller/internal/esa"
)

type patroller struct {
	debug    bool
	client   *esa.Client
	query    string
	checker  *checker
	template *template.Template
	logger   *log.Logger
}

type PatrollerOptionFn func(*patroller)

func WithDebug() PatrollerOptionFn {
	return func(p *patroller) {
		p.debug = true
	}
}

func WithCheckerThreshold(day int) PatrollerOptionFn {
	return func(p *patroller) {
		p.checker.threshold = day
	}
}

func New(esaApiKey, esaTeam, query string, opts ...PatrollerOptionFn) patroller {
	p := patroller{
		client:  esa.NewClient(esaTeam, esaApiKey),
		query:   query,
		checker: &checker{threshold: 90},
		logger:  log.Default(),
	}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

func (p patroller) Debugf(format string, v ...interface{}) {
	if p.debug {
		p.logger.Printf("[debug] "+format, v...)
	}
}

func (p patroller) Infof(format string, v ...interface{}) {
	p.logger.Printf("[info] "+format, v...)
}

func (p patroller) Warnf(format string, v ...interface{}) {
	p.logger.Printf("[warn] "+format, v...)
}

func (p patroller) Patrol(ctx context.Context) (*Result, error) {
	page := 1
	result := &Result{}
	for {
		resp, err := p.client.ListPosts(
			ctx,
			esa.WithListPostsOptionQuery(p.query),
			esa.WithListPostsOptionPage(page),
			esa.WithListPostsOptionPerPage(esa.MaxElementsPerPage),
		)
		if err != nil {
			return nil, err
		}
		p.Debugf("Hit %d posts", len(resp.Posts))
		for _, post := range resp.Posts {
			mo, err := p.checker.Check(post.BodyMarkdown)
			if err != nil {
				result.Warnings = append(result.Warnings, &Warning{
					Title:  post.Name,
					URL:    post.URL,
					Reason: err.Error(),
				})
				continue
			}
			if mo == nil {
				continue
			}
			mo.Title = post.Name
			mo.URL = post.URL
			result.Items = append(result.Items, mo)
		}
		if resp.NextPage == nil {
			break
		}
		page = *resp.NextPage
	}
	return result, nil
}

type Result struct {
	Items    []*MaybeOutdated `json:"items"`
	Warnings []*Warning       `json:"warnings"`
}

type MaybeOutdated struct {
	Title         string    `json:"title"`
	URL           string    `json:"url"`
	LastCheckedAt time.Time `json:"last_checked_at"`
	Owners        []string  `json:"owners"`
}

type Warning struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Reason string `json:"reason"`
}
