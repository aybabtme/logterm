package parser

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
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
			if string(m) != input.Msg {
				t.Fatalf("want msg %q, got %q", input.Msg, m)
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

func TestCanParseTestTable(t *testing.T) {
	for n, tt := range tests {
		t.Logf("test %d", n)
		canParseTestTable(t, tt.input, tt.want)
	}
}

func canParseTestTable(t *testing.T, input string, want map[string]Field) {

	parser := NewParser(strings.NewReader(input))

	if !parser.Next() {
		t.Fatalf("should have one entry to parse, got %v", parser.Err())
	}

	entry := parser.LogEntry()

	checkEntryMatch(t, want, entry)

	if parser.Next() {
		t.Fatalf("should have no more entry to parse, got %v", parser.LogEntry())
	}

	err := parser.Err()
	if err != nil {
		t.Fatalf("got parsing error: %v", err)
	}

}

func checkEntryMatch(t *testing.T, want map[string]Field, entry *Entry) {
	var wantNames []string
	for key := range want {
		wantNames = append(wantNames, key)
	}
	gotNames := entry.FieldNames()
	sort.Strings(wantNames)
	sort.Strings(gotNames)

	if !reflect.DeepEqual(wantNames, gotNames) {
		t.Logf("want=%v", wantNames)
		t.Logf(" got=%v", gotNames)
		t.Fatal("different fields")
	}

	for name, wantField := range want {
		gotField, ok := entry.Field(name)
		if !ok {
			t.Fatalf("entry should have field %q", name)
		}

		if reflect.TypeOf(wantField) != reflect.TypeOf(gotField) {
			t.Logf("want=%v", wantField)
			t.Logf(" got=%v", gotField)
			t.Fatalf("field %q, want type %T, got %T", name, wantField, gotField)
		}

		var match bool
		switch w := wantField.(type) {
		case TimeField:
			g := gotField.(TimeField)
			match = w.Time.Equal(g.Time)

		default:
			match = reflect.DeepEqual(wantField, gotField)
		}

		if !match {
			t.Logf("want=%v", wantField)
			t.Logf(" got=%v", gotField)
			t.Fatalf("field %q's value differs", name)
		}
	}
}

var tests = []struct {
	input string
	want  map[string]Field
}{

	// raw stuff

	{
		want: map[string]Field{
			DefaultRaw: RawField(`[raw log]2014/10/27 18:38:45 warn: map[string]interface {}{"Romano":[]uint8{0x24, 0x5a, 0x9b, 0xd8, 0x97, 0xc1, 0x66, 0x13, 0xec, 0x8f, 0x24, 0x8d, 0x87, 0x53, 0x22, 0x9, 0x3f, 0x5d, 0x42, 0x66, 0xad, 0x2c, 0xd2, 0x56, 0x27, 0x17, 0x2b, 0xa7, 0x9b, 0x6c, 0x96, 0x34, 0xc0, 0x85, 0x2a, 0xb2, 0x80, 0x81, 0xae, 0x17, 0x9e, 0x2e, 0xa, 0x98, 0x39, 0x95, 0x47, 0x5a, 0xc5, 0xeb, 0xf8, 0xf5, 0xa1, 0x67, 0x27, 0x2f, 0xa1, 0x87, 0xa5, 0x42, 0x67, 0x53, 0xe, 0xf, 0x5e, 0xaa, 0x10, 0x8d, 0x91, 0xbb, 0x7b, 0x62, 0x88, 0x7c, 0x2f, 0xf5, 0x2b, 0x2e, 0x10, 0x23, 0x13, 0xb3, 0xbb, 0x6e, 0xaa, 0xf6, 0xa3, 0xc, 0x9c, 0x4d, 0x66, 0x1a, 0x1b, 0x13, 0x90, 0x5f, 0x37, 0x52, 0x59, 0xbc, 0x94, 0xea, 0xbd, 0x33, 0x5a, 0x6a, 0xb8, 0xe7, 0xd3, 0xa5, 0xf7, 0x4e, 0x26, 0x45, 0xfd, 0x1e, 0xe6, 0x81, 0x4b, 0x76, 0x22, 0x28, 0x2c, 0x78, 0xdd, 0xce, 0xfe, 0x2c, 0xc7, 0xb4}, "rem":[]uint8{0x1e, 0x73, 0x5d, 0x93, 0x9f, 0x3d, 0xd2, 0x1a, 0x9a, 0xe8, 0xac, 0x2c, 0xc7, 0x1e, 0xbc, 0x4a, 0xa, 0xe9, 0xb4, 0xea, 0xe2, 0xca, 0xd7, 0x14, 0x3d, 0xe0, 0x10, 0x7f, 0x98, 0x1b, 0x5c, 0x7, 0x27, 0xac, 0xcf, 0x7e, 0x74, 0x45, 0x60, 0x52, 0xfa, 0x5a, 0xad, 0x84, 0xde, 0xee, 0xab, 0x68, 0xe, 0x71, 0x7e, 0x4c, 0xa9, 0x66, 0x23, 0x20, 0xf0, 0x41, 0x8e, 0x35, 0xd2, 0x6e, 0x63, 0x3e, 0xda, 0x4b, 0x57, 0xf, 0x8c, 0x4, 0xa2, 0x70, 0x1a, 0xd9, 0x6d, 0x73, 0x30, 0x35, 0xad, 0x55, 0x7d, 0x41, 0x9f, 0x54, 0x8e, 0x44, 0x7d, 0x7d, 0xab, 0x28, 0x47, 0xd2, 0x4b, 0xb9, 0x99, 0x45, 0x3, 0xbe, 0xaf, 0x26, 0xdd, 0x9a, 0xc4, 0x4a, 0x49, 0xb4, 0x32, 0xdb, 0xd5, 0x0, 0x9e, 0x6c, 0x56, 0x80, 0x88, 0x21, 0x8e, 0x76, 0x13, 0x4a, 0xf8, 0xdb, 0x52, 0x2a, 0x7f, 0x83, 0x12, 0x7b, 0x8f, 0x1b, 0xec, 0xb4, 0x79, 0x8a, 0x18, 0x1a, 0x5c, 0xcf, 0x3f, 0xe3, 0xf0, 0xf9, 0x4, 0x42, 0x96, 0x69, 0x90, 0xf9, 0x95, 0xa4, 0x94, 0x42, 0xac, 0x48, 0xc6, 0x33, 0xbc, 0x71, 0x73, 0x71, 0x52, 0xcf, 0x50, 0x41, 0x1, 0x32, 0xaa, 0x14, 0xe0, 0xf3, 0x5a, 0x1f, 0xf5, 0xab, 0x30, 0x12, 0x19, 0x48, 0x11, 0x34, 0x99, 0x48, 0xd1, 0x1b, 0xb3, 0x9a, 0x26, 0x4c, 0x13, 0xcf, 0x75, 0x1a, 0xd9, 0xbf, 0x29, 0x86, 0xf6, 0xa4, 0xef, 0x24, 0x3c, 0x64, 0x5a, 0x39, 0xcf, 0xb0, 0x7e, 0x47, 0x11, 0xb5, 0xc7, 0xec, 0x2a, 0x24, 0xd3, 0x71, 0xd1, 0x98, 0xaf, 0xa4, 0xd9, 0x6d, 0x40, 0x9, 0x76, 0x4c, 0xb6, 0x82, 0xb}, "si":true, "malivoli":[]uint8{0x24, 0x99, 0xae, 0x4f, 0x37, 0xda, 0x57, 0x74, 0xfc, 0xb1, 0x6d, 0x89, 0x7, 0x6e, 0x46, 0x5c, 0xbe, 0xd3, 0x39, 0xc6, 0xcc, 0x7c, 0xba, 0x44, 0x7e, 0x9f, 0x70, 0x2c, 0x55, 0x67, 0x27, 0x1e, 0x25, 0x43, 0xd9, 0xfb, 0x2c, 0x24, 0xfe, 0x3c, 0x5e, 0xdf, 0x5f, 0x44, 0x58, 0x26, 0xdb, 0x3, 0x1c, 0xa0, 0xfa, 0x8a, 0x7f, 0x86, 0xd2, 0x51, 0x4a, 0xd2, 0xd, 0xfc, 0x1a, 0xbe, 0xad, 0x40, 0xa, 0xe1, 0x49, 0x93, 0x53, 0x24, 0xa, 0x54, 0x3e, 0xc2, 0x55, 0xf9, 0x4d, 0x8b}, "videtur":0, "etiam":[]uint8{0xae, 0xbe, 0xd5, 0x45, 0x1f, 0x49, 0xcf, 0x3c, 0xd, 0x3b, 0x4b, 0x5f, 0xcf, 0xbe, 0xc3, 0x1b, 0xe5, 0x28, 0xb5, 0x8c, 0xe8, 0x2a, 0x5d, 0x9, 0x13, 0x53, 0xfb, 0x48, 0xc, 0x45, 0x93, 0x60, 0x5c, 0x8, 0xec, 0xe5, 0x48, 0xbb, 0xcb, 0x5c, 0xee, 0x89, 0x36, 0xb, 0x42, 0xa8, 0x9f, 0x10, 0x90, 0x5b, 0x92, 0x36, 0xd1, 0x4c, 0xc2, 0x1b, 0xf0, 0x55, 0x23, 0x7, 0xec, 0xda, 0x8a, 0x36, 0x31, 0xe1, 0x91, 0x44, 0xf, 0x79, 0x4e, 0x56, 0xf5, 0x1, 0xc6, 0x28, 0xfa, 0x35, 0x37, 0x2c, 0xad, 0xac, 0x3d, 0x76, 0x78, 0x16, 0xb3, 0x28, 0x3b, 0x48, 0x67, 0x91, 0xf9, 0x31, 0x72, 0x47, 0xca, 0x71, 0x22, 0xd0, 0xa9, 0x3c, 0x5, 0x57, 0xe9, 0x62, 0x2a, 0x73, 0xd5, 0x98, 0x48, 0x6b, 0x90, 0x4d, 0xd8, 0xa6, 0x7a, 0xe0, 0x3c, 0x5d, 0xd3, 0x70, 0x86, 0xd2, 0x48, 0x2a, 0x3d, 0x2, 0x0, 0x3f, 0x7d, 0x63, 0x82, 0x99, 0xda, 0x29, 0xd1, 0x68, 0xa8, 0x8, 0x36, 0x47, 0x7a, 0x1b, 0xc3, 0x61, 0xff, 0x56, 0x1c, 0xa, 0xb4, 0x3a, 0x8a, 0x77, 0xb0, 0x94, 0xd7, 0xa3, 0x8f, 0x7, 0xf2, 0x58, 0x34, 0xdc, 0x4a, 0xf9, 0xb, 0x0, 0xc3, 0x39}}`),
		},
		input: `[raw log]2014/10/27 18:38:45 warn: map[string]interface {}{"Romano":[]uint8{0x24, 0x5a, 0x9b, 0xd8, 0x97, 0xc1, 0x66, 0x13, 0xec, 0x8f, 0x24, 0x8d, 0x87, 0x53, 0x22, 0x9, 0x3f, 0x5d, 0x42, 0x66, 0xad, 0x2c, 0xd2, 0x56, 0x27, 0x17, 0x2b, 0xa7, 0x9b, 0x6c, 0x96, 0x34, 0xc0, 0x85, 0x2a, 0xb2, 0x80, 0x81, 0xae, 0x17, 0x9e, 0x2e, 0xa, 0x98, 0x39, 0x95, 0x47, 0x5a, 0xc5, 0xeb, 0xf8, 0xf5, 0xa1, 0x67, 0x27, 0x2f, 0xa1, 0x87, 0xa5, 0x42, 0x67, 0x53, 0xe, 0xf, 0x5e, 0xaa, 0x10, 0x8d, 0x91, 0xbb, 0x7b, 0x62, 0x88, 0x7c, 0x2f, 0xf5, 0x2b, 0x2e, 0x10, 0x23, 0x13, 0xb3, 0xbb, 0x6e, 0xaa, 0xf6, 0xa3, 0xc, 0x9c, 0x4d, 0x66, 0x1a, 0x1b, 0x13, 0x90, 0x5f, 0x37, 0x52, 0x59, 0xbc, 0x94, 0xea, 0xbd, 0x33, 0x5a, 0x6a, 0xb8, 0xe7, 0xd3, 0xa5, 0xf7, 0x4e, 0x26, 0x45, 0xfd, 0x1e, 0xe6, 0x81, 0x4b, 0x76, 0x22, 0x28, 0x2c, 0x78, 0xdd, 0xce, 0xfe, 0x2c, 0xc7, 0xb4}, "rem":[]uint8{0x1e, 0x73, 0x5d, 0x93, 0x9f, 0x3d, 0xd2, 0x1a, 0x9a, 0xe8, 0xac, 0x2c, 0xc7, 0x1e, 0xbc, 0x4a, 0xa, 0xe9, 0xb4, 0xea, 0xe2, 0xca, 0xd7, 0x14, 0x3d, 0xe0, 0x10, 0x7f, 0x98, 0x1b, 0x5c, 0x7, 0x27, 0xac, 0xcf, 0x7e, 0x74, 0x45, 0x60, 0x52, 0xfa, 0x5a, 0xad, 0x84, 0xde, 0xee, 0xab, 0x68, 0xe, 0x71, 0x7e, 0x4c, 0xa9, 0x66, 0x23, 0x20, 0xf0, 0x41, 0x8e, 0x35, 0xd2, 0x6e, 0x63, 0x3e, 0xda, 0x4b, 0x57, 0xf, 0x8c, 0x4, 0xa2, 0x70, 0x1a, 0xd9, 0x6d, 0x73, 0x30, 0x35, 0xad, 0x55, 0x7d, 0x41, 0x9f, 0x54, 0x8e, 0x44, 0x7d, 0x7d, 0xab, 0x28, 0x47, 0xd2, 0x4b, 0xb9, 0x99, 0x45, 0x3, 0xbe, 0xaf, 0x26, 0xdd, 0x9a, 0xc4, 0x4a, 0x49, 0xb4, 0x32, 0xdb, 0xd5, 0x0, 0x9e, 0x6c, 0x56, 0x80, 0x88, 0x21, 0x8e, 0x76, 0x13, 0x4a, 0xf8, 0xdb, 0x52, 0x2a, 0x7f, 0x83, 0x12, 0x7b, 0x8f, 0x1b, 0xec, 0xb4, 0x79, 0x8a, 0x18, 0x1a, 0x5c, 0xcf, 0x3f, 0xe3, 0xf0, 0xf9, 0x4, 0x42, 0x96, 0x69, 0x90, 0xf9, 0x95, 0xa4, 0x94, 0x42, 0xac, 0x48, 0xc6, 0x33, 0xbc, 0x71, 0x73, 0x71, 0x52, 0xcf, 0x50, 0x41, 0x1, 0x32, 0xaa, 0x14, 0xe0, 0xf3, 0x5a, 0x1f, 0xf5, 0xab, 0x30, 0x12, 0x19, 0x48, 0x11, 0x34, 0x99, 0x48, 0xd1, 0x1b, 0xb3, 0x9a, 0x26, 0x4c, 0x13, 0xcf, 0x75, 0x1a, 0xd9, 0xbf, 0x29, 0x86, 0xf6, 0xa4, 0xef, 0x24, 0x3c, 0x64, 0x5a, 0x39, 0xcf, 0xb0, 0x7e, 0x47, 0x11, 0xb5, 0xc7, 0xec, 0x2a, 0x24, 0xd3, 0x71, 0xd1, 0x98, 0xaf, 0xa4, 0xd9, 0x6d, 0x40, 0x9, 0x76, 0x4c, 0xb6, 0x82, 0xb}, "si":true, "malivoli":[]uint8{0x24, 0x99, 0xae, 0x4f, 0x37, 0xda, 0x57, 0x74, 0xfc, 0xb1, 0x6d, 0x89, 0x7, 0x6e, 0x46, 0x5c, 0xbe, 0xd3, 0x39, 0xc6, 0xcc, 0x7c, 0xba, 0x44, 0x7e, 0x9f, 0x70, 0x2c, 0x55, 0x67, 0x27, 0x1e, 0x25, 0x43, 0xd9, 0xfb, 0x2c, 0x24, 0xfe, 0x3c, 0x5e, 0xdf, 0x5f, 0x44, 0x58, 0x26, 0xdb, 0x3, 0x1c, 0xa0, 0xfa, 0x8a, 0x7f, 0x86, 0xd2, 0x51, 0x4a, 0xd2, 0xd, 0xfc, 0x1a, 0xbe, 0xad, 0x40, 0xa, 0xe1, 0x49, 0x93, 0x53, 0x24, 0xa, 0x54, 0x3e, 0xc2, 0x55, 0xf9, 0x4d, 0x8b}, "videtur":0, "etiam":[]uint8{0xae, 0xbe, 0xd5, 0x45, 0x1f, 0x49, 0xcf, 0x3c, 0xd, 0x3b, 0x4b, 0x5f, 0xcf, 0xbe, 0xc3, 0x1b, 0xe5, 0x28, 0xb5, 0x8c, 0xe8, 0x2a, 0x5d, 0x9, 0x13, 0x53, 0xfb, 0x48, 0xc, 0x45, 0x93, 0x60, 0x5c, 0x8, 0xec, 0xe5, 0x48, 0xbb, 0xcb, 0x5c, 0xee, 0x89, 0x36, 0xb, 0x42, 0xa8, 0x9f, 0x10, 0x90, 0x5b, 0x92, 0x36, 0xd1, 0x4c, 0xc2, 0x1b, 0xf0, 0x55, 0x23, 0x7, 0xec, 0xda, 0x8a, 0x36, 0x31, 0xe1, 0x91, 0x44, 0xf, 0x79, 0x4e, 0x56, 0xf5, 0x1, 0xc6, 0x28, 0xfa, 0x35, 0x37, 0x2c, 0xad, 0xac, 0x3d, 0x76, 0x78, 0x16, 0xb3, 0x28, 0x3b, 0x48, 0x67, 0x91, 0xf9, 0x31, 0x72, 0x47, 0xca, 0x71, 0x22, 0xd0, 0xa9, 0x3c, 0x5, 0x57, 0xe9, 0x62, 0x2a, 0x73, 0xd5, 0x98, 0x48, 0x6b, 0x90, 0x4d, 0xd8, 0xa6, 0x7a, 0xe0, 0x3c, 0x5d, 0xd3, 0x70, 0x86, 0xd2, 0x48, 0x2a, 0x3d, 0x2, 0x0, 0x3f, 0x7d, 0x63, 0x82, 0x99, 0xda, 0x29, 0xd1, 0x68, 0xa8, 0x8, 0x36, 0x47, 0x7a, 0x1b, 0xc3, 0x61, 0xff, 0x56, 0x1c, 0xa, 0xb4, 0x3a, 0x8a, 0x77, 0xb0, 0x94, 0xd7, 0xa3, 0x8f, 0x7, 0xf2, 0x58, 0x34, 0xdc, 0x4a, 0xf9, 0xb, 0x0, 0xc3, 0x39}}`,
	},
	{
		want: map[string]Field{
			DefaultRaw: RawField(`[raw log]2014/10/27 18:38:45 warn: map[string]interface {}{"nemo":false, "alteram":0.8879334884220434, "omnium":time.Time{sec:8724124955, nsec:413, loc:(*time.Location)(0x25acc0)}, "vero":time.Time{sec:39763127421, nsec:927, loc:(*time.Location)(0x25acc0)}, "quidem":0.04749476992250518, "autem":"Quodsi vitam omnem perturbari videmus errore et inscientia, sapientiamque esse solam, quae nos a libidinum impetu et a formidinum terrore vindicet et ipsius fortunae modice ferre doceat iniurias et omnis monstret vias, quae ad quietem et ad tranquillitatem ferant, quid est cur dubitemus dicere et sapientiam propter voluptates expetendam et insipientiam propter molestias esse fugiendam?", "erat":0.3939079663168565, "semper":time.Time{sec:52034192405, nsec:552, loc:(*time.Location)(0x25acc0)}, "abhorreant":true}`),
		},
		input: `[raw log]2014/10/27 18:38:45 warn: map[string]interface {}{"nemo":false, "alteram":0.8879334884220434, "omnium":time.Time{sec:8724124955, nsec:413, loc:(*time.Location)(0x25acc0)}, "vero":time.Time{sec:39763127421, nsec:927, loc:(*time.Location)(0x25acc0)}, "quidem":0.04749476992250518, "autem":"Quodsi vitam omnem perturbari videmus errore et inscientia, sapientiamque esse solam, quae nos a libidinum impetu et a formidinum terrore vindicet et ipsius fortunae modice ferre doceat iniurias et omnis monstret vias, quae ad quietem et ad tranquillitatem ferant, quid est cur dubitemus dicere et sapientiam propter voluptates expetendam et insipientiam propter molestias esse fugiendam?", "erat":0.3939079663168565, "semper":time.Time{sec:52034192405, nsec:552, loc:(*time.Location)(0x25acc0)}, "abhorreant":true}`,
	},
	{
		want: map[string]Field{
			DefaultRaw: RawField(`[raw log]2014/10/27 18:38:45 warn: map[string]interface {}{"praesertim":interface {}(nil)}`),
		},
		input: `[raw log]2014/10/27 18:38:45 warn: map[string]interface {}{"praesertim":interface {}(nil)}`,
	},

	// json stuff
	{
		want: map[string]Field{
			"level": StringField("warning"),
			"mea":   NumberField(0.6645600532184904),
			"msg":   StringField("esse quidem veri viros diligant summum esse ut"),
			"time":  TimeField{time.Date(2014, 10, 27, 18, 31, 40, 0, time.FixedZone("EDT", -4*60*60))},
		},
		input: `{"level":"warning","mea":0.6645600532184904,"msg":"esse quidem veri viros diligant summum esse ut","time":"2014-10-27T18:31:40-04:00"}`,
	},
	{
		want: map[string]Field{
			"diu":        BooleanField(true),
			"enim":       BooleanField(false),
			"est":        NilField{},
			"et":         NumberField(0.6408415373414884),
			"expetendam": BooleanField(false),
			"huius":      BooleanField(false),
			"legimus":    NilField{},
			"level":      StringField("warning"),
			"msg":        StringField("sic stare"),
			"optinere":   NumberField(0.3473303755727788),
			"quo":        StringField("Ita, quae mutat, ea corrumpit, quae sequitur sunt tota Democriti, atomi, inane, imagines, quae eidola nominant, quorum incursione non solum videamus, sed etiam cogitemus; infinitio ipsa, quam apeirian vocant, tota ab illo est, tum innumerabiles munuperatum."),
			"quodsi":     BooleanField(false),
			"rerum":      TimeField{time.Date(1799, 11, 03, 12, 21, 50, 422, time.UTC)},
			"time":       TimeField{time.Date(2014, 10, 27, 18, 38, 45, 0, time.FixedZone("EDT", -4*60*60))},
			"ut":         StringField("At etiam Atheatum: 'Numquidnam manus tua sic affecta, quem ad modum affecta nunc est, desiderat?' -- Nihil sane. -- 'At, si voluptas esset bonum, desideraret.' -- Ita credo. -- 'Non est igitur voluptas bonum.' Hoc ne statuaml dolere, primum tibi recte, Chrysippe, concessum est nihil desiderare manum, cum ita esset affecta, secundum non recte, si voluptas esset bonum, fuisse desideraturam. idcirco enim non desideraret, quia, quod dolore caret, id in voluptate est."),
		},
		input: `{"diu":true,"enim":false,"est":null,"et":0.6408415373414884,"expetendam":false,"huius":false,"legimus":null,"level":"warning","msg":"sic stare","optinere":0.3473303755727788,"quo":"Ita, quae mutat, ea corrumpit, quae sequitur sunt tota Democriti, atomi, inane, imagines, quae eidola nominant, quorum incursione non solum videamus, sed etiam cogitemus; infinitio ipsa, quam apeirian vocant, tota ab illo est, tum innumerabiles munuperatum.","quodsi":false,"rerum":"1799-11-03T12:21:50.000000422Z","time":"2014-10-27T18:38:45-04:00","ut":"At etiam Atheatum: 'Numquidnam manus tua sic affecta, quem ad modum affecta nunc est, desiderat?' -- Nihil sane. -- 'At, si voluptas esset bonum, desideraret.' -- Ita credo. -- 'Non est igitur voluptas bonum.' Hoc ne statuaml dolere, primum tibi recte, Chrysippe, concessum est nihil desiderare manum, cum ita esset affecta, secundum non recte, si voluptas esset bonum, fuisse desideraturam. idcirco enim non desideraret, quia, quod dolore caret, id in voluptate est."}`,
	},
	{
		want: map[string]Field{
			"apud":       NumberField(0.04118484163828559),
			"ardore":     TimeField{time.Date(1486, 3, 5, 8, 23, 42, 525, time.UTC)},
			"conducunt":  StringField("OI18KPtbDwYlJI95ycLGZuJegXmcW84BSfhE85XDS3gGLrDgyxxBaijEbMb4b4VP+9tLARgXVWdB"),
			"dein":       StringField("Quae fuerit causa, mox videro; interea hoc tenebo, si ob aliquam causam ista, quae sine dubio praeclara sunt, fecerint, virtutem iis per se ipsam causam non fuisse. -- Torquem detraxit hosti. -- Et quidem se texit, ne interiret. -- At magnum periculum adiit. -- In oculis quidem exercitus. -- Quid ex eo est consecutus? -- Laudem et caritatem, quae sunt vitae sine metu degendae praesidia firmissima. -- Filium morte multavit. -- Si sine causa, nollem me ab eo ortum, tam inportuno tamque crudeli; sin, ut dolore suo sanciret militaris imperii disciplinam exercitumque in gravissimo bello animadversionis metu contineret, saluti prospexit civium, qua intellegebat contineri suam. atque haec ratio late patet."),
			"dolor":      TimeField{time.Date(915, 02, 28, 18, 35, 8, 194, time.UTC)},
			"enim":       BooleanField(true),
			"esse":       StringField("Nisi mihi Phaedrum, inquam, tu mentitum aut Zenonem putas, quorum utrumque audivi, cum mihi nihil sane praeter sedulitatem probarent, omnes mihi Epicuri sententiae satis notae sunt. atque eos, quos nominavi, cum Attico nostro frequenter audivi, cum miraretur ille quidem utrumque, Phaedrum autem etiam amaret, cotidieque inter nos ea, quae audiebamus, conferebamus, neque erat umquam controversia, quid ego intellegerem, sed quid probarem."),
			"et":         StringField("TEIgdBLCXQk1ysa6Y11aA2ZohpEVG7wGBMFM3VLt7Bt6ausXv1tdCg=="),
			"exultat":    NilField{},
			"hae":        NumberField(0),
			"in":         NumberField(0.5627995383118168),
			"level":      StringField("error"),
			"me":         NilField{},
			"msg":        StringField("nos omnium postulet doctrinis amicitia iis nullas"),
			"nobis":      NilField{},
			"plerisque":  BooleanField(false),
			"quem":       NumberField(0),
			"sola":       StringField("y5xCSn+1yTfcqouq4FVQeC/pKArTKHRlOLqRwC/TeXcs9LY="),
			"time":       TimeField{time.Date(2014, 10, 27, 18, 38, 45, 0, time.FixedZone("EDT", -4*60*60))},
			"volunt":     NilField{},
			"voluptates": StringField("caEehN56R1l11f5Hq8GgVHcOVq5rwWN03DUay+zBBxGw7xsdviPabU+g/bWBAyYHZCqCaKt/GlAH7oj+D28sF/FhWYDORVllIs2uYJbfB0wYViqVH1ZbKwl2al7VEtwi+tI0Fo9X9iP2xSgopBL4Zr4T/V7mxodtroAeFUz3fV0z1gS3rJvWONfBC2Gap3QG16l66rIdTgbWMBzoKdXKW04921tHSJtn/CMJcKXvRD0A02TbReON"),
		},
		input: `{"apud":0.04118484163828559,"ardore":"1486-03-05T08:23:42.000000525Z","conducunt":"OI18KPtbDwYlJI95ycLGZuJegXmcW84BSfhE85XDS3gGLrDgyxxBaijEbMb4b4VP+9tLARgXVWdB","dein":"Quae fuerit causa, mox videro; interea hoc tenebo, si ob aliquam causam ista, quae sine dubio praeclara sunt, fecerint, virtutem iis per se ipsam causam non fuisse. -- Torquem detraxit hosti. -- Et quidem se texit, ne interiret. -- At magnum periculum adiit. -- In oculis quidem exercitus. -- Quid ex eo est consecutus? -- Laudem et caritatem, quae sunt vitae sine metu degendae praesidia firmissima. -- Filium morte multavit. -- Si sine causa, nollem me ab eo ortum, tam inportuno tamque crudeli; sin, ut dolore suo sanciret militaris imperii disciplinam exercitumque in gravissimo bello animadversionis metu contineret, saluti prospexit civium, qua intellegebat contineri suam. atque haec ratio late patet.","dolor":"0915-02-28T18:35:08.000000194Z","enim":true,"esse":"Nisi mihi Phaedrum, inquam, tu mentitum aut Zenonem putas, quorum utrumque audivi, cum mihi nihil sane praeter sedulitatem probarent, omnes mihi Epicuri sententiae satis notae sunt. atque eos, quos nominavi, cum Attico nostro frequenter audivi, cum miraretur ille quidem utrumque, Phaedrum autem etiam amaret, cotidieque inter nos ea, quae audiebamus, conferebamus, neque erat umquam controversia, quid ego intellegerem, sed quid probarem.","et":"TEIgdBLCXQk1ysa6Y11aA2ZohpEVG7wGBMFM3VLt7Bt6ausXv1tdCg==","exultat":null,"hae":0,"in":0.5627995383118168,"level":"error","me":null,"msg":"nos omnium postulet doctrinis amicitia iis nullas","nobis":null,"plerisque":false,"quem":0,"sola":"y5xCSn+1yTfcqouq4FVQeC/pKArTKHRlOLqRwC/TeXcs9LY=","time":"2014-10-27T18:38:45-04:00","volunt":null,"voluptates":"caEehN56R1l11f5Hq8GgVHcOVq5rwWN03DUay+zBBxGw7xsdviPabU+g/bWBAyYHZCqCaKt/GlAH7oj+D28sF/FhWYDORVllIs2uYJbfB0wYViqVH1ZbKwl2al7VEtwi+tI0Fo9X9iP2xSgopBL4Zr4T/V7mxodtroAeFUz3fV0z1gS3rJvWONfBC2Gap3QG16l66rIdTgbWMBzoKdXKW04921tHSJtn/CMJcKXvRD0A02TbReON"}`,
	},
	{
		want: map[string]Field{
			"alias":     TimeField{time.Date(673, 9, 8, 23, 45, 16, 224, time.UTC)},
			"autem":     StringField("Quid? si nos non interpretum fungimur munere, sed tuemur ea, quae dicta sunt ab iis quos probamus, eisque nostrum iudicium et nostrum scribendi ordinem adiungimus, quid habent, cur Graeca anteponant iis, quae et splendide dicta sint neque sint conversa de Graecis? nam si dicent ab illis has res esse tractatas, ne ipsos quidem Graecos est cur tam multos legant, quam legendi sunt. quid enim est a Chrysippo praetermissum in Stoicis? legimus tamen Diogenem, Antipatrum, Mnesarchum, Panaetium, multos alios in primisque familiarem nostrum Posidonium. quid? Theophrastus mediocriterne delectat, cum tractat locos ab Aristotele ante tractatos? quid? Epicurei num desistunt de isdem, de quibus et ab Epicuro scriptum est et ab antiquis, ad arbitrium suum scribere? quodsi Graeci leguntur a Graecis isdem de rebus alia ratione compositis, quid est, cur nostri a nostris non legantur?"),
			"causa":     NilField{},
			"doloribus": TimeField{time.Date(229, 8, 27, 8, 5, 9, 348, time.UTC)},
			"fabulis":   NumberField(0),
			"level":     StringField("info"),
			"msg":       StringField("appareat opes doloris eius conficiuntque intuemur posse"),
			"rebus":     NumberField(0),
			"time":      TimeField{time.Date(2014, 10, 27, 18, 38, 45, 0, time.FixedZone("EDT", -4*60*60))},
		},
		input: `{"alias":"0673-09-08T23:45:16.000000224Z","autem":"Quid? si nos non interpretum fungimur munere, sed tuemur ea, quae dicta sunt ab iis quos probamus, eisque nostrum iudicium et nostrum scribendi ordinem adiungimus, quid habent, cur Graeca anteponant iis, quae et splendide dicta sint neque sint conversa de Graecis? nam si dicent ab illis has res esse tractatas, ne ipsos quidem Graecos est cur tam multos legant, quam legendi sunt. quid enim est a Chrysippo praetermissum in Stoicis? legimus tamen Diogenem, Antipatrum, Mnesarchum, Panaetium, multos alios in primisque familiarem nostrum Posidonium. quid? Theophrastus mediocriterne delectat, cum tractat locos ab Aristotele ante tractatos? quid? Epicurei num desistunt de isdem, de quibus et ab Epicuro scriptum est et ab antiquis, ad arbitrium suum scribere? quodsi Graeci leguntur a Graecis isdem de rebus alia ratione compositis, quid est, cur nostri a nostris non legantur?","causa":null,"doloribus":"0229-08-27T08:05:09.000000348Z","fabulis":0,"level":"info","msg":"appareat opes doloris eius conficiuntque intuemur posse","rebus":0,"time":"2014-10-27T18:38:45-04:00"}`,
	},

	// logfmt stuff

	{
		want: map[string]Field{
			"level":    StringField("warning"),
			"msg":      StringField("naturam reprehenderit dicturam victi a autem et omittam bene"),
			"time":     TimeField{time.Date(2014, 10, 27, 18, 31, 40, 0, time.FixedZone("EDT", -4*60*60))},
			"In":       BooleanField(true),
			"opinemur": StringField("At etiam Athenis, ut e patre audiebam facete et urbane Stoicos irridente, statua est in Ceramico Chrysippi sedentis porrecta manu, quae manus significet illum in hae esse rogatiuncula delectatum: 'Numquidnam manus tua sic affecta, quem ad modum affecta nunc est, desiderat?' -- Nihil sane. -- 'At, si voluptas esset bonum, desideraret.' -- Ita credo. -- 'Non est igitur voluptas bonum.' Hoc ne statuam quidem dicturam pater aiebat, si loqui posset. conclusum est enim contra Cyrenaicos satis acute, nihil ad Epicurum. nam si ea sola voluptas esset, quae quasi titillaret sensus, ut ita dicam, et ad eos cum suavitate afflueret et illaberetur, nec manus esse contenta posset nec ulla pars vacuitate doloris sine iucundo motu voluptatis. sin autem summa voluptas est, ut Epicuro placet, nihil dolere, primum tibi recte, Chrysippe, concessum est nihil desiderare manum, cum ita esset affecta, secundum non recte, si voluptas esset bonum, fuisse desideraturam. idcirco enim non desideraret, quia, quod dolore caret, id in voluptate est."),
		},
		input: `time="2014-10-27T18:31:40-04:00" level="warning" msg="naturam reprehenderit dicturam victi a autem et omittam bene" In=true opinemur="At etiam Athenis, ut e patre audiebam facete et urbane Stoicos irridente, statua est in Ceramico Chrysippi sedentis porrecta manu, quae manus significet illum in hae esse rogatiuncula delectatum: 'Numquidnam manus tua sic affecta, quem ad modum affecta nunc est, desiderat?' -- Nihil sane. -- 'At, si voluptas esset bonum, desideraret.' -- Ita credo. -- 'Non est igitur voluptas bonum.' Hoc ne statuam quidem dicturam pater aiebat, si loqui posset. conclusum est enim contra Cyrenaicos satis acute, nihil ad Epicurum. nam si ea sola voluptas esset, quae quasi titillaret sensus, ut ita dicam, et ad eos cum suavitate afflueret et illaberetur, nec manus esse contenta posset nec ulla pars vacuitate doloris sine iucundo motu voluptatis. sin autem summa voluptas est, ut Epicuro placet, nihil dolere, primum tibi recte, Chrysippe, concessum est nihil desiderare manum, cum ita esset affecta, secundum non recte, si voluptas esset bonum, fuisse desideraturam. idcirco enim non desideraret, quia, quod dolore caret, id in voluptate est."`,
	},
	{
		want: map[string]Field{
			"time":    TimeField{time.Date(2014, 10, 27, 18, 38, 45, 0, time.FixedZone("EDT", -4*60*60))},
			"level":   StringField("warning"),
			"msg":     StringField("faciunt aequum bene quidem et esse parabilis"),
			"atque":   BooleanField(false),
			"et":      StringField("Quid tibi, Torquate, quid huic Triario litterae, quid historiae cognitioque rerum, quid poetarum evolutio, quid tanta tot versuum memoria voluptatis affert? nec mihi illud dixeris: 'Haec enim ipsa mihi sunt voluptati, et erant illa Torquatis.' Numquam hoc ita defendit Epicurus neque Metrodorus aut quisquam eorum, qui aut saperet aliquid aut ista didicisset. et quod quaeritur saepe, cur tam multi sint Epicurei, sunt aliae quoque causae, sed multitudinem haec maxime allicit, quod ita putant dici ab illo, recta et honesta quae sint, ea facere ipsa per se laetitiam, id est voluptatem. homines optimi non intellegunt totam rationem everti, si ita res se habeat. nam si concederetur, etiamsi ad corpus nihil referatur, ista sua sponte et per se esse iucunda, per se esset et virtus et cognitio rerum, quod minime ille vult expetenda."),
			"factis":  NumberField(0),
			"legimus": NumberField(0.6380484281330108),
			"morati":  BooleanField(true),
			"tribuat": TimeField{time.Date(1474, 12, 23, 03, 10, 55, 246, time.UTC)},
			"tum":     NumberField(0.36542830164796136),
		},
		input: `time="2014-10-27T18:38:45-04:00" level="warning" msg="faciunt aequum bene quidem et esse parabilis" atque=false et="Quid tibi, Torquate, quid huic Triario litterae, quid historiae cognitioque rerum, quid poetarum evolutio, quid tanta tot versuum memoria voluptatis affert? nec mihi illud dixeris: 'Haec enim ipsa mihi sunt voluptati, et erant illa Torquatis.' Numquam hoc ita defendit Epicurus neque Metrodorus aut quisquam eorum, qui aut saperet aliquid aut ista didicisset. et quod quaeritur saepe, cur tam multi sint Epicurei, sunt aliae quoque causae, sed multitudinem haec maxime allicit, quod ita putant dici ab illo, recta et honesta quae sint, ea facere ipsa per se laetitiam, id est voluptatem. homines optimi non intellegunt totam rationem everti, si ita res se habeat. nam si concederetur, etiamsi ad corpus nihil referatur, ista sua sponte et per se esse iucunda, per se esset et virtus et cognitio rerum, quod minime ille vult expetenda." factis=0 legimus=0.6380484281330108 morati=true tribuat=1474-12-23 03:10:55.000000246 +0000 UTC tum=0.36542830164796136 `,
	},
	{
		want: map[string]Field{
			"time":      TimeField{time.Date(2014, 10, 27, 18, 38, 45, 0, time.FixedZone("EDT", -4*60*60))},
			"level":     StringField("error"),
			"msg":       StringField("qua"),
			"dixi":      StringField("Ego autem quem timeam lectorem, cum ad te ne Graecis quidem cedentem in philosophia audeam scribere? quamquam a te ipso id quidem facio provocatus gratissimo mihi libro, quem ad me de virtute misisti. Sed ex eo credo quibusdam usu venire; ut abhorreant a Latinis, quod inciderint in inculta quaedam et horrida, de malis Graecis Latine scripta deterius. quibus ego assentior, dum modo de isdem rebus ne Graecos quidem legendos putent. res vero bonas verbis electis graviter ornateque dictas quis non legat? nisi qui se plane Graecum dici velit, ut a Scaevola est praetore salutatus Athenis Albucius."),
			"effecerit": StringField("Sed ut iis bonis erigimur, quae expectamus, sic laetamur iis, quae recordamur. stulti autem malorum memoria torquentur, sapientes bona praeterita grata recordatione renovata delectant. est autem situm in nobis ut et adversa quasi perpetua oblivione obruamus et secunda iucunde ac suaviter meminerimus. sed cum ea, quae praeterierunt, acri animo et attento intuemur, tum fit ut aegritudo sequatur, si illa mala sint, laetitia, si bona."),
			"nec":       NilField{},
			"ob":        StringField("Ut autem a facillimis ordiamur, prima veniat in medium Epicuri ratio, quae plerisque notissima est. quam a nobis sic intelleges eitam, ut ab ipsis, qui eam disciplinam probant, non soleat accuratius explicari; verum enim invenire volumus, non tamquam adversarium aliquem convincere. accurate autem quondam a L. Torquato, homine omni doctrina erudito, defensa est Epicuri sententia de voluptate, a meque ei responsum, cum C. Triarius, in primis gravis et doctus adolescens, ei disputationi interesset."),
			"quidem":    BooleanField(false),
			"quod":      NumberField(0),
			"se":        BooleanField(true),
			"sequitur":  NumberField(0.13815772053664252),
			"summum":    NumberField(0.9507713169389637),
			"tamen":     StringField("Hanc ego cum teneam sententiam, quid est cur verear, ne ad eam non possim accommodare Torquatos nostros? quos tu paulo ante cum memoriter, tum etiam erga nos amice et benivole collegisti, nec me tamen laudandis maioribus meis corrupisti nec segniorem ad respondendum reddidisti. quorum facta quem ad modum, quaeso, interpretaris? sicine eos censes aut in armatum hostem impetum fecisse aut in liberos atque in sanguinem suum tam crudelis fuisse, nihil ut de utilitatibus, nihil ut de commodis suis cogitarent? at id ne ferae quidem faciunt, ut ita ruant itaque turbent, ut earum motus et impetus quo pertineant non intellegamus, tu tam egregios viros censes tantas res gessisse sine causa?"),
			"ut":        NumberField(0.9220774561265302),
		},
		input: `time="2014-10-27T18:38:45-04:00" level="error" msg="qua" dixi="Ego autem quem timeam lectorem, cum ad te ne Graecis quidem cedentem in philosophia audeam scribere? quamquam a te ipso id quidem facio provocatus gratissimo mihi libro, quem ad me de virtute misisti. Sed ex eo credo quibusdam usu venire; ut abhorreant a Latinis, quod inciderint in inculta quaedam et horrida, de malis Graecis Latine scripta deterius. quibus ego assentior, dum modo de isdem rebus ne Graecos quidem legendos putent. res vero bonas verbis electis graviter ornateque dictas quis non legat? nisi qui se plane Graecum dici velit, ut a Scaevola est praetore salutatus Athenis Albucius." effecerit="Sed ut iis bonis erigimur, quae expectamus, sic laetamur iis, quae recordamur. stulti autem malorum memoria torquentur, sapientes bona praeterita grata recordatione renovata delectant. est autem situm in nobis ut et adversa quasi perpetua oblivione obruamus et secunda iucunde ac suaviter meminerimus. sed cum ea, quae praeterierunt, acri animo et attento intuemur, tum fit ut aegritudo sequatur, si illa mala sint, laetitia, si bona." nec=<nil> ob="Ut autem a facillimis ordiamur, prima veniat in medium Epicuri ratio, quae plerisque notissima est. quam a nobis sic intelleges eitam, ut ab ipsis, qui eam disciplinam probant, non soleat accuratius explicari; verum enim invenire volumus, non tamquam adversarium aliquem convincere. accurate autem quondam a L. Torquato, homine omni doctrina erudito, defensa est Epicuri sententia de voluptate, a meque ei responsum, cum C. Triarius, in primis gravis et doctus adolescens, ei disputationi interesset." quidem=false quod=0 se=true sequitur=0.13815772053664252 summum=0.9507713169389637 tamen="Hanc ego cum teneam sententiam, quid est cur verear, ne ad eam non possim accommodare Torquatos nostros? quos tu paulo ante cum memoriter, tum etiam erga nos amice et benivole collegisti, nec me tamen laudandis maioribus meis corrupisti nec segniorem ad respondendum reddidisti. quorum facta quem ad modum, quaeso, interpretaris? sicine eos censes aut in armatum hostem impetum fecisse aut in liberos atque in sanguinem suum tam crudelis fuisse, nihil ut de utilitatibus, nihil ut de commodis suis cogitarent? at id ne ferae quidem faciunt, ut ita ruant itaque turbent, ut earum motus et impetus quo pertineant non intellegamus, tu tam egregios viros censes tantas res gessisse sine causa?" ut=0.9220774561265302 `,
	},
	{
		want: map[string]Field{

			"time":       TimeField{time.Date(2014, 10, 27, 18, 38, 45, 0, time.FixedZone("EDT", -4*60*60))},
			"level":      StringField("error"),
			"msg":        StringField("paulo esse sibi ab ratio saepe"),
			"":           NumberField(0),
			"Graeci":     NumberField(0),
			"arare":      StringField("Quam ob rem tandem, inquit, non satisfacit? te enim iudicem aequum puto, modo quae dicat ille bene noris."),
			"enim":       NilField{},
			"facere":     BooleanField(false),
			"in":         StringField("Id qui in una virtute ponunt et splendore nominis capti quid natura postulet non intellegunt, errore maximo, si Epicurum audire voluerint, liberabuntur: istae enim vestrae eximiae pulchraeque virtutes nisi voluptatem efficerent, quis eas aut laudabilis aut expetendas arbitraretur? ut enim medicorum scientiam non ipsius artis, sed bonae valetudinis causa probamus, et gubernatoris ars, quia bene navigandi rationem habet, utilitate, non arte laudatur, sic sapientia, quae ars vivendi putanda est, non expeteretur, si nihil efficeret; nunc expetitur, quod est tamquam artifex conquirendae et comparandae voluptatis --"),
			"inpendente": StringField("Quid tibi, Torquate, quid huic Triario litterae, quid historiae cognitioque rerum, quid poetarum evolutio, quid tanta tot versuum memoria voluptatis affert? nec mihi illud dixeris: 'Haec enim ipsa mihi sunt voluptati, et erant illa Torquatis.' Numquam hoc ita defendit Epicurus neque Metrodorus aut quisquam eorum, qui aut saperet aliquid aut ista didicisset. et quod quaeritur saepe, cur tam multi sint Epicurei, sunt aliae quoque causae, sed multitudinem haec maxime allicit, quod ita putant dici ab illo, recta et honesta quae sint, ea facere ipsa per se laetitiam, id est voluptatem. homines optimi non intellegunt totam rationem everti, si ita res se habeat. nam si concederetur, etiamsi ad corpus nihil referatur, ista sua sponte et per se esse iucunda, per se esset et virtus et cognitio rerum, quod minime ille vult expetenda."),
			"inviti":     NumberField(0.5559171606675798),
			"ipsam":      NumberField(0.20452342141212618),
			"iudicia":    NumberField(0.3127147602988959),
			"non":        StringField("Quam ob rem tandem, inquit, non satisfacit? te enim iudicem aequum puto, modo quae dicat ille bene noris."),
			"quidem":     BooleanField(true),
			"se":         NilField{},
			"sed":        StringField("Iam in altera philosophiae parte. quae est quaerendi ac disserendi, quae logikh dicitur, iste vester plane, ut mihi quidem videtur, inermis ac nudus est. tollit definitiones, nihil de dividendo ac partiendo docet, non quo modo efficiatur concludaturque ratio tradit, non qua via captiosa solvantur ambigua distinguantur ostendit; iudicia rerum in sensibus ponit, quibus si semel aliquid falsi pro vero probatum sit, sublatum esse omne iudicium veri et falsi putat."),
			"sustulisti": TimeField{time.Date(1127, 5, 14, 1, 23, 2, 560, time.UTC)},
			"tenere":     NumberField(0),
			"verbis":     NilField{},
			"voluptatis": StringField("Quamquam, si plane sic verterem Platonem aut Aristotelem, ut verterunt nostri poetae fabulas, male, credo, mererer de meis civibus, si ad eorum cognitionem divina illa ingenia transferrem. sed id neque feci adhuc nec mihi tamen, ne faciam, interdictum puto. locos quidem quosdam, si videbitur, transferam, et maxime ab iis, quos modo nominavi, cum inciderit, ut id apte fieri possit, ut ab Homero Ennius, Afranius a Menandro solet. Nec vero, ut noster Lucilius, recusabo, quo minus omnes mea legant. utinam esset ille Persius, Scipio vero et Rutilius multo etiam magis, quorum ille iudicium reformidans Tarentinis ait se et Consentinis et Siculis scribere. facete is quidem, sicut alia; sed neque tam docti tum erant, ad quorum iudicium elaboraret, et sunt illius scripta leviora, ut urbanitas summa appareat, doctrina mediocris."),
		},
		input: `time="2014-10-27T18:38:45-04:00" level="error" msg="paulo esse sibi ab ratio saepe" =0 Graeci=0 arare="Quam ob rem tandem, inquit, non satisfacit? te enim iudicem aequum puto, modo quae dicat ille bene noris." enim=<nil> facere=false in="Id qui in una virtute ponunt et splendore nominis capti quid natura postulet non intellegunt, errore maximo, si Epicurum audire voluerint, liberabuntur: istae enim vestrae eximiae pulchraeque virtutes nisi voluptatem efficerent, quis eas aut laudabilis aut expetendas arbitraretur? ut enim medicorum scientiam non ipsius artis, sed bonae valetudinis causa probamus, et gubernatoris ars, quia bene navigandi rationem habet, utilitate, non arte laudatur, sic sapientia, quae ars vivendi putanda est, non expeteretur, si nihil efficeret; nunc expetitur, quod est tamquam artifex conquirendae et comparandae voluptatis --" inpendente="Quid tibi, Torquate, quid huic Triario litterae, quid historiae cognitioque rerum, quid poetarum evolutio, quid tanta tot versuum memoria voluptatis affert? nec mihi illud dixeris: 'Haec enim ipsa mihi sunt voluptati, et erant illa Torquatis.' Numquam hoc ita defendit Epicurus neque Metrodorus aut quisquam eorum, qui aut saperet aliquid aut ista didicisset. et quod quaeritur saepe, cur tam multi sint Epicurei, sunt aliae quoque causae, sed multitudinem haec maxime allicit, quod ita putant dici ab illo, recta et honesta quae sint, ea facere ipsa per se laetitiam, id est voluptatem. homines optimi non intellegunt totam rationem everti, si ita res se habeat. nam si concederetur, etiamsi ad corpus nihil referatur, ista sua sponte et per se esse iucunda, per se esset et virtus et cognitio rerum, quod minime ille vult expetenda." inviti=0.5559171606675798 ipsam=0.20452342141212618 iudicia=0.3127147602988959 non="Quam ob rem tandem, inquit, non satisfacit? te enim iudicem aequum puto, modo quae dicat ille bene noris." quidem=true se=<nil> sed="Iam in altera philosophiae parte. quae est quaerendi ac disserendi, quae logikh dicitur, iste vester plane, ut mihi quidem videtur, inermis ac nudus est. tollit definitiones, nihil de dividendo ac partiendo docet, non quo modo efficiatur concludaturque ratio tradit, non qua via captiosa solvantur ambigua distinguantur ostendit; iudicia rerum in sensibus ponit, quibus si semel aliquid falsi pro vero probatum sit, sublatum esse omne iudicium veri et falsi putat." sustulisti=1127-05-14 01:23:02.00000056 +0000 UTC tenere=0 verbis=<nil> voluptatis="Quamquam, si plane sic verterem Platonem aut Aristotelem, ut verterunt nostri poetae fabulas, male, credo, mererer de meis civibus, si ad eorum cognitionem divina illa ingenia transferrem. sed id neque feci adhuc nec mihi tamen, ne faciam, interdictum puto. locos quidem quosdam, si videbitur, transferam, et maxime ab iis, quos modo nominavi, cum inciderit, ut id apte fieri possit, ut ab Homero Ennius, Afranius a Menandro solet. Nec vero, ut noster Lucilius, recusabo, quo minus omnes mea legant. utinam esset ille Persius, Scipio vero et Rutilius multo etiam magis, quorum ille iudicium reformidans Tarentinis ait se et Consentinis et Siculis scribere. facete is quidem, sicut alia; sed neque tam docti tum erant, ad quorum iudicium elaboraret, et sunt illius scripta leviora, ut urbanitas summa appareat, doctrina mediocris." `,
	},
	{
		want: map[string]Field{
			"time":      TimeField{time.Date(2014, 10, 27, 18, 38, 45, 0, time.FixedZone("EDT", -4*60*60))},
			"level":     StringField("warning"),
			"msg":       StringField("est innumerabiles ab nobis"),
			"inesse":    StringField("Quae enim cupiditates a natura proficiscuntur, facile explentur sine ulla iniuria, quae autem inanes sunt, iis parendum non est. nihil enim desiderabile concupiscunt, plusque in ipsa iniuria detrimenti est quam in iis rebus emolumenti, quae pariuntur iniuria. Itaque ne iustitiam quidem recte quis dixerit per se ipsam optabilem, sed quia iucunditatis vel plurimum afferat. nam diligi et carum esse iucundum est propterea, quia tutiorem vitam et voluptatem pleniorem efficit. itaque non ob ea solum incommoda, quae eveniunt inprobis, fugiendam inprobitatem putamus, sed multo etiam magis, quod, cuius in animo versatur, numquam sinit eum respirare, numquam adquiescere."),
			"non":       BooleanField(true),
			"praeclare": NumberField(0.5869776997848529),
		},
		input: `time="2014-10-27T18:38:45-04:00" level="warning" msg="est innumerabiles ab nobis" inesse="Quae enim cupiditates a natura proficiscuntur, facile explentur sine ulla iniuria, quae autem inanes sunt, iis parendum non est. nihil enim desiderabile concupiscunt, plusque in ipsa iniuria detrimenti est quam in iis rebus emolumenti, quae pariuntur iniuria. Itaque ne iustitiam quidem recte quis dixerit per se ipsam optabilem, sed quia iucunditatis vel plurimum afferat. nam diligi et carum esse iucundum est propterea, quia tutiorem vitam et voluptatem pleniorem efficit. itaque non ob ea solum incommoda, quae eveniunt inprobis, fugiendam inprobitatem putamus, sed multo etiam magis, quod, cuius in animo versatur, numquam sinit eum respirare, numquam adquiescere." non=true praeclare=0.5869776997848529 `,
	},
}

func TestCanParseFewLines(t *testing.T) {

	var input = `[raw log]2014/10/27 18:38:45 warn: map[string]interface {}{"Romano":[]uint8{0x24, 0x5a, 0x9b, 0xd8, 0x97, 0xc1, 0x66, 0x13, 0xec, 0x8f, 0x24, 0x8d, 0x87, 0x53, 0x22, 0x9, 0x3f, 0x5d, 0x42, 0x66, 0xad, 0x2c, 0xd2, 0x56, 0x27, 0x17, 0x2b, 0xa7, 0x9b, 0x6c, 0x96, 0x34, 0xc0, 0x85, 0x2a, 0xb2, 0x80, 0x81, 0xae, 0x17, 0x9e, 0x2e, 0xa, 0x98, 0x39, 0x95, 0x47, 0x5a, 0xc5, 0xeb, 0xf8, 0xf5, 0xa1, 0x67, 0x27, 0x2f, 0xa1, 0x87, 0xa5, 0x42, 0x67, 0x53, 0xe, 0xf, 0x5e, 0xaa, 0x10, 0x8d, 0x91, 0xbb, 0x7b, 0x62, 0x88, 0x7c, 0x2f, 0xf5, 0x2b, 0x2e, 0x10, 0x23, 0x13, 0xb3, 0xbb, 0x6e, 0xaa, 0xf6, 0xa3, 0xc, 0x9c, 0x4d, 0x66, 0x1a, 0x1b, 0x13, 0x90, 0x5f, 0x37, 0x52, 0x59, 0xbc, 0x94, 0xea, 0xbd, 0x33, 0x5a, 0x6a, 0xb8, 0xe7, 0xd3, 0xa5, 0xf7, 0x4e, 0x26, 0x45, 0xfd, 0x1e, 0xe6, 0x81, 0x4b, 0x76, 0x22, 0x28, 0x2c, 0x78, 0xdd, 0xce, 0xfe, 0x2c, 0xc7, 0xb4}, "rem":[]uint8{0x1e, 0x73, 0x5d, 0x93, 0x9f, 0x3d, 0xd2, 0x1a, 0x9a, 0xe8, 0xac, 0x2c, 0xc7, 0x1e, 0xbc, 0x4a, 0xa, 0xe9, 0xb4, 0xea, 0xe2, 0xca, 0xd7, 0x14, 0x3d, 0xe0, 0x10, 0x7f, 0x98, 0x1b, 0x5c, 0x7, 0x27, 0xac, 0xcf, 0x7e, 0x74, 0x45, 0x60, 0x52, 0xfa, 0x5a, 0xad, 0x84, 0xde, 0xee, 0xab, 0x68, 0xe, 0x71, 0x7e, 0x4c, 0xa9, 0x66, 0x23, 0x20, 0xf0, 0x41, 0x8e, 0x35, 0xd2, 0x6e, 0x63, 0x3e, 0xda, 0x4b, 0x57, 0xf, 0x8c, 0x4, 0xa2, 0x70, 0x1a, 0xd9, 0x6d, 0x73, 0x30, 0x35, 0xad, 0x55, 0x7d, 0x41, 0x9f, 0x54, 0x8e, 0x44, 0x7d, 0x7d, 0xab, 0x28, 0x47, 0xd2, 0x4b, 0xb9, 0x99, 0x45, 0x3, 0xbe, 0xaf, 0x26, 0xdd, 0x9a, 0xc4, 0x4a, 0x49, 0xb4, 0x32, 0xdb, 0xd5, 0x0, 0x9e, 0x6c, 0x56, 0x80, 0x88, 0x21, 0x8e, 0x76, 0x13, 0x4a, 0xf8, 0xdb, 0x52, 0x2a, 0x7f, 0x83, 0x12, 0x7b, 0x8f, 0x1b, 0xec, 0xb4, 0x79, 0x8a, 0x18, 0x1a, 0x5c, 0xcf, 0x3f, 0xe3, 0xf0, 0xf9, 0x4, 0x42, 0x96, 0x69, 0x90, 0xf9, 0x95, 0xa4, 0x94, 0x42, 0xac, 0x48, 0xc6, 0x33, 0xbc, 0x71, 0x73, 0x71, 0x52, 0xcf, 0x50, 0x41, 0x1, 0x32, 0xaa, 0x14, 0xe0, 0xf3, 0x5a, 0x1f, 0xf5, 0xab, 0x30, 0x12, 0x19, 0x48, 0x11, 0x34, 0x99, 0x48, 0xd1, 0x1b, 0xb3, 0x9a, 0x26, 0x4c, 0x13, 0xcf, 0x75, 0x1a, 0xd9, 0xbf, 0x29, 0x86, 0xf6, 0xa4, 0xef, 0x24, 0x3c, 0x64, 0x5a, 0x39, 0xcf, 0xb0, 0x7e, 0x47, 0x11, 0xb5, 0xc7, 0xec, 0x2a, 0x24, 0xd3, 0x71, 0xd1, 0x98, 0xaf, 0xa4, 0xd9, 0x6d, 0x40, 0x9, 0x76, 0x4c, 0xb6, 0x82, 0xb}, "si":true, "malivoli":[]uint8{0x24, 0x99, 0xae, 0x4f, 0x37, 0xda, 0x57, 0x74, 0xfc, 0xb1, 0x6d, 0x89, 0x7, 0x6e, 0x46, 0x5c, 0xbe, 0xd3, 0x39, 0xc6, 0xcc, 0x7c, 0xba, 0x44, 0x7e, 0x9f, 0x70, 0x2c, 0x55, 0x67, 0x27, 0x1e, 0x25, 0x43, 0xd9, 0xfb, 0x2c, 0x24, 0xfe, 0x3c, 0x5e, 0xdf, 0x5f, 0x44, 0x58, 0x26, 0xdb, 0x3, 0x1c, 0xa0, 0xfa, 0x8a, 0x7f, 0x86, 0xd2, 0x51, 0x4a, 0xd2, 0xd, 0xfc, 0x1a, 0xbe, 0xad, 0x40, 0xa, 0xe1, 0x49, 0x93, 0x53, 0x24, 0xa, 0x54, 0x3e, 0xc2, 0x55, 0xf9, 0x4d, 0x8b}, "videtur":0, "etiam":[]uint8{0xae, 0xbe, 0xd5, 0x45, 0x1f, 0x49, 0xcf, 0x3c, 0xd, 0x3b, 0x4b, 0x5f, 0xcf, 0xbe, 0xc3, 0x1b, 0xe5, 0x28, 0xb5, 0x8c, 0xe8, 0x2a, 0x5d, 0x9, 0x13, 0x53, 0xfb, 0x48, 0xc, 0x45, 0x93, 0x60, 0x5c, 0x8, 0xec, 0xe5, 0x48, 0xbb, 0xcb, 0x5c, 0xee, 0x89, 0x36, 0xb, 0x42, 0xa8, 0x9f, 0x10, 0x90, 0x5b, 0x92, 0x36, 0xd1, 0x4c, 0xc2, 0x1b, 0xf0, 0x55, 0x23, 0x7, 0xec, 0xda, 0x8a, 0x36, 0x31, 0xe1, 0x91, 0x44, 0xf, 0x79, 0x4e, 0x56, 0xf5, 0x1, 0xc6, 0x28, 0xfa, 0x35, 0x37, 0x2c, 0xad, 0xac, 0x3d, 0x76, 0x78, 0x16, 0xb3, 0x28, 0x3b, 0x48, 0x67, 0x91, 0xf9, 0x31, 0x72, 0x47, 0xca, 0x71, 0x22, 0xd0, 0xa9, 0x3c, 0x5, 0x57, 0xe9, 0x62, 0x2a, 0x73, 0xd5, 0x98, 0x48, 0x6b, 0x90, 0x4d, 0xd8, 0xa6, 0x7a, 0xe0, 0x3c, 0x5d, 0xd3, 0x70, 0x86, 0xd2, 0x48, 0x2a, 0x3d, 0x2, 0x0, 0x3f, 0x7d, 0x63, 0x82, 0x99, 0xda, 0x29, 0xd1, 0x68, 0xa8, 0x8, 0x36, 0x47, 0x7a, 0x1b, 0xc3, 0x61, 0xff, 0x56, 0x1c, 0xa, 0xb4, 0x3a, 0x8a, 0x77, 0xb0, 0x94, 0xd7, 0xa3, 0x8f, 0x7, 0xf2, 0x58, 0x34, 0xdc, 0x4a, 0xf9, 0xb, 0x0, 0xc3, 0x39}}
[raw log]2014/10/27 18:38:45 warn: map[string]interface {}{"nemo":false, "alteram":0.8879334884220434, "omnium":time.Time{sec:8724124955, nsec:413, loc:(*time.Location)(0x25acc0)}, "vero":time.Time{sec:39763127421, nsec:927, loc:(*time.Location)(0x25acc0)}, "quidem":0.04749476992250518, "autem":"Quodsi vitam omnem perturbari videmus errore et inscientia, sapientiamque esse solam, quae nos a libidinum impetu et a formidinum terrore vindicet et ipsius fortunae modice ferre doceat iniurias et omnis monstret vias, quae ad quietem et ad tranquillitatem ferant, quid est cur dubitemus dicere et sapientiam propter voluptates expetendam et insipientiam propter molestias esse fugiendam?", "erat":0.3939079663168565, "semper":time.Time{sec:52034192405, nsec:552, loc:(*time.Location)(0x25acc0)}, "abhorreant":true}
[raw log]2014/10/27 18:38:45 warn: map[string]interface {}{"praesertim":interface {}(nil)}
{"level":"warning","mea":0.6645600532184904,"msg":"esse quidem veri viros diligant summum esse ut","time":"2014-10-27T18:31:40-04:00"}
{"diu":true,"enim":false,"est":null,"et":0.6408415373414884,"expetendam":false,"huius":false,"legimus":null,"level":"warning","msg":"sic stare","optinere":0.3473303755727788,"quo":"Ita, quae mutat, ea corrumpit, quae sequitur sunt tota Democriti, atomi, inane, imagines, quae eidola nominant, quorum incursione non solum videamus, sed etiam cogitemus; infinitio ipsa, quam apeirian vocant, tota ab illo est, tum innumerabiles munuperatum.","quodsi":false,"rerum":"1799-11-03T12:21:50.000000422Z","time":"2014-10-27T18:38:45-04:00","ut":"At etiam Atheatum: 'Numquidnam manus tua sic affecta, quem ad modum affecta nunc est, desiderat?' -- Nihil sane. -- 'At, si voluptas esset bonum, desideraret.' -- Ita credo. -- 'Non est igitur voluptas bonum.' Hoc ne statuaml dolere, primum tibi recte, Chrysippe, concessum est nihil desiderare manum, cum ita esset affecta, secundum non recte, si voluptas esset bonum, fuisse desideraturam. idcirco enim non desideraret, quia, quod dolore caret, id in voluptate est."}
{"apud":0.04118484163828559,"ardore":"1486-03-05T08:23:42.000000525Z","conducunt":"OI18KPtbDwYlJI95ycLGZuJegXmcW84BSfhE85XDS3gGLrDgyxxBaijEbMb4b4VP+9tLARgXVWdB","dein":"Quae fuerit causa, mox videro; interea hoc tenebo, si ob aliquam causam ista, quae sine dubio praeclara sunt, fecerint, virtutem iis per se ipsam causam non fuisse. -- Torquem detraxit hosti. -- Et quidem se texit, ne interiret. -- At magnum periculum adiit. -- In oculis quidem exercitus. -- Quid ex eo est consecutus? -- Laudem et caritatem, quae sunt vitae sine metu degendae praesidia firmissima. -- Filium morte multavit. -- Si sine causa, nollem me ab eo ortum, tam inportuno tamque crudeli; sin, ut dolore suo sanciret militaris imperii disciplinam exercitumque in gravissimo bello animadversionis metu contineret, saluti prospexit civium, qua intellegebat contineri suam. atque haec ratio late patet.","dolor":"0915-02-28T18:35:08.000000194Z","enim":true,"esse":"Nisi mihi Phaedrum, inquam, tu mentitum aut Zenonem putas, quorum utrumque audivi, cum mihi nihil sane praeter sedulitatem probarent, omnes mihi Epicuri sententiae satis notae sunt. atque eos, quos nominavi, cum Attico nostro frequenter audivi, cum miraretur ille quidem utrumque, Phaedrum autem etiam amaret, cotidieque inter nos ea, quae audiebamus, conferebamus, neque erat umquam controversia, quid ego intellegerem, sed quid probarem.","et":"TEIgdBLCXQk1ysa6Y11aA2ZohpEVG7wGBMFM3VLt7Bt6ausXv1tdCg==","exultat":null,"hae":0,"in":0.5627995383118168,"level":"error","me":null,"msg":"nos omnium postulet doctrinis amicitia iis nullas","nobis":null,"plerisque":false,"quem":0,"sola":"y5xCSn+1yTfcqouq4FVQeC/pKArTKHRlOLqRwC/TeXcs9LY=","time":"2014-10-27T18:38:45-04:00","volunt":null,"voluptates":"caEehN56R1l11f5Hq8GgVHcOVq5rwWN03DUay+zBBxGw7xsdviPabU+g/bWBAyYHZCqCaKt/GlAH7oj+D28sF/FhWYDORVllIs2uYJbfB0wYViqVH1ZbKwl2al7VEtwi+tI0Fo9X9iP2xSgopBL4Zr4T/V7mxodtroAeFUz3fV0z1gS3rJvWONfBC2Gap3QG16l66rIdTgbWMBzoKdXKW04921tHSJtn/CMJcKXvRD0A02TbReON"}
{"alias":"0673-09-08T23:45:16.000000224Z","autem":"Quid? si nos non interpretum fungimur munere, sed tuemur ea, quae dicta sunt ab iis quos probamus, eisque nostrum iudicium et nostrum scribendi ordinem adiungimus, quid habent, cur Graeca anteponant iis, quae et splendide dicta sint neque sint conversa de Graecis? nam si dicent ab illis has res esse tractatas, ne ipsos quidem Graecos est cur tam multos legant, quam legendi sunt. quid enim est a Chrysippo praetermissum in Stoicis? legimus tamen Diogenem, Antipatrum, Mnesarchum, Panaetium, multos alios in primisque familiarem nostrum Posidonium. quid? Theophrastus mediocriterne delectat, cum tractat locos ab Aristotele ante tractatos? quid? Epicurei num desistunt de isdem, de quibus et ab Epicuro scriptum est et ab antiquis, ad arbitrium suum scribere? quodsi Graeci leguntur a Graecis isdem de rebus alia ratione compositis, quid est, cur nostri a nostris non legantur?","causa":null,"doloribus":"0229-08-27T08:05:09.000000348Z","fabulis":0,"level":"info","msg":"appareat opes doloris eius conficiuntque intuemur posse","rebus":0,"time":"2014-10-27T18:38:45-04:00"}
time="2014-10-27T18:31:40-04:00" level="warning" msg="naturam reprehenderit dicturam victi a autem et omittam bene" In=true opinemur="At etiam Athenis, ut e patre audiebam facete et urbane Stoicos irridente, statua est in Ceramico Chrysippi sedentis porrecta manu, quae manus significet illum in hae esse rogatiuncula delectatum: 'Numquidnam manus tua sic affecta, quem ad modum affecta nunc est, desiderat?' -- Nihil sane. -- 'At, si voluptas esset bonum, desideraret.' -- Ita credo. -- 'Non est igitur voluptas bonum.' Hoc ne statuam quidem dicturam pater aiebat, si loqui posset. conclusum est enim contra Cyrenaicos satis acute, nihil ad Epicurum. nam si ea sola voluptas esset, quae quasi titillaret sensus, ut ita dicam, et ad eos cum suavitate afflueret et illaberetur, nec manus esse contenta posset nec ulla pars vacuitate doloris sine iucundo motu voluptatis. sin autem summa voluptas est, ut Epicuro placet, nihil dolere, primum tibi recte, Chrysippe, concessum est nihil desiderare manum, cum ita esset affecta, secundum non recte, si voluptas esset bonum, fuisse desideraturam. idcirco enim non desideraret, quia, quod dolore caret, id in voluptate est."
time="2014-10-27T18:38:45-04:00" level="warning" msg="faciunt aequum bene quidem et esse parabilis" atque=false et="Quid tibi, Torquate, quid huic Triario litterae, quid historiae cognitioque rerum, quid poetarum evolutio, quid tanta tot versuum memoria voluptatis affert? nec mihi illud dixeris: 'Haec enim ipsa mihi sunt voluptati, et erant illa Torquatis.' Numquam hoc ita defendit Epicurus neque Metrodorus aut quisquam eorum, qui aut saperet aliquid aut ista didicisset. et quod quaeritur saepe, cur tam multi sint Epicurei, sunt aliae quoque causae, sed multitudinem haec maxime allicit, quod ita putant dici ab illo, recta et honesta quae sint, ea facere ipsa per se laetitiam, id est voluptatem. homines optimi non intellegunt totam rationem everti, si ita res se habeat. nam si concederetur, etiamsi ad corpus nihil referatur, ista sua sponte et per se esse iucunda, per se esset et virtus et cognitio rerum, quod minime ille vult expetenda." factis=0 legimus=0.6380484281330108 morati=true tribuat=1474-12-23 03:10:55.000000246 +0000 UTC tum=0.36542830164796136
time="2014-10-27T18:38:45-04:00" level="error" msg="qua" dixi="Ego autem quem timeam lectorem, cum ad te ne Graecis quidem cedentem in philosophia audeam scribere? quamquam a te ipso id quidem facio provocatus gratissimo mihi libro, quem ad me de virtute misisti. Sed ex eo credo quibusdam usu venire; ut abhorreant a Latinis, quod inciderint in inculta quaedam et horrida, de malis Graecis Latine scripta deterius. quibus ego assentior, dum modo de isdem rebus ne Graecos quidem legendos putent. res vero bonas verbis electis graviter ornateque dictas quis non legat? nisi qui se plane Graecum dici velit, ut a Scaevola est praetore salutatus Athenis Albucius." effecerit="Sed ut iis bonis erigimur, quae expectamus, sic laetamur iis, quae recordamur. stulti autem malorum memoria torquentur, sapientes bona praeterita grata recordatione renovata delectant. est autem situm in nobis ut et adversa quasi perpetua oblivione obruamus et secunda iucunde ac suaviter meminerimus. sed cum ea, quae praeterierunt, acri animo et attento intuemur, tum fit ut aegritudo sequatur, si illa mala sint, laetitia, si bona." nec=<nil> ob="Ut autem a facillimis ordiamur, prima veniat in medium Epicuri ratio, quae plerisque notissima est. quam a nobis sic intelleges eitam, ut ab ipsis, qui eam disciplinam probant, non soleat accuratius explicari; verum enim invenire volumus, non tamquam adversarium aliquem convincere. accurate autem quondam a L. Torquato, homine omni doctrina erudito, defensa est Epicuri sententia de voluptate, a meque ei responsum, cum C. Triarius, in primis gravis et doctus adolescens, ei disputationi interesset." quidem=false quod=0 se=true sequitur=0.13815772053664252 summum=0.9507713169389637 tamen="Hanc ego cum teneam sententiam, quid est cur verear, ne ad eam non possim accommodare Torquatos nostros? quos tu paulo ante cum memoriter, tum etiam erga nos amice et benivole collegisti, nec me tamen laudandis maioribus meis corrupisti nec segniorem ad respondendum reddidisti. quorum facta quem ad modum, quaeso, interpretaris? sicine eos censes aut in armatum hostem impetum fecisse aut in liberos atque in sanguinem suum tam crudelis fuisse, nihil ut de utilitatibus, nihil ut de commodis suis cogitarent? at id ne ferae quidem faciunt, ut ita ruant itaque turbent, ut earum motus et impetus quo pertineant non intellegamus, tu tam egregios viros censes tantas res gessisse sine causa?" ut=0.9220774561265302
time="2014-10-27T18:38:45-04:00" level="error" msg="paulo esse sibi ab ratio saepe" =0 Graeci=0 arare="Quam ob rem tandem, inquit, non satisfacit? te enim iudicem aequum puto, modo quae dicat ille bene noris." enim=<nil> facere=false in="Id qui in una virtute ponunt et splendore nominis capti quid natura postulet non intellegunt, errore maximo, si Epicurum audire voluerint, liberabuntur: istae enim vestrae eximiae pulchraeque virtutes nisi voluptatem efficerent, quis eas aut laudabilis aut expetendas arbitraretur? ut enim medicorum scientiam non ipsius artis, sed bonae valetudinis causa probamus, et gubernatoris ars, quia bene navigandi rationem habet, utilitate, non arte laudatur, sic sapientia, quae ars vivendi putanda est, non expeteretur, si nihil efficeret; nunc expetitur, quod est tamquam artifex conquirendae et comparandae voluptatis --" inpendente="Quid tibi, Torquate, quid huic Triario litterae, quid historiae cognitioque rerum, quid poetarum evolutio, quid tanta tot versuum memoria voluptatis affert? nec mihi illud dixeris: 'Haec enim ipsa mihi sunt voluptati, et erant illa Torquatis.' Numquam hoc ita defendit Epicurus neque Metrodorus aut quisquam eorum, qui aut saperet aliquid aut ista didicisset. et quod quaeritur saepe, cur tam multi sint Epicurei, sunt aliae quoque causae, sed multitudinem haec maxime allicit, quod ita putant dici ab illo, recta et honesta quae sint, ea facere ipsa per se laetitiam, id est voluptatem. homines optimi non intellegunt totam rationem everti, si ita res se habeat. nam si concederetur, etiamsi ad corpus nihil referatur, ista sua sponte et per se esse iucunda, per se esset et virtus et cognitio rerum, quod minime ille vult expetenda." inviti=0.5559171606675798 ipsam=0.20452342141212618 iudicia=0.3127147602988959 non="Quam ob rem tandem, inquit, non satisfacit? te enim iudicem aequum puto, modo quae dicat ille bene noris." quidem=true se=<nil> sed="Iam in altera philosophiae parte. quae est quaerendi ac disserendi, quae logikh dicitur, iste vester plane, ut mihi quidem videtur, inermis ac nudus est. tollit definitiones, nihil de dividendo ac partiendo docet, non quo modo efficiatur concludaturque ratio tradit, non qua via captiosa solvantur ambigua distinguantur ostendit; iudicia rerum in sensibus ponit, quibus si semel aliquid falsi pro vero probatum sit, sublatum esse omne iudicium veri et falsi putat." sustulisti=1127-05-14 01:23:02.00000056 +0000 UTC tenere=0 verbis=<nil> voluptatis="Quamquam, si plane sic verterem Platonem aut Aristotelem, ut verterunt nostri poetae fabulas, male, credo, mererer de meis civibus, si ad eorum cognitionem divina illa ingenia transferrem. sed id neque feci adhuc nec mihi tamen, ne faciam, interdictum puto. locos quidem quosdam, si videbitur, transferam, et maxime ab iis, quos modo nominavi, cum inciderit, ut id apte fieri possit, ut ab Homero Ennius, Afranius a Menandro solet. Nec vero, ut noster Lucilius, recusabo, quo minus omnes mea legant. utinam esset ille Persius, Scipio vero et Rutilius multo etiam magis, quorum ille iudicium reformidans Tarentinis ait se et Consentinis et Siculis scribere. facete is quidem, sicut alia; sed neque tam docti tum erant, ad quorum iudicium elaboraret, et sunt illius scripta leviora, ut urbanitas summa appareat, doctrina mediocris."
time="2014-10-27T18:38:45-04:00" level="warning" msg="est innumerabiles ab nobis" inesse="Quae enim cupiditates a natura proficiscuntur, facile explentur sine ulla iniuria, quae autem inanes sunt, iis parendum non est. nihil enim desiderabile concupiscunt, plusque in ipsa iniuria detrimenti est quam in iis rebus emolumenti, quae pariuntur iniuria. Itaque ne iustitiam quidem recte quis dixerit per se ipsam optabilem, sed quia iucunditatis vel plurimum afferat. nam diligi et carum esse iucundum est propterea, quia tutiorem vitam et voluptatem pleniorem efficit. itaque non ob ea solum incommoda, quae eveniunt inprobis, fugiendam inprobitatem putamus, sed multo etiam magis, quod, cuius in animo versatur, numquam sinit eum respirare, numquam adquiescere." non=true praeclare=0.5869776997848529 `

	want := []map[string]Field{
		{DefaultRaw: RawField(`[raw log]2014/10/27 18:38:45 warn: map[string]interface {}{"Romano":[]uint8{0x24, 0x5a, 0x9b, 0xd8, 0x97, 0xc1, 0x66, 0x13, 0xec, 0x8f, 0x24, 0x8d, 0x87, 0x53, 0x22, 0x9, 0x3f, 0x5d, 0x42, 0x66, 0xad, 0x2c, 0xd2, 0x56, 0x27, 0x17, 0x2b, 0xa7, 0x9b, 0x6c, 0x96, 0x34, 0xc0, 0x85, 0x2a, 0xb2, 0x80, 0x81, 0xae, 0x17, 0x9e, 0x2e, 0xa, 0x98, 0x39, 0x95, 0x47, 0x5a, 0xc5, 0xeb, 0xf8, 0xf5, 0xa1, 0x67, 0x27, 0x2f, 0xa1, 0x87, 0xa5, 0x42, 0x67, 0x53, 0xe, 0xf, 0x5e, 0xaa, 0x10, 0x8d, 0x91, 0xbb, 0x7b, 0x62, 0x88, 0x7c, 0x2f, 0xf5, 0x2b, 0x2e, 0x10, 0x23, 0x13, 0xb3, 0xbb, 0x6e, 0xaa, 0xf6, 0xa3, 0xc, 0x9c, 0x4d, 0x66, 0x1a, 0x1b, 0x13, 0x90, 0x5f, 0x37, 0x52, 0x59, 0xbc, 0x94, 0xea, 0xbd, 0x33, 0x5a, 0x6a, 0xb8, 0xe7, 0xd3, 0xa5, 0xf7, 0x4e, 0x26, 0x45, 0xfd, 0x1e, 0xe6, 0x81, 0x4b, 0x76, 0x22, 0x28, 0x2c, 0x78, 0xdd, 0xce, 0xfe, 0x2c, 0xc7, 0xb4}, "rem":[]uint8{0x1e, 0x73, 0x5d, 0x93, 0x9f, 0x3d, 0xd2, 0x1a, 0x9a, 0xe8, 0xac, 0x2c, 0xc7, 0x1e, 0xbc, 0x4a, 0xa, 0xe9, 0xb4, 0xea, 0xe2, 0xca, 0xd7, 0x14, 0x3d, 0xe0, 0x10, 0x7f, 0x98, 0x1b, 0x5c, 0x7, 0x27, 0xac, 0xcf, 0x7e, 0x74, 0x45, 0x60, 0x52, 0xfa, 0x5a, 0xad, 0x84, 0xde, 0xee, 0xab, 0x68, 0xe, 0x71, 0x7e, 0x4c, 0xa9, 0x66, 0x23, 0x20, 0xf0, 0x41, 0x8e, 0x35, 0xd2, 0x6e, 0x63, 0x3e, 0xda, 0x4b, 0x57, 0xf, 0x8c, 0x4, 0xa2, 0x70, 0x1a, 0xd9, 0x6d, 0x73, 0x30, 0x35, 0xad, 0x55, 0x7d, 0x41, 0x9f, 0x54, 0x8e, 0x44, 0x7d, 0x7d, 0xab, 0x28, 0x47, 0xd2, 0x4b, 0xb9, 0x99, 0x45, 0x3, 0xbe, 0xaf, 0x26, 0xdd, 0x9a, 0xc4, 0x4a, 0x49, 0xb4, 0x32, 0xdb, 0xd5, 0x0, 0x9e, 0x6c, 0x56, 0x80, 0x88, 0x21, 0x8e, 0x76, 0x13, 0x4a, 0xf8, 0xdb, 0x52, 0x2a, 0x7f, 0x83, 0x12, 0x7b, 0x8f, 0x1b, 0xec, 0xb4, 0x79, 0x8a, 0x18, 0x1a, 0x5c, 0xcf, 0x3f, 0xe3, 0xf0, 0xf9, 0x4, 0x42, 0x96, 0x69, 0x90, 0xf9, 0x95, 0xa4, 0x94, 0x42, 0xac, 0x48, 0xc6, 0x33, 0xbc, 0x71, 0x73, 0x71, 0x52, 0xcf, 0x50, 0x41, 0x1, 0x32, 0xaa, 0x14, 0xe0, 0xf3, 0x5a, 0x1f, 0xf5, 0xab, 0x30, 0x12, 0x19, 0x48, 0x11, 0x34, 0x99, 0x48, 0xd1, 0x1b, 0xb3, 0x9a, 0x26, 0x4c, 0x13, 0xcf, 0x75, 0x1a, 0xd9, 0xbf, 0x29, 0x86, 0xf6, 0xa4, 0xef, 0x24, 0x3c, 0x64, 0x5a, 0x39, 0xcf, 0xb0, 0x7e, 0x47, 0x11, 0xb5, 0xc7, 0xec, 0x2a, 0x24, 0xd3, 0x71, 0xd1, 0x98, 0xaf, 0xa4, 0xd9, 0x6d, 0x40, 0x9, 0x76, 0x4c, 0xb6, 0x82, 0xb}, "si":true, "malivoli":[]uint8{0x24, 0x99, 0xae, 0x4f, 0x37, 0xda, 0x57, 0x74, 0xfc, 0xb1, 0x6d, 0x89, 0x7, 0x6e, 0x46, 0x5c, 0xbe, 0xd3, 0x39, 0xc6, 0xcc, 0x7c, 0xba, 0x44, 0x7e, 0x9f, 0x70, 0x2c, 0x55, 0x67, 0x27, 0x1e, 0x25, 0x43, 0xd9, 0xfb, 0x2c, 0x24, 0xfe, 0x3c, 0x5e, 0xdf, 0x5f, 0x44, 0x58, 0x26, 0xdb, 0x3, 0x1c, 0xa0, 0xfa, 0x8a, 0x7f, 0x86, 0xd2, 0x51, 0x4a, 0xd2, 0xd, 0xfc, 0x1a, 0xbe, 0xad, 0x40, 0xa, 0xe1, 0x49, 0x93, 0x53, 0x24, 0xa, 0x54, 0x3e, 0xc2, 0x55, 0xf9, 0x4d, 0x8b}, "videtur":0, "etiam":[]uint8{0xae, 0xbe, 0xd5, 0x45, 0x1f, 0x49, 0xcf, 0x3c, 0xd, 0x3b, 0x4b, 0x5f, 0xcf, 0xbe, 0xc3, 0x1b, 0xe5, 0x28, 0xb5, 0x8c, 0xe8, 0x2a, 0x5d, 0x9, 0x13, 0x53, 0xfb, 0x48, 0xc, 0x45, 0x93, 0x60, 0x5c, 0x8, 0xec, 0xe5, 0x48, 0xbb, 0xcb, 0x5c, 0xee, 0x89, 0x36, 0xb, 0x42, 0xa8, 0x9f, 0x10, 0x90, 0x5b, 0x92, 0x36, 0xd1, 0x4c, 0xc2, 0x1b, 0xf0, 0x55, 0x23, 0x7, 0xec, 0xda, 0x8a, 0x36, 0x31, 0xe1, 0x91, 0x44, 0xf, 0x79, 0x4e, 0x56, 0xf5, 0x1, 0xc6, 0x28, 0xfa, 0x35, 0x37, 0x2c, 0xad, 0xac, 0x3d, 0x76, 0x78, 0x16, 0xb3, 0x28, 0x3b, 0x48, 0x67, 0x91, 0xf9, 0x31, 0x72, 0x47, 0xca, 0x71, 0x22, 0xd0, 0xa9, 0x3c, 0x5, 0x57, 0xe9, 0x62, 0x2a, 0x73, 0xd5, 0x98, 0x48, 0x6b, 0x90, 0x4d, 0xd8, 0xa6, 0x7a, 0xe0, 0x3c, 0x5d, 0xd3, 0x70, 0x86, 0xd2, 0x48, 0x2a, 0x3d, 0x2, 0x0, 0x3f, 0x7d, 0x63, 0x82, 0x99, 0xda, 0x29, 0xd1, 0x68, 0xa8, 0x8, 0x36, 0x47, 0x7a, 0x1b, 0xc3, 0x61, 0xff, 0x56, 0x1c, 0xa, 0xb4, 0x3a, 0x8a, 0x77, 0xb0, 0x94, 0xd7, 0xa3, 0x8f, 0x7, 0xf2, 0x58, 0x34, 0xdc, 0x4a, 0xf9, 0xb, 0x0, 0xc3, 0x39}}`)},
		{DefaultRaw: RawField(`[raw log]2014/10/27 18:38:45 warn: map[string]interface {}{"nemo":false, "alteram":0.8879334884220434, "omnium":time.Time{sec:8724124955, nsec:413, loc:(*time.Location)(0x25acc0)}, "vero":time.Time{sec:39763127421, nsec:927, loc:(*time.Location)(0x25acc0)}, "quidem":0.04749476992250518, "autem":"Quodsi vitam omnem perturbari videmus errore et inscientia, sapientiamque esse solam, quae nos a libidinum impetu et a formidinum terrore vindicet et ipsius fortunae modice ferre doceat iniurias et omnis monstret vias, quae ad quietem et ad tranquillitatem ferant, quid est cur dubitemus dicere et sapientiam propter voluptates expetendam et insipientiam propter molestias esse fugiendam?", "erat":0.3939079663168565, "semper":time.Time{sec:52034192405, nsec:552, loc:(*time.Location)(0x25acc0)}, "abhorreant":true}`)},
		{DefaultRaw: RawField(`[raw log]2014/10/27 18:38:45 warn: map[string]interface {}{"praesertim":interface {}(nil)}`)},
		{
			"level": StringField("warning"),
			"mea":   NumberField(0.6645600532184904),
			"msg":   StringField("esse quidem veri viros diligant summum esse ut"),
			"time":  TimeField{time.Date(2014, 10, 27, 18, 31, 40, 0, time.FixedZone("EDT", -4*60*60))},
		}, {
			"diu":        BooleanField(true),
			"enim":       BooleanField(false),
			"est":        NilField{},
			"et":         NumberField(0.6408415373414884),
			"expetendam": BooleanField(false),
			"huius":      BooleanField(false),
			"legimus":    NilField{},
			"level":      StringField("warning"),
			"msg":        StringField("sic stare"),
			"optinere":   NumberField(0.3473303755727788),
			"quo":        StringField("Ita, quae mutat, ea corrumpit, quae sequitur sunt tota Democriti, atomi, inane, imagines, quae eidola nominant, quorum incursione non solum videamus, sed etiam cogitemus; infinitio ipsa, quam apeirian vocant, tota ab illo est, tum innumerabiles munuperatum."),
			"quodsi":     BooleanField(false),
			"rerum":      TimeField{time.Date(1799, 11, 03, 12, 21, 50, 422, time.UTC)},
			"time":       TimeField{time.Date(2014, 10, 27, 18, 38, 45, 0, time.FixedZone("EDT", -4*60*60))},
			"ut":         StringField("At etiam Atheatum: 'Numquidnam manus tua sic affecta, quem ad modum affecta nunc est, desiderat?' -- Nihil sane. -- 'At, si voluptas esset bonum, desideraret.' -- Ita credo. -- 'Non est igitur voluptas bonum.' Hoc ne statuaml dolere, primum tibi recte, Chrysippe, concessum est nihil desiderare manum, cum ita esset affecta, secundum non recte, si voluptas esset bonum, fuisse desideraturam. idcirco enim non desideraret, quia, quod dolore caret, id in voluptate est."),
		}, {
			"apud":       NumberField(0.04118484163828559),
			"ardore":     TimeField{time.Date(1486, 3, 5, 8, 23, 42, 525, time.UTC)},
			"conducunt":  StringField("OI18KPtbDwYlJI95ycLGZuJegXmcW84BSfhE85XDS3gGLrDgyxxBaijEbMb4b4VP+9tLARgXVWdB"),
			"dein":       StringField("Quae fuerit causa, mox videro; interea hoc tenebo, si ob aliquam causam ista, quae sine dubio praeclara sunt, fecerint, virtutem iis per se ipsam causam non fuisse. -- Torquem detraxit hosti. -- Et quidem se texit, ne interiret. -- At magnum periculum adiit. -- In oculis quidem exercitus. -- Quid ex eo est consecutus? -- Laudem et caritatem, quae sunt vitae sine metu degendae praesidia firmissima. -- Filium morte multavit. -- Si sine causa, nollem me ab eo ortum, tam inportuno tamque crudeli; sin, ut dolore suo sanciret militaris imperii disciplinam exercitumque in gravissimo bello animadversionis metu contineret, saluti prospexit civium, qua intellegebat contineri suam. atque haec ratio late patet."),
			"dolor":      TimeField{time.Date(915, 02, 28, 18, 35, 8, 194, time.UTC)},
			"enim":       BooleanField(true),
			"esse":       StringField("Nisi mihi Phaedrum, inquam, tu mentitum aut Zenonem putas, quorum utrumque audivi, cum mihi nihil sane praeter sedulitatem probarent, omnes mihi Epicuri sententiae satis notae sunt. atque eos, quos nominavi, cum Attico nostro frequenter audivi, cum miraretur ille quidem utrumque, Phaedrum autem etiam amaret, cotidieque inter nos ea, quae audiebamus, conferebamus, neque erat umquam controversia, quid ego intellegerem, sed quid probarem."),
			"et":         StringField("TEIgdBLCXQk1ysa6Y11aA2ZohpEVG7wGBMFM3VLt7Bt6ausXv1tdCg=="),
			"exultat":    NilField{},
			"hae":        NumberField(0),
			"in":         NumberField(0.5627995383118168),
			"level":      StringField("error"),
			"me":         NilField{},
			"msg":        StringField("nos omnium postulet doctrinis amicitia iis nullas"),
			"nobis":      NilField{},
			"plerisque":  BooleanField(false),
			"quem":       NumberField(0),
			"sola":       StringField("y5xCSn+1yTfcqouq4FVQeC/pKArTKHRlOLqRwC/TeXcs9LY="),
			"time":       TimeField{time.Date(2014, 10, 27, 18, 38, 45, 0, time.FixedZone("EDT", -4*60*60))},
			"volunt":     NilField{},
			"voluptates": StringField("caEehN56R1l11f5Hq8GgVHcOVq5rwWN03DUay+zBBxGw7xsdviPabU+g/bWBAyYHZCqCaKt/GlAH7oj+D28sF/FhWYDORVllIs2uYJbfB0wYViqVH1ZbKwl2al7VEtwi+tI0Fo9X9iP2xSgopBL4Zr4T/V7mxodtroAeFUz3fV0z1gS3rJvWONfBC2Gap3QG16l66rIdTgbWMBzoKdXKW04921tHSJtn/CMJcKXvRD0A02TbReON"),
		}, {
			"alias":     TimeField{time.Date(673, 9, 8, 23, 45, 16, 224, time.UTC)},
			"autem":     StringField("Quid? si nos non interpretum fungimur munere, sed tuemur ea, quae dicta sunt ab iis quos probamus, eisque nostrum iudicium et nostrum scribendi ordinem adiungimus, quid habent, cur Graeca anteponant iis, quae et splendide dicta sint neque sint conversa de Graecis? nam si dicent ab illis has res esse tractatas, ne ipsos quidem Graecos est cur tam multos legant, quam legendi sunt. quid enim est a Chrysippo praetermissum in Stoicis? legimus tamen Diogenem, Antipatrum, Mnesarchum, Panaetium, multos alios in primisque familiarem nostrum Posidonium. quid? Theophrastus mediocriterne delectat, cum tractat locos ab Aristotele ante tractatos? quid? Epicurei num desistunt de isdem, de quibus et ab Epicuro scriptum est et ab antiquis, ad arbitrium suum scribere? quodsi Graeci leguntur a Graecis isdem de rebus alia ratione compositis, quid est, cur nostri a nostris non legantur?"),
			"causa":     NilField{},
			"doloribus": TimeField{time.Date(229, 8, 27, 8, 5, 9, 348, time.UTC)},
			"fabulis":   NumberField(0),
			"level":     StringField("info"),
			"msg":       StringField("appareat opes doloris eius conficiuntque intuemur posse"),
			"rebus":     NumberField(0),
			"time":      TimeField{time.Date(2014, 10, 27, 18, 38, 45, 0, time.FixedZone("EDT", -4*60*60))},
		}, {
			"level":    StringField("warning"),
			"msg":      StringField("naturam reprehenderit dicturam victi a autem et omittam bene"),
			"time":     TimeField{time.Date(2014, 10, 27, 18, 31, 40, 0, time.FixedZone("EDT", -4*60*60))},
			"In":       BooleanField(true),
			"opinemur": StringField("At etiam Athenis, ut e patre audiebam facete et urbane Stoicos irridente, statua est in Ceramico Chrysippi sedentis porrecta manu, quae manus significet illum in hae esse rogatiuncula delectatum: 'Numquidnam manus tua sic affecta, quem ad modum affecta nunc est, desiderat?' -- Nihil sane. -- 'At, si voluptas esset bonum, desideraret.' -- Ita credo. -- 'Non est igitur voluptas bonum.' Hoc ne statuam quidem dicturam pater aiebat, si loqui posset. conclusum est enim contra Cyrenaicos satis acute, nihil ad Epicurum. nam si ea sola voluptas esset, quae quasi titillaret sensus, ut ita dicam, et ad eos cum suavitate afflueret et illaberetur, nec manus esse contenta posset nec ulla pars vacuitate doloris sine iucundo motu voluptatis. sin autem summa voluptas est, ut Epicuro placet, nihil dolere, primum tibi recte, Chrysippe, concessum est nihil desiderare manum, cum ita esset affecta, secundum non recte, si voluptas esset bonum, fuisse desideraturam. idcirco enim non desideraret, quia, quod dolore caret, id in voluptate est."),
		}, {
			"time":    TimeField{time.Date(2014, 10, 27, 18, 38, 45, 0, time.FixedZone("EDT", -4*60*60))},
			"level":   StringField("warning"),
			"msg":     StringField("faciunt aequum bene quidem et esse parabilis"),
			"atque":   BooleanField(false),
			"et":      StringField("Quid tibi, Torquate, quid huic Triario litterae, quid historiae cognitioque rerum, quid poetarum evolutio, quid tanta tot versuum memoria voluptatis affert? nec mihi illud dixeris: 'Haec enim ipsa mihi sunt voluptati, et erant illa Torquatis.' Numquam hoc ita defendit Epicurus neque Metrodorus aut quisquam eorum, qui aut saperet aliquid aut ista didicisset. et quod quaeritur saepe, cur tam multi sint Epicurei, sunt aliae quoque causae, sed multitudinem haec maxime allicit, quod ita putant dici ab illo, recta et honesta quae sint, ea facere ipsa per se laetitiam, id est voluptatem. homines optimi non intellegunt totam rationem everti, si ita res se habeat. nam si concederetur, etiamsi ad corpus nihil referatur, ista sua sponte et per se esse iucunda, per se esset et virtus et cognitio rerum, quod minime ille vult expetenda."),
			"factis":  NumberField(0),
			"legimus": NumberField(0.6380484281330108),
			"morati":  BooleanField(true),
			"tribuat": TimeField{time.Date(1474, 12, 23, 03, 10, 55, 246, time.UTC)},
			"tum":     NumberField(0.36542830164796136),
		}, {
			"time":      TimeField{time.Date(2014, 10, 27, 18, 38, 45, 0, time.FixedZone("EDT", -4*60*60))},
			"level":     StringField("error"),
			"msg":       StringField("qua"),
			"dixi":      StringField("Ego autem quem timeam lectorem, cum ad te ne Graecis quidem cedentem in philosophia audeam scribere? quamquam a te ipso id quidem facio provocatus gratissimo mihi libro, quem ad me de virtute misisti. Sed ex eo credo quibusdam usu venire; ut abhorreant a Latinis, quod inciderint in inculta quaedam et horrida, de malis Graecis Latine scripta deterius. quibus ego assentior, dum modo de isdem rebus ne Graecos quidem legendos putent. res vero bonas verbis electis graviter ornateque dictas quis non legat? nisi qui se plane Graecum dici velit, ut a Scaevola est praetore salutatus Athenis Albucius."),
			"effecerit": StringField("Sed ut iis bonis erigimur, quae expectamus, sic laetamur iis, quae recordamur. stulti autem malorum memoria torquentur, sapientes bona praeterita grata recordatione renovata delectant. est autem situm in nobis ut et adversa quasi perpetua oblivione obruamus et secunda iucunde ac suaviter meminerimus. sed cum ea, quae praeterierunt, acri animo et attento intuemur, tum fit ut aegritudo sequatur, si illa mala sint, laetitia, si bona."),
			"nec":       NilField{},
			"ob":        StringField("Ut autem a facillimis ordiamur, prima veniat in medium Epicuri ratio, quae plerisque notissima est. quam a nobis sic intelleges eitam, ut ab ipsis, qui eam disciplinam probant, non soleat accuratius explicari; verum enim invenire volumus, non tamquam adversarium aliquem convincere. accurate autem quondam a L. Torquato, homine omni doctrina erudito, defensa est Epicuri sententia de voluptate, a meque ei responsum, cum C. Triarius, in primis gravis et doctus adolescens, ei disputationi interesset."),
			"quidem":    BooleanField(false),
			"quod":      NumberField(0),
			"se":        BooleanField(true),
			"sequitur":  NumberField(0.13815772053664252),
			"summum":    NumberField(0.9507713169389637),
			"tamen":     StringField("Hanc ego cum teneam sententiam, quid est cur verear, ne ad eam non possim accommodare Torquatos nostros? quos tu paulo ante cum memoriter, tum etiam erga nos amice et benivole collegisti, nec me tamen laudandis maioribus meis corrupisti nec segniorem ad respondendum reddidisti. quorum facta quem ad modum, quaeso, interpretaris? sicine eos censes aut in armatum hostem impetum fecisse aut in liberos atque in sanguinem suum tam crudelis fuisse, nihil ut de utilitatibus, nihil ut de commodis suis cogitarent? at id ne ferae quidem faciunt, ut ita ruant itaque turbent, ut earum motus et impetus quo pertineant non intellegamus, tu tam egregios viros censes tantas res gessisse sine causa?"),
			"ut":        NumberField(0.9220774561265302),
		}, {
			"time":       TimeField{time.Date(2014, 10, 27, 18, 38, 45, 0, time.FixedZone("EDT", -4*60*60))},
			"level":      StringField("error"),
			"msg":        StringField("paulo esse sibi ab ratio saepe"),
			"":           NumberField(0),
			"Graeci":     NumberField(0),
			"arare":      StringField("Quam ob rem tandem, inquit, non satisfacit? te enim iudicem aequum puto, modo quae dicat ille bene noris."),
			"enim":       NilField{},
			"facere":     BooleanField(false),
			"in":         StringField("Id qui in una virtute ponunt et splendore nominis capti quid natura postulet non intellegunt, errore maximo, si Epicurum audire voluerint, liberabuntur: istae enim vestrae eximiae pulchraeque virtutes nisi voluptatem efficerent, quis eas aut laudabilis aut expetendas arbitraretur? ut enim medicorum scientiam non ipsius artis, sed bonae valetudinis causa probamus, et gubernatoris ars, quia bene navigandi rationem habet, utilitate, non arte laudatur, sic sapientia, quae ars vivendi putanda est, non expeteretur, si nihil efficeret; nunc expetitur, quod est tamquam artifex conquirendae et comparandae voluptatis --"),
			"inpendente": StringField("Quid tibi, Torquate, quid huic Triario litterae, quid historiae cognitioque rerum, quid poetarum evolutio, quid tanta tot versuum memoria voluptatis affert? nec mihi illud dixeris: 'Haec enim ipsa mihi sunt voluptati, et erant illa Torquatis.' Numquam hoc ita defendit Epicurus neque Metrodorus aut quisquam eorum, qui aut saperet aliquid aut ista didicisset. et quod quaeritur saepe, cur tam multi sint Epicurei, sunt aliae quoque causae, sed multitudinem haec maxime allicit, quod ita putant dici ab illo, recta et honesta quae sint, ea facere ipsa per se laetitiam, id est voluptatem. homines optimi non intellegunt totam rationem everti, si ita res se habeat. nam si concederetur, etiamsi ad corpus nihil referatur, ista sua sponte et per se esse iucunda, per se esset et virtus et cognitio rerum, quod minime ille vult expetenda."),
			"inviti":     NumberField(0.5559171606675798),
			"ipsam":      NumberField(0.20452342141212618),
			"iudicia":    NumberField(0.3127147602988959),
			"non":        StringField("Quam ob rem tandem, inquit, non satisfacit? te enim iudicem aequum puto, modo quae dicat ille bene noris."),
			"quidem":     BooleanField(true),
			"se":         NilField{},
			"sed":        StringField("Iam in altera philosophiae parte. quae est quaerendi ac disserendi, quae logikh dicitur, iste vester plane, ut mihi quidem videtur, inermis ac nudus est. tollit definitiones, nihil de dividendo ac partiendo docet, non quo modo efficiatur concludaturque ratio tradit, non qua via captiosa solvantur ambigua distinguantur ostendit; iudicia rerum in sensibus ponit, quibus si semel aliquid falsi pro vero probatum sit, sublatum esse omne iudicium veri et falsi putat."),
			"sustulisti": TimeField{time.Date(1127, 5, 14, 1, 23, 2, 560, time.UTC)},
			"tenere":     NumberField(0),
			"verbis":     NilField{},
			"voluptatis": StringField("Quamquam, si plane sic verterem Platonem aut Aristotelem, ut verterunt nostri poetae fabulas, male, credo, mererer de meis civibus, si ad eorum cognitionem divina illa ingenia transferrem. sed id neque feci adhuc nec mihi tamen, ne faciam, interdictum puto. locos quidem quosdam, si videbitur, transferam, et maxime ab iis, quos modo nominavi, cum inciderit, ut id apte fieri possit, ut ab Homero Ennius, Afranius a Menandro solet. Nec vero, ut noster Lucilius, recusabo, quo minus omnes mea legant. utinam esset ille Persius, Scipio vero et Rutilius multo etiam magis, quorum ille iudicium reformidans Tarentinis ait se et Consentinis et Siculis scribere. facete is quidem, sicut alia; sed neque tam docti tum erant, ad quorum iudicium elaboraret, et sunt illius scripta leviora, ut urbanitas summa appareat, doctrina mediocris."),
		}, {
			"time":      TimeField{time.Date(2014, 10, 27, 18, 38, 45, 0, time.FixedZone("EDT", -4*60*60))},
			"level":     StringField("warning"),
			"msg":       StringField("est innumerabiles ab nobis"),
			"inesse":    StringField("Quae enim cupiditates a natura proficiscuntur, facile explentur sine ulla iniuria, quae autem inanes sunt, iis parendum non est. nihil enim desiderabile concupiscunt, plusque in ipsa iniuria detrimenti est quam in iis rebus emolumenti, quae pariuntur iniuria. Itaque ne iustitiam quidem recte quis dixerit per se ipsam optabilem, sed quia iucunditatis vel plurimum afferat. nam diligi et carum esse iucundum est propterea, quia tutiorem vitam et voluptatem pleniorem efficit. itaque non ob ea solum incommoda, quae eveniunt inprobis, fugiendam inprobitatem putamus, sed multo etiam magis, quod, cuius in animo versatur, numquam sinit eum respirare, numquam adquiescere."),
			"non":       BooleanField(true),
			"praeclare": NumberField(0.5869776997848529),
		},
	}

	parser := NewParser(strings.NewReader(input))
	i := 0
	for ; parser.Next(); i++ {
		gotE := parser.LogEntry()
		if i >= len(want) {
			t.Errorf("parsed too many lines: want %d, got %d so far", len(want), i)
			continue
		}
		wantE := want[i]
		checkEntryMatch(t, wantE, gotE)
	}
	if i != len(want) {
		t.Fatalf("parsed too wrong number of lines: want %d, got %d", len(want), i)
	}

	err := parser.Err()
	if err != nil {
		t.Fatalf("got parsing error: %v", err)
	}
}
