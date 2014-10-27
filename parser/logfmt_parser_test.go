package parser

import (
	"bytes"
	"reflect"
	"sort"
	"testing"
)

type byKeyName []kv

func (b byKeyName) Len() int           { return len(b) }
func (b byKeyName) Less(i, j int) bool { return bytes.Compare(b[i].key, b[j].key) == -1 }
func (b byKeyName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func TestScanKeyValue(t *testing.T) {
	var tests = []struct {
		input string
		want  []kv
	}{
		{
			input: "hello=bye",
			want: []kv{
				{key: []byte("hello"), val: []byte("bye")},
			},
		},
		{
			input: "hello=bye crap crap",
			want: []kv{
				{key: []byte("hello"), val: []byte("bye crap crap")},
			},
		},
		{
			input: "hello=bye crap crap    ",
			want: []kv{
				{key: []byte("hello"), val: []byte("bye crap crap    ")},
			},
		},
		{
			input: "hello=bye crap crap allo=more crap",
			want: []kv{
				{key: []byte("hello"), val: []byte("bye crap crap")},
				{key: []byte("allo"), val: []byte("more crap")},
			},
		},
		{
			input: `hello="bye crap crap" allo=more crap`,
			want: []kv{
				{key: []byte("hello"), val: []byte(`"bye crap crap"`)},
				{key: []byte("allo"), val: []byte("more crap")},
			},
		},
		{
			input: `hello="bye crap\" crap" allo=more crap`,
			want: []kv{
				{key: []byte("hello"), val: []byte(`"bye crap\" crap"`)},
				{key: []byte("allo"), val: []byte("more crap")},
			},
		},
		{
			input: `hello="bye crap\\" allo=more crap`,
			want: []kv{
				{key: []byte("hello"), val: []byte(`"bye crap\\"`)},
				{key: []byte("allo"), val: []byte("more crap")},
			},
		},
		{
			input: " hello=bye",
			want: []kv{
				{key: []byte("hello"), val: []byte("bye")},
			},
		},
		{
			input: " hello=bye crap crap",
			want: []kv{
				{key: []byte("hello"), val: []byte("bye crap crap")},
			},
		},
		{
			input: " hello=bye crap crap    ",
			want: []kv{
				{key: []byte("hello"), val: []byte("bye crap crap    ")},
			},
		},
		{
			input: " hello=bye crap crap allo=more crap",
			want: []kv{
				{key: []byte("hello"), val: []byte("bye crap crap")},
				{key: []byte("allo"), val: []byte("more crap")},
			},
		},
		{
			input: ` hello="bye crap crap" allo=more crap`,
			want: []kv{
				{key: []byte("hello"), val: []byte(`"bye crap crap"`)},
				{key: []byte("allo"), val: []byte("more crap")},
			},
		},

		{
			input: ` hello="bye crap=crap" allo=more crap`,
			want: []kv{
				{key: []byte("hello"), val: []byte(`"bye crap=crap"`)},
				{key: []byte("allo"), val: []byte("more crap")},
			},
		},
	}

	for n, tt := range tests {
		t.Logf("test #%d", n)
		got := scanAllKeyValue([]byte(tt.input))
		sort.Sort(byKeyName(tt.want))
		sort.Sort(byKeyName(got))
		if !reflect.DeepEqual(tt.want, got) {
			t.Logf("want=%v", tt.want)
			t.Logf(" got=%v", got)
			t.Fatalf("different KVs for %q", tt.input)
		}
	}
}

func TestFindWordFollowedBy(t *testing.T) {
	var tests = []struct {
		input string
		from  int
		found bool
		want  string
	}{
		{
			input: "hello=bye",
			from:  0,
			found: true,
			want:  "hello",
		},
		{
			input: "hello=bye aloa allo=bye",
			from:  6,
			found: true,
			want:  "allo",
		},
		{
			input: "hello=bye aloa allo.fr=bye",
			from:  6,
			found: true,
			want:  "allo.fr",
		},
		{
			input: " hello=bye",
			from:  0,
			found: true,
			want:  "hello",
		},
		{
			input: " hello=bye",
			from:  1,
			found: true,
			want:  "hello",
		},
		{
			input: "allo hello=bye",
			from:  0,
			found: true,
			want:  "hello",
		},
		{
			input: "allo hello=bye",
			from:  1,
			found: true,
			want:  "hello",
		},
		{
			input: " allo hello=bye",
			from:  0,
			found: true,
			want:  "hello",
		},
		{
			input: " allo hello=bye",
			from:  1,
			found: true,
			want:  "hello",
		},
		{
			input: " hello ",
			from:  0,
			found: false,
		},
		{
			input: " hello ",
			from:  1,
			found: false,
		},
		{
			input: " hello =bye",
			from:  0,
			found: true,
			want:  "",
		},
		{
			input: "hello =bye",
			from:  0,
			found: true,
			want:  "",
		},
		{
			input: " =bye",
			from:  0,
			found: true,
			want:  "",
		},
		{
			input: "=bye",
			from:  0,
			found: true,
			want:  "",
		},
		{
			input: "",
			from:  0,
			found: false,
		},
		{
			input: "=",
			from:  0,
			found: true,
			want:  "",
		},
	}

	for n, tt := range tests {
		t.Logf("test #%d", n)
		start, end, found := findWordFollowedBy('=', []byte(tt.input), tt.from)
		if found != tt.found {
			t.Errorf("want found %v, got %v", tt.found, found)
		}

		if !found {
			continue
		}

		got := string([]byte(tt.input)[start:end])

		if got != tt.want {
			t.Fatalf("want start %q, got %q", tt.want, got)
		}

	}
}

func TestFindUnescaped(t *testing.T) {
	var tests = []struct {
		input    string
		find     rune
		escape   rune
		from     int
		found    bool
		wantRest string
	}{
		{
			input:  "input",
			find:   '"',
			escape: '\\',
			from:   0,
			found:  false,
		},
		{
			input:    `inp"ut`,
			find:     '"',
			escape:   '\\',
			from:     0,
			found:    true,
			wantRest: `"ut`,
		},
		{
			input:  `inp\"ut`,
			find:   '"',
			escape: '\\',
			from:   0,
			found:  false,
		},
		{
			input:    `inp\\"ut`,
			find:     '"',
			escape:   '\\',
			from:     0,
			found:    true,
			wantRest: `"ut`,
		},
		{
			input:  `inp\\\"ut`,
			find:   '"',
			escape: '\\',
			from:   0,
			found:  false,
		},
	}

	for n, tt := range tests {
		t.Logf("test #%d", n)
		idx := findUnescaped(tt.find, tt.escape, []byte(tt.input), tt.from)
		if idx == -1 && tt.found {
			t.Fatalf("should have found %q in %q", tt.wantRest, tt.input)
		}
		if !tt.found {
			continue
		}
		gotRest := string([]byte(tt.input)[idx:])
		if tt.wantRest != gotRest {
			t.Fatalf("want %q, got %q", tt.wantRest, gotRest)
		}
	}
}
