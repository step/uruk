package streamClient

type StreamClient interface {
	Add(sName, key, value string) error
	Read() error
}