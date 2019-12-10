package streamClient

import (
	"github.com/step/saurontypes"
)

type StreamClient interface {
	Add(sName string, entries []saurontypes.Entry) error
	Read(streams []string) []saurontypes.StreamEvent
}