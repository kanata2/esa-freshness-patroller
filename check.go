package patroller

import (
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/russross/blackfriday/v2"
)

type Checker interface {
	Check(post string) (*MaybeOutdated, error)
}

type checker struct {
	parser             *blackfriday.Markdown
	threshold          int
	enableSimplyFormat bool
}

type annotation struct {
	Owners            []string          `yaml:"owners"`
	LastCheckedAt     string            `yaml:"last_checked_at"`
	CheckIntervalDays int               `yaml:"check_interval_days"`
	Skip              bool              `yaml:"skip"`
	Custom            map[string]string `yaml:"custom"`
}

func (c *checker) Check(post string) (*MaybeOutdated, error) {
	var (
		mo  *MaybeOutdated
		err error
	)
	// NOTE(kanata2): The esa.io system uses CRLF as a newline. So convert to LF for Blackfriday.
	normalized := strings.NewReplacer("\r\n", "\n").Replace(post)
	parser := blackfriday.New(blackfriday.WithExtensions(blackfriday.FencedCode))
	ast := parser.Parse([]byte(normalized))
	ast.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if !hasEsaFreshnessPatrollerInfoString(n) {
			return blackfriday.GoToNext
		}
		var annotation annotation
		if werr := yaml.Unmarshal(n.Literal, &annotation); werr != nil {
			err = werr
			return blackfriday.Terminate
		}
		var werr error
		mo, werr = c.check(annotation)
		if werr != nil {
			err = werr
		}
		return blackfriday.Terminate

	})
	if err != nil {
		return nil, err
	}
	if mo == nil {
		if !c.enableSimplyFormat {
			return nil, nil
		}
		parser := blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions))
		ast := parser.Parse([]byte(post))
		ast.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
			if !hasEsaFreshnessPatrollerSimpleAnnotation(n) {
				return blackfriday.GoToNext
			}
			var werr error
			a, werr := c.extractSimpleAnnotation(n)
			if werr != nil {
				err = werr
				return blackfriday.Terminate
			}
			mo, werr = c.check(a)
			if werr != nil {
				err = werr
			}
			return blackfriday.Terminate
		})
		if err != nil {
			return nil, err
		}
		return mo, nil
	}
	return mo, nil
}

func (c *checker) check(a annotation) (*MaybeOutdated, error) {
	date, err := parseDate(a.LastCheckedAt)
	if err != nil {
		return nil, err
	}

	if a.CheckIntervalDays == 0 {
		a.CheckIntervalDays = c.threshold
	}
	if a.Skip || date.AddDate(0, 0, a.CheckIntervalDays).After(time.Now()) {
		return nil, nil
	}
	return &MaybeOutdated{Owners: a.Owners, LastCheckedAt: date}, nil
}

func (c *checker) extractSimpleAnnotation(n *blackfriday.Node) (annotation, error) {
	line := strings.Split(string(n.Literal), "\n")[0]
	words := strings.Split(line, " ")
	// Last checked at YYYY/MM/DD by @username1, @username2,...
	// 0    1       2  3          4  5...
	// <-------------- 5 --------->  <------ 1~ -------------->
	if len(words) < 6 {
		return annotation{}, fmt.Errorf("checker.extractSimpleAnnotation: cannot parse simple format against '%s'", string(n.Literal))
	}

	splitOwners := strings.Split(strings.Join(words[5:], ""), ",")
	owners := make([]string, 0, len(splitOwners))
	for i := range splitOwners {
		normalized := strings.TrimSpace(splitOwners[i])
		if normalized == "" {
			continue
		}
		owners = append(owners, normalized)
	}
	return annotation{
		Owners:        owners,
		LastCheckedAt: words[3],
	}, nil
}

func hasEsaFreshnessPatrollerInfoString(n *blackfriday.Node) bool {
	return n.Type == blackfriday.CodeBlock && strings.HasPrefix(string(n.CodeBlockData.Info), "esa-freshness-patroller")
}

func hasEsaFreshnessPatrollerSimpleAnnotation(n *blackfriday.Node) bool {
	return n.Type == blackfriday.Text && strings.HasPrefix(string(n.Literal), "Last checked at")
}

func parseDate(s string) (time.Time, error) {
	date, err := time.Parse("2006/01/02", s)
	if err == nil {
		return date, nil
	}
	date, err = time.Parse("2006-01-02", s)
	if err == nil {
		return date, nil
	}
	return time.Time{}, err
}
