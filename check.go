package patroller

import (
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/russross/blackfriday/v2"
)

type Checker interface {
	Check(post string) (*MaybeOutdated, error)
}

type checker struct {
	parser    *blackfriday.Markdown
	threshold int
}

type annotation struct {
	Owners        []string          `yaml:"owners"`
	LastCheckedAt string            `yaml:"last_checked_at"`
	Interval      int               `yaml:"interval"`
	Skip          bool              `yaml:"skip"`
	Custom        map[string]string `yaml:"custom"`
}

func (c *checker) Check(post string) (*MaybeOutdated, error) {
	var (
		mo  MaybeOutdated
		err error
	)
	parser := blackfriday.New(blackfriday.WithExtensions(blackfriday.FencedCode))
	ast := parser.Parse([]byte(post))
	ast.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if !hasEsaFreshnessPatrollerInfoString(n) {
			return blackfriday.GoToNext
		}
		var annotation annotation
		if werr := yaml.Unmarshal(n.Literal, &annotation); werr != nil {
			err = werr
			return blackfriday.Terminate
		}
		date, werr := time.Parse("2006/01/02", annotation.LastCheckedAt)
		if werr != nil {
			err = werr
			return blackfriday.Terminate
		}

		if annotation.Interval == 0 {
			annotation.Interval = c.threshold
		}
		if annotation.Skip || date.AddDate(0, 0, annotation.Interval).After(time.Now()) {
			return blackfriday.Terminate
		}
		mo.LastCheckedAt = date
		mo.Owners = annotation.Owners

		return blackfriday.Terminate
	})
	if err != nil {
		return nil, err
	}
	if len(mo.Owners) == 0 {
		return nil, nil
	}
	return &mo, nil
}

func hasEsaFreshnessPatrollerInfoString(n *blackfriday.Node) bool {
	return n.Type == blackfriday.CodeBlock && strings.HasPrefix(string(n.CodeBlockData.Info), "esa-freshness-patroller")
}
