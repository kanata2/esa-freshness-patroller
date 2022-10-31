package patroller

import (
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
)

type Checker interface {
	Check(post string) (*MaybeOutdated, error)
}

type checker struct {
	parser         *blackfriday.Markdown
	checkThreshold int
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
		// Last checked at YYYY/MM/DD by @username1, @username2,...
		// 0    1       2  3          4  5...
		// <-------------- 5 --------->  <------ 1~ -------------->
		if len(words) < 6 {
			return blackfriday.Terminate
		}
		date, perr := time.Parse("2006/01/02", words[3])
		if perr != nil {
			err = perr
			return blackfriday.Terminate
		}
		if date.AddDate(0, 0, c.checkThreshold).After(time.Now()) {
			return blackfriday.Terminate
		}
		mo.LastCheckedAt = date
		mo.Owners = normalizeOwners(strings.Join(words[5:], ""))

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

func normalizeOwners(s string) []string {
	splitOwners := strings.Split(s, ",")
	owners := make([]string, 0, len(splitOwners))
	for i := range splitOwners {
		normalized := strings.TrimSpace(splitOwners[i])
		if normalized == "" {
			continue
		}
		owners = append(owners, normalized)
	}
	return owners
}
