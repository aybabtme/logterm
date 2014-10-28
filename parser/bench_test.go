package parser

import (
	"bytes"
	"github.com/aybabtme/logterm/testutil/genlogs"
	"io"
	"testing"
)

func BenchmarkParseRandomLogs(b *testing.B) {
	buf := prepareBenchData(10000, genlogs.NewRandLogger)
	benchmarkLogParse(b, buf)
}

func BenchmarkParseJSONLogs(b *testing.B) {
	buf := prepareBenchData(10000, genlogs.NewJSONLogger)
	benchmarkLogParse(b, buf)
}

func BenchmarkParseFmtLogs(b *testing.B) {
	buf := prepareBenchData(10000, genlogs.NewFmtLogger)
	benchmarkLogParse(b, buf)
}

func BenchmarkParseRawLogs(b *testing.B) {
	buf := prepareBenchData(10000, genlogs.NewRawLogger)
	benchmarkLogParse(b, buf)
}

func prepareBenchData(lines int, logFunc func(io.Writer) *genlogs.RandLogger) *bytes.Reader {
	buf := bytes.NewBuffer(nil)
	rlog := logFunc(buf)
	for i := 0; i < lines; i++ {
		rlog.Log()
	}
	return bytes.NewReader(buf.Bytes())
}

func benchmarkLogParse(b *testing.B, buf *bytes.Reader) {
	b.SetBytes(int64(buf.Len()))
	b.ReportAllocs()
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = buf.Seek(0, 0)
		parser := NewParser(buf)
		b.StartTimer()
		for parser.Next() {
			e := parser.LogEntry()
			_ = len(e.FieldNames())
		}
		b.StopTimer()
	}
}
