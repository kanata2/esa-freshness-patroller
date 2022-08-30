package main

import (
	"context"
	"io"
	"os"

	"github.com/kanata2/esa-freshness-patroller/internal/esa"
)

type Sender interface {
	Send(io.Reader) error
}

type stdoutSender struct{}

func (s *stdoutSender) Send(r io.Reader) error {
	_, err := io.Copy(os.Stdout, r)
	return err
}

type esaSender struct {
	client           *esa.Client
	reportPostNumber int
}

func (s *esaSender) Send(r io.Reader) error {
	bb, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	_, err = s.client.UpdatePost(context.Background(), s.reportPostNumber, esa.WithPostParamsOptionBody(string(bb)))
	return err
}
