package parser

import "time"

type Entry struct {
	names  []string
	fields map[string]Field
}

func newEntry() *Entry {
	return &Entry{
		fields: make(map[string]Field),
	}
}

func (e *Entry) setField(name string, f Field) {
	if _, ok := e.fields[name]; !ok {
		e.names = append(e.names, name)
		e.fields[name] = f
	}
}

func (e *Entry) Field(name string) (Field, bool) {
	f, ok := e.fields[name]
	return f, ok
}

func (e *Entry) FieldNames() []string {
	return e.names
}

type Field interface{}

type NilField struct{}

type UnknownField Field

type RawField []byte

type StringField string

type NumberField float64

type DurationField struct {
	time.Duration
}

type TimeField struct {
	time.Time
}

type BooleanField bool
