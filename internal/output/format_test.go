package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewOptionsSelectJSONFormat(t *testing.T) {
	w := New(Options{})
	if w.IsJSON() {
		t.Fatal("default should be text")
	}
	for _, opts := range []Options{
		{JSON: true},
		{JQ: ".foo"},
		{Template: "{{.}}"},
	} {
		if !New(opts).IsJSON() {
			t.Fatalf("%+v should imply JSON", opts)
		}
	}
}

func TestJSONIndented(t *testing.T) {
	var buf bytes.Buffer
	w := &Writer{Out: &buf, Format: FormatJSON}
	if err := w.JSON(map[string]int{"a": 1}); err != nil {
		t.Fatal(err)
	}
	var round map[string]int
	if err := json.Unmarshal(buf.Bytes(), &round); err != nil || round["a"] != 1 {
		t.Fatalf("roundtrip: %v %+v\n%s", err, round, buf.String())
	}
	if !strings.Contains(buf.String(), "  \"a\"") {
		t.Fatal("indent missing")
	}
}

func TestTemplatePipe(t *testing.T) {
	var buf bytes.Buffer
	w := &Writer{Out: &buf, Format: FormatJSON, Template: "{{.name}}"}
	if err := w.JSON(map[string]string{"name": "ok"}); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "ok" {
		t.Fatalf("template out: %q", buf.String())
	}
}

func TestWordWrap(t *testing.T) {
	got := WordWrap("one two three four five", 10)
	if len(got) < 2 {
		t.Fatalf("no wrap: %v", got)
	}
	if short := WordWrap("short", 10); len(short) != 1 || short[0] != "short" {
		t.Fatalf("short: %v", short)
	}
}
