package parser

import (
	"bytes"
	"github.com/kr/logfmt"
	"strconv"
	"time"
)

type logfmtEntry struct {
	*Entry
}

func (l *logfmtEntry) HandleLogfmt(key, val []byte) error {
	name := string(key)

	if len(val) == 0 {
		l.setField(name, NilField{})
		return nil
	}

	switch {
	case bytes.Equal(val, []byte("true")):
		l.setField(name, BooleanField{Value: true})
		return nil
	case bytes.Equal(val, []byte("false")):
		l.setField(name, BooleanField{Value: false})
		return nil
	}

	strVal := string(val)
	unquotedVal, err := strconv.Unquote(strVal)
	if err != nil {
		f, ok := parseStringTypes(unquotedVal)
		if !ok {
			f = StringField{Value: unquotedVal}
		}
		l.setField(name, f)
		return nil
	}

	f, ok := parseStringTypes(unquotedVal)
	if !ok {
		f = RawField{Value: val}
	}
	l.setField(name, f)

	return nil
}

func parseStringTypes(str string) (Field, bool) {
	if f, err := strconv.ParseFloat(str, 64); err == nil {
		return NumberField{Value: f}, true
	}

	if t, err := tryParseTime(str); err == nil {
		return TimeField{Value: t}, true
	}
	if d, err := time.ParseDuration(str); err == nil {
		return DurationField{Value: d}, true
	}
	return nil, false
}

func parseLogFmt(data []byte) (*Entry, bool) {
	l := logfmtEntry{newEntry()}
	err := logfmt.Unmarshal(data, &l)
	return l.Entry, err == nil
}
