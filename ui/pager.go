package ui

import (
	"bytes"
	"github.com/aybabtme/linehistory"
	"io"
	"unicode/utf8"
)

const maxRuneLen = 4 // assume the max length of a utf8 rune is 4

// PagerBox shows text from an io.ReadSeeker.
type PagerBox struct {
	win    *Window
	width  int
	height int
	lines  linehistory.History
}

func NewPagerBox(win *Window) *PagerBox {
	maxSize := win.Width() * win.Height() * maxRuneLen
	return &PagerBox{
		win:   win,
		lines: linehistory.NewRing(maxSize, '\n'),
	}
}

func (p *PagerBox) canShowRunes() int {
	return p.win.Height() * p.win.Width()
}

func (p *PagerBox) Write(b []byte) (int, error) {
	p.lines.Add(b)
	p.Refresh()
	return len(b), nil
}

func (p *PagerBox) Refresh() {
	var lines []*bytes.Buffer
	var totalSpan int
	p.lines.Walk(func(line []byte) {
		span := (utf8.RuneCount(line)-1)/p.win.Width() + 1
		totalSpan += span
		// when a line is wider than the pager, break it in many lines
		for i := 0; i < span; i++ {
			from := i * p.win.Width()
			to := imin(from+p.win.Width(), len(line))
			lines = append(lines, bytes.NewBuffer(line[from:to]))
		}
	})

	// discard lines that are higher than what the pager can show
	start := totalSpan - p.win.Height()
	if start >= 0 {
		lines = lines[start:]
	} else {
		// need to pad with empty lines
		missing := -1 * start
		empties := make([]*bytes.Buffer, 0, missing)
		for ; missing > 0; missing-- {
			empties = append(empties, bytes.NewBuffer([]byte("\n")))
		}
		lines = append(empties, lines...)
	}

	p.drawLines(lines)
}

func (p *PagerBox) drawLines(lines []*bytes.Buffer) {

	var (
		x   int
		r   rune
		err error
	)
	for y, line := range lines {

		for err != io.EOF {
			x++
			r, _, err = line.ReadRune()
			p.win.Draw(x, y, r, 0, 0)
		}
		for ; x < p.win.Width(); x++ {
			p.win.Draw(x, y, ' ', 0, 0)
		}
		err = nil
		x = 0
	}
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
