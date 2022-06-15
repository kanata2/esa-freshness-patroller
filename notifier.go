package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/slack-go/slack"
)

type Notifier interface {
	Notify(os []*MaybeOutdated) error
}

type defaultNotifier struct {
	out io.Writer
}

func (n *defaultNotifier) Notify(os []*MaybeOutdated) error {
	table := tablewriter.NewWriter(n.out)
	table.SetBorder(false)
	table.SetHeader([]string{"OWNER", "TITLE", "URL", "LAST REVIEWED AT"})
	rows := make([][]string, 0, len(os))
	for _, o := range os {
		rows = append(rows, []string{o.Owner, o.Title, o.URL, o.LastCheckedAt.Format("2006-01-02")})
	}
	table.AppendBulk(rows)
	table.Render()
	return nil
}

type slackNotifier struct {
	client  *slack.Client
	channel string
}

func (n *slackNotifier) Notify(os []*MaybeOutdated) error {
	buf := new(bytes.Buffer)
	sn := defaultNotifier{out: buf}
	if err := sn.Notify(os); err != nil {
		return err
	}
	text := ":closed_book: *esa-freshness-patroller's result* : \n"
	for _, o := range os {
		text += fmt.Sprintf(
			"- <%s|%s> maybe outdated. Last checked at %s by %s\n",
			o.URL, o.Title, o.LastCheckedAt.Format("2006-01-02"), o.Owner,
		)
	}
	_, _, err := n.client.PostMessage(n.channel, slack.MsgOptionText(text, false))
	return err
}

type emailNotifier struct {
}

func (n *emailNotifier) Notify(os []*MaybeOutdated) error {
	return fmt.Errorf("not implemented yet")
}
