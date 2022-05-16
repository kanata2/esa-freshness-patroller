package esa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	Number         int       `json:"number"`
	Name           string    `json:"name"`
	FullName       string    `json:"full_name"`
	WIP            bool      `json:"wip"`
	BodyMarkdown   string    `json:"body_md"`
	BodyHTML       string    `json:"body_html"`
	CreatedAt      time.Time `json:"created_at"`
	Message        string    `json:"message"`
	URL            string    `json:"url"`
	UpdatedAt      time.Time `json:"updated_at"`
	Tags           []string  `json:"tags"`
	Category       string    `json:"category"`
	RevisionNumber int       `json:"revision_number"`
	CreatedBy      struct {
		Myself     bool   `json:"myself"`
		Name       string `json:"name"`
		ScreenName string `json:"screen_name"`
		Icon       string `json:"icon"`
	} `json:"created_by"`
	UpdatedBy struct {
		Myself     bool   `json:"myself"`
		Name       string `json:"name"`
		ScreenName string `json:"screen_name"`
		Icon       string `json:"icon"`
	} `json:"updated_by"`
}

type ListPostsResponse struct {
	Posts      []*Post `json:"posts"`
	PrevPage   *int    `json:"prev_page"`
	NextPage   *int    `json:"next_page"`
	TotalCount int     `json:"total_count"`
	Page       int     `json:"page"`
	PerPage    int     `json:"per_page"`
	MaxPerPage int     `json:"max_per_page"`
}

type ListPostsOption func(url.Values) error

func WithListPostsOptionOrder(o string) ListPostsOption {
	return func(v url.Values) error {
		if o != "desc" || o != "asc" {
			return fmt.Errorf("%s is wrong ordering value. must be specify desc or asc", o)
		}
		v.Add("order", o)
		return nil
	}
}

func WithListPostsOptionSort(s string) ListPostsOption {
	return func(v url.Values) error {
		isSortable := s == "updated" ||
			s == "created" ||
			s == "number" ||
			s == "stars" ||
			s == "watches" ||
			s == "comments" ||
			s == "best_match"
		if !isSortable {
			return fmt.Errorf("%s is not sortable property", s)
		}
		v.Add("sort", s)
		return nil
	}
}

func WithListPostsOptionInclude(fields ...string) ListPostsOption {
	return func(v url.Values) error {
		v.Add("include", strings.Join(fields, ","))
		return nil
	}
}

func WithListPostsOptionQuery(q string) ListPostsOption {
	return func(v url.Values) error {
		v.Add("q", q)
		return nil
	}
}

func WithListPostsOptionPage(n int) ListPostsOption {
	return func(v url.Values) error {
		v.Add("page", strconv.Itoa(n))
		return nil
	}
}

func WithListPostsOptionPerPage(n int) ListPostsOption {
	return func(v url.Values) error {
		v.Add("per_page", strconv.Itoa(n))
		return nil
	}
}

func (c *Client) ListPosts(ctx context.Context, team string, opts ...ListPostsOption) (*ListPostsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path.Join("teams", team, "posts"), nil)
	if err != nil {
		return nil, fmt.Errorf("ListPosts: %w", err)
	}
	q := req.URL.Query()
	for _, opt := range opts {
		if err := opt(q); err != nil {
			return nil, err
		}
	}
	req.URL.RawQuery = q.Encode()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ListPosts: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ListPosts: %s", resp.Status)
	}
	var ret *ListPostsResponse
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("ListPosts: %w", err)
	}
	return ret, nil
}
