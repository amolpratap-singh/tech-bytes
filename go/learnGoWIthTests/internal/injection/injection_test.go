package injection

import (
	"bytes"
	"testing"
)

func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{}
	Greet(&buffer, "Amolpratap")

	got := buffer.String()
	want := "Hello, Amolpratap"

	if got != want {
		t.Errorf("got %q but want %q", got, want)
	}
}