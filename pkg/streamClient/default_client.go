package streamClient

import (
	"github.com/step/saurontypes"
)

type DefaultStreamClient struct {
	Stream *[]saurontypes.Entry
}

func (d DefaultStreamClient) Add(sClient string, entries []saurontypes.Entry) error {
	*d.Stream = append(*d.Stream, entries...)
	return nil
}

func (d DefaultStreamClient) Read() error {
	return nil
}

func NewDefaultClient() DefaultStreamClient {
	return DefaultStreamClient{&[]saurontypes.Entry{}}
}
