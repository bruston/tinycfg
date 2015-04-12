package tinycfg

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

var happyCase = `server=irc.example.com
name=Example
// Example comment

port=6667
valid=a string containing =
`

var expected = Config{
	vals: map[string]string{
		"server": "irc.example.com",
		"name":   "Example",
		"port":   "6667",
		"valid":  "a string containing =",
	},
}

func TestDecode(t *testing.T) {
	cfg, err := Decode(strings.NewReader(happyCase))
	if err != nil {
		t.Fatalf("expecting nil error, got %s", err)
	}
	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("expecting %#v\n received %#v", expected, cfg)
	}
}

func TestEncode(t *testing.T) {
	expected := `age=29
name=joe
team=gopher
`
	cfg := New()
	cfg.Set("name", "joe")
	cfg.Set("age", "29")
	cfg.Set("team", "gopher")
	var buf bytes.Buffer
	if err := cfg.Encode(&buf); err != nil {
		t.Fatalf("expecting nil error, got %s", err)
	}
	if buf.String() != expected {
		t.Errorf("expected: %s\nreceived: %s", expected, buf.String())
	}
}

// TODO: Refactor and test bad inputs.
