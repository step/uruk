package streamClient_test

import (
	"reflect"
	"testing"
	
	s "github.com/step/uruk/pkg/streamClient"
)

func TestAdd(t *testing.T)  {
	client := s.NewDefaultClient()

	err := client.Add("","foo","bar")

	if err != nil {
		t.Errorf("Error happened while adding to stream \n%s", err)
	}

	actual := client.Stream
	expected := s.DefaultStreamClient{[]s.Entry{{
		"foo", "bar",
	}}}.Stream

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nActual => %s\nExpected => %s\n", actual, expected)
	}
}

