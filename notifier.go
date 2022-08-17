package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/slack-go/slack"
)

type Notifier interface {
	Notify(os []*MaybeOutdated, template *template.Template) error
}

type defaultNotifier struct {
	out io.Writer
}

func (n *defaultNotifier) Notify(os []*MaybeOutdated, tmpl *template.Template) error {
	return tmpl.Execute(n.out, os)
}

type slackNotifier struct {
	client  *slack.Client
	channel string
}

func (n *slackNotifier) Notify(os []*MaybeOutdated, tmpl *template.Template) error {
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, os); err != nil {
		return err
	}
	_, _, err := n.client.PostMessage(n.channel, slack.MsgOptionText(buf.String(), false))
	return err
}

type emailNotifier struct {
}

func (n *emailNotifier) Notify(os []*MaybeOutdated) error {
	return fmt.Errorf("not implemented yet")
}
