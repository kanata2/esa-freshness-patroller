package patroller

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestChecker_Check(t *testing.T) {
	cases := []struct {
		filename      string
		ExpectedValue *MaybeOutdated
		ExpectedError error
	}{
		{
			"post-1.md",
			&MaybeOutdated{
				Title:         "",
				URL:           "",
				LastCheckedAt: time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC),
				Owners:        []string{"@kanata2"},
			},
			nil,
		},
		{
			"post-2.md",
			&MaybeOutdated{
				Title:         "",
				URL:           "",
				LastCheckedAt: time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC),
				Owners:        []string{"@kanata2"},
			},
			nil},
		{
			"post-3.md",
			&MaybeOutdated{
				Title:         "",
				URL:           "",
				LastCheckedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				Owners:        []string{"@kanata2", "@kanata1"},
			},
			nil},
		// {"post-4.md", nil, parseError},
	}
	checker := &checker{}
	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			b, err := os.ReadFile("testdata/" + tc.filename)
			if err != nil {
				t.Fatal(err)
			}
			actual, err := checker.Check(string(b))
			if d := cmp.Diff(tc.ExpectedValue, actual); d != "" {
				t.Errorf("Returned values are mismatch (-expected +actual):\n%s", d)
			}
		})
	}

}
