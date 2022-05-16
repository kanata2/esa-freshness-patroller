package main

import (
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
)

type Checker interface {
	Check(post string) (*MaybeOutdated, error)
}

type checker struct {
	parser *blackfriday.Markdown
}

func (c *checker) Check(post string) (*MaybeOutdated, error) {
	var (
		mo  MaybeOutdated
		err error
	)
	parser := blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions))
	ast := parser.Parse([]byte(post))
	ast.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if n.Type != blackfriday.Text {
			return blackfriday.GoToNext
		}
		text := string(n.Literal)
		if !strings.HasPrefix(text, "Last checked at") {
			return blackfriday.GoToNext
		}
		words := strings.Split(string(n.Literal), " ")
		// Last checked at YYYY/MM/DD by @username
		// 0    1       2  3          4  5
		// <---------------- 6 ------------------>
		if len(words) < 6 {
			return blackfriday.Terminate
		}
		date, perr := time.Parse("2006/01/02", words[3])
		if perr != nil {
			err = perr
			return blackfriday.Terminate
		}
		// 3 month
		if date.AddDate(0, 3, 0).After(time.Now()) {
			return blackfriday.Terminate
		}
		mo.LastCheckedAt = date
		mo.Owner = strings.TrimSpace(words[5])

		return blackfriday.Terminate
	})
	if err != nil {
		return nil, err
	}
	if mo.Owner == "" {
		return nil, nil
	}
	return &mo, nil
}
