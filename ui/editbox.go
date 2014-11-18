package ui

import (
	"github.com/nsf/termbox-go"
	"log"
	"unicode"
)

var (
	_ ResizeHandler = &EditBox{}
	_ InputHandler  = &EditBox{}
)

type EditBox struct {
	win    *Window
	width  int
	buffer []rune
	lines  []rune
	start  int
	cursor int
}

func NewEditBox(win *Window) *EditBox {
	lines := make([]rune, win.Width())
	for i := range lines {
		lines[i] = ' '
	}
	e := &EditBox{
		win:   win,
		width: win.Width(),
		lines: lines,
	}
	e.drawLine()
	return e
}

func (e *EditBox) Resize(x, y, width, height int) {
	if width > len(e.lines) {
		e.lines = e.lines[:width]
	}
	if width < len(e.lines) {
		e.lines = append(e.lines, make([]rune, len(e.lines)-width)...)
	}
}

func (e *EditBox) KeyPress(ch rune, key termbox.Key, mod termbox.Modifier) {
	log.Printf("ch=%#v\tkey=%#v\tmod=%#v", ch, key, mod)

	switch key {
	case 0:
		//continue
	case 0xffea: // right arrow
		if e.cursor < len(e.buffer) {
			e.cursor++
			e.drawLine()
		}
		return
	case 0xffeb: // left arrow
		if e.cursor > 0 {
			e.cursor--
			e.drawLine()
		}
		return
	case 0x7f: // backspace
		if e.cursor > 0 {
			e.cursor--
			buf := make([]rune, len(e.buffer)-1)
			copy(buf, e.buffer[:e.cursor])
			copy(buf[e.cursor:], e.buffer[e.cursor+1:])
			e.buffer = buf
			e.drawLine()
		}
		return
	case 0x1: // ctrl+a
		e.cursor = 0
		e.drawLine()
		return
	case 0x5: // ctrl+e
		e.cursor = len(e.buffer)
		e.drawLine()
		return
	case 0x20: // space
		ch = ' '
	}

	if mod != 0 {
		return
	}

	if unicode.IsPrint(ch) {
		if e.cursor == len(e.buffer) {
			e.buffer = append(e.buffer, ch)
		} else {
			pre := e.buffer[:e.cursor]
			post := e.buffer[e.cursor:]
			line := make([]rune, len(e.buffer)+1)
			copy(line, pre)
			line[e.cursor] = ch
			copy(line[e.cursor+1:], post)
			e.buffer = line
		}
		e.cursor++
	}
	e.drawLine()
}

func (e *EditBox) Mouse(termbox.Event) {}

func (e *EditBox) drawLine() {

	if len(e.buffer) > len(e.lines) {
		buf := e.buffer[len(e.buffer)-len(e.lines) : len(e.buffer)]
		copy(e.lines, buf)
	} else {
		copy(e.lines, e.buffer)
	}

	for i := len(e.buffer); i < len(e.lines); i++ {
		e.lines[i] = 0x0
	}
	if e.cursor == len(e.lines) {
		e.win.Draw(len(e.lines), 0, ' ', termbox.ColorWhite, termbox.ColorBlue)
	}
	for i, r := range e.lines {
		if i == e.cursor {
			e.win.Draw(i, 0, r, termbox.ColorWhite, termbox.ColorBlue)
		} else if r == 0x0 {
			e.win.Draw(i, 0, ' ', termbox.ColorBlack, termbox.ColorWhite)
		} else {
			e.win.Draw(i, 0, r, termbox.ColorBlack, termbox.ColorWhite)
		}
	}
}
