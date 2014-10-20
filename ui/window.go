package ui

import (
	"github.com/nsf/termbox-go"
	"sync"
)

type Refresher interface {
	Refresh()
}

type Window struct {
	mu     sync.RWMutex
	canvas *Canvas
	x      int
	y      int
	width  int
	height int
}

func (w *Window) Resize(x, y, width, height int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.x = x
	w.y = y
	w.width = width
	w.height = height
}

func (w *Window) Draw(x, y int, ch rune, fg, bg termbox.Attribute) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.canvas.Set(w.x+x, w.y+y, ch, fg, bg)
}

func (w *Window) Width() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.width
}

func (w *Window) Height() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.height
}
