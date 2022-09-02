package output

import (
	"context"
	"io"
	"os"

	"github.com/kanata2/esa-freshness-patroller/internal/esa"
)

type Outputter interface {
	Output(io.Reader) error
}

type stdoutOutputter struct{}

func NewStdoutOutputter() *stdoutOutputter {
	return &stdoutOutputter{}
}

func (s *stdoutOutputter) Output(r io.Reader) error {
	_, err := io.Copy(os.Stdout, r)
	return err
}

type esaOutputter struct {
	client           *esa.Client
	reportPostNumber int
}

func NewEsaOutputter(c *esa.Client, reportPostNumber int) *esaOutputter {
	return &esaOutputter{
		client:           c,
		reportPostNumber: reportPostNumber,
	}
}

func (s *esaOutputter) Output(r io.Reader) error {
	bb, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	_, err = s.client.UpdatePost(
		context.Background(),
		s.reportPostNumber,
		esa.WithPostParamsOptionBody(string(bb)),
	)
	return err
}
