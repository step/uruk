package streamClient_test

import (
	"github.com/step/saurontypes"
	"reflect"
	"testing"
	
	s "github.com/step/uruk/pkg/streamClient"
)

func TestAdd(t *testing.T)  {
	client := s.NewDefaultClient()
	entry := saurontypes.Entry{
		Key: "foo",
		Value: "bar",
	}
	entries := make([]saurontypes.Entry,1)
	entries[0] = entry

	err := client.Add("", entries)

	if err != nil {
		t.Errorf("Error happened while adding to stream \n%s", err)
	}

	actual := client.Stream
	expected := s.DefaultStreamClient{&[]saurontypes.Entry{{
		Key: "foo",
		Value: "bar",
	}}}.Stream

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nActual => %s\nExpected => %s\n", actual, expected)
	}
}

