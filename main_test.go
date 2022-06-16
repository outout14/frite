package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestFixuri(t *testing.T) {
	tests := []string{"test/demo", "test/demo/", "/test/demo"}

	for _, uri := range tests {
		want := "/test/demo/"
		msg := fixuri(uri)
		if want != msg {
			t.Fatalf("fixuri('%s') = '%s', wanted '%s'", uri, msg, want)
		}
	}
}

func TestGetlinks(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("test https://example.com")

	a := App{}
	err := a.parselinks(&buffer)

	if err != nil {
		t.Fatalf("parselinks() got error : %s", err)
	}
	for _, l := range a.links {
		want := link{short: "test", to: "https://example.com"}
		if l != want {
			t.Fatalf("Got '%s', wanted '%s'", l, want)
		}
	}
}

func TestGetlinksErr(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("test ")

	a := App{}
	err := a.parselinks(&buffer)

	if !strings.HasPrefix(err.Error(), "invalid syntax") {
		t.Fatalf("parselinks() got unexcepted error : %s", err)
	}
}
