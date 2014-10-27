package parser

import (
	"bytes"
	"strconv"
	"time"
	"unicode"
	"unicode/utf8"
)

type kv struct{ key, val []byte }

func (k kv) String() string {
	return string(k.key) + "=" + string(k.val)
}

func parseLogFmt(data []byte) (*Entry, bool) {
	// don't try to parse logfmt if there's no `mykey=` in the
	// first 100 bytes
	if !startsWithStringEqual(data, 100) {
		return nil, false
	}
	kvs := scanAllKeyValue(data)
	e := newEntry()
	for _, kv := range kvs {
		e.setField(string(kv.key), inferValueField(kv.val))
	}
	return e, true
}

func inferValueField(val []byte) Field {
	if len(val) == 0 {
		return NilField{}
	}

	switch {
	case bytes.Equal(val, []byte("true")):
		return BooleanField(true)
	case bytes.Equal(val, []byte("false")):
		return BooleanField(false)
	case bytes.Equal(val, []byte(`null`)),
		bytes.Equal(val, []byte(`nil`)),
		bytes.Equal(val, []byte(`<nil>`)):
		return NilField{}
	}

	strVal := string(bytes.TrimSpace(val))
	unquotedVal, err := strconv.Unquote(strVal)
	if err != nil {
		f, ok := parseStringTypes(strVal)
		if !ok {
			return StringField(strVal)
		}
		return f
	}

	f, ok := parseStringTypes(unquotedVal)
	if !ok {
		return StringField(unquotedVal)
	}
	return f
}

func parseStringTypes(str string) (Field, bool) {
	f, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return NumberField(f), true
	}

	if t, err := tryParseTime(str); err == nil {
		return TimeField{t}, true
	}
	if d, err := time.ParseDuration(str); err == nil {
		return DurationField{d}, true
	}
	return nil, false
}

func startsWithStringEqual(data []byte, atMost int) bool {
	var i int
	for i < len(data) && i < atMost {
		r, sz := utf8.DecodeRune(data[i:])
		if unicode.IsLetter(r) {
			i += sz
		} else if r == '=' {
			return true
		} else {
			return false
		}
	}
	return false
}

func scanAllKeyValue(data []byte) []kv {
	i := 0
	var kvs []kv
	for i < len(data) {
		keyStart, keyEnd, valStart, valEnd, found := scanKeyValue(data, i)
		if !found {
			return kvs
		}
		kvs = append(kvs, kv{
			key: data[keyStart:keyEnd],
			val: data[valStart:valEnd],
		})
		i = valEnd + 1
	}
	return kvs
}

func scanKeyValue(data []byte, from int) (keyStart, keyEnd, valStart, valEnd int, found bool) {

	keyStart, keyEnd, found = findWordFollowedBy('=', data, from)
	if !found {
		return
	}
	valStart = keyEnd + 1
	if r, sz := utf8.DecodeRune(data[valStart:]); r == '"' {
		// find next unescaped `"`
		valEnd = findUnescaped('"', '\\', data, valStart+sz)
		found = valEnd != -1
		valEnd++
		return
	}

	nextKeyStart, _, nextFound := findWordFollowedBy('=', data, keyEnd+1)

	if nextFound {
		valEnd = nextKeyStart - 1
	} else {
		valEnd = len(data)
	}

	return
}

func findWordFollowedBy(by rune, data []byte, from int) (start int, end int, found bool) {
	i := bytes.IndexRune(data[from:], by)
	if i == -1 {
		return i, i, false
	}
	i += from
	// loop for all letters before the `by`, stop at the first space
	for j := i - 1; j >= from; j-- {
		if !utf8.RuneStart(data[j]) {
			continue
		}
		r, _ := utf8.DecodeRune(data[j:])
		if unicode.IsSpace(r) {
			j++
			return j, i, true //j < i
		}
	}
	return from, i, true //from < i
}

func findUnescaped(toFind, escape rune, data []byte, from int) int {
	for i := from; i < len(data); {
		r, sz := utf8.DecodeRune(data[i:])
		i += sz
		if r == escape {
			// skip next char
			_, sz = utf8.DecodeRune(data[i:])
			i += sz
		} else if r == toFind {
			return i - sz
		}
	}
	return -1
}
