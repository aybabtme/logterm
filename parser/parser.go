package parser

import (
	"bufio"
	"io"
)

const DefaultRaw = "raw"

type Parser struct {
	scan          *bufio.Scanner
	allowEmptyKey bool
}

func NewParser(r io.Reader) *Parser {
	scan := bufio.NewScanner(r)
	scan.Split(bufio.ScanLines)
	return &Parser{scan: scan, allowEmptyKey: true}
}

func (p *Parser) Next() bool { return p.scan.Scan() }

func (p *Parser) Err() error { return p.scan.Err() }

func (p *Parser) LogEntry() *Entry {
	data := p.scan.Bytes()
	switch data[0] {
	case byte('{'):
		e, ok := parseJSON(data)
		if ok {
			return e
		}
	}

	if e, ok := parseLogFmt(data, p.allowEmptyKey); ok {
		return e
	}

	e := newEntry()
	e.setField(DefaultRaw, RawField(data))

	return e
}
