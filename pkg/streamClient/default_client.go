package streamClient

type DefaultStreamClient struct {
	Stream []Entry
}

type Entry struct {
	Key   string
	Value string
}

func (d *DefaultStreamClient) Add(sClient, key, value string) error {
	elem := Entry{key, value}
	d.Stream = append(d.Stream, elem)
	return nil
}

func (d *DefaultStreamClient) Read() error {
	return nil
}

func NewDefaultClient() DefaultStreamClient {
	return DefaultStreamClient{[]Entry{}}
}
