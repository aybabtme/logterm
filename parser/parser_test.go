package parser

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

type logEntry struct {
	Msg  string    `json:"msg"`
	Time time.Time `json:"time"`
}

func TestCanParseJSONLog(t *testing.T) {
	input := logEntry{
		Msg:  "hello",
		Time: time.Date(1900, 01, 01, 01, 01, 01, 01, time.UTC),
	}
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(input)
	if err != nil {
		t.Fatalf("couldn't encode test input: %v", err)
	}
	parser := NewParser(buf)

	for parser.Next() {
		entry := parser.LogEntry()
		msg, ok := entry.Field("msg")
		if !ok {
			t.Fatalf("entry doesn't contain a `msg` field")
		}
		switch m := msg.(type) {
		case StringField:
			if m.Value != input.Msg {
				t.Fatalf("want msg %q, got %q", input.Msg, m.Value)
			}
		case NumberField, TimeField, DurationField, RawField:
			t.Fatalf("want a StringEntry, msg was %T", msg)
		}
	}

	err = parser.Err()
	if err != nil {
		t.Fatalf("got parsing error: %v", err)
	}

}
