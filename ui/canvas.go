package ui

import (
	"github.com/nsf/termbox-go"
	"log"
	"sync"
	"time"
)

type ResizeHandler interface {
	Resize(x, y, width, height int)
}

type InputHandler interface {
	KeyPress(termbox.Key, termbox.Modifier)
	Mouse(termbox.Event)
}

type Canvas struct {
	mu            *sync.Mutex
	done          chan struct{}
	tick          *time.Ticker
	dirty         bool
	width, height int
	cells         []termbox.Cell
}

// NewCanvas that repaints at a frequency, if the canvas has changed.
func NewCanvas(refreshFreqHz int) (*Canvas, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}
	w, h := termbox.Size()
	return &Canvas{
		mu:     &sync.Mutex{},
		done:   make(chan struct{}),
		tick:   time.NewTicker(time.Second / time.Duration(refreshFreqHz)),
		dirty:  true, // force first refresh
		width:  w,
		height: h,
		cells:  termbox.CellBuffer(),
	}, nil
}

func (c *Canvas) Run(resizer []ResizeHandler, inputer []InputHandler) error {
	go c.draw()
	return c.pollEvents(resizer, inputer)
}

// Close stops the screen refreshes and
func (c *Canvas) Close() {
	select {
	case <-c.done:
	default:
		termbox.Close()
		close(c.done)
		c.tick.Stop()
	}
}

func (c *Canvas) Fullscreen() *Window {
	return &Window{
		canvas: c,
		x:      0,
		y:      0,
		width:  c.width,
		height: c.height,
	}
}

func (c *Canvas) Size() (width, height int) {
	return termbox.Size()
}

func (c *Canvas) Set(x, y int, ch rune, fg, bg termbox.Attribute) {
	c.mu.Lock()
	defer c.mu.Unlock()
	i := c.computeIndex(x, y)
	if i >= len(c.cells) {
		return
	}
	cell := c.cells[i]
	newCell := termbox.Cell{Ch: ch, Fg: fg, Bg: bg}
	if cell != newCell {
		c.dirty = true
		c.cells[i] = newCell
	}
}

func (c *Canvas) draw() {
	defer c.Close()
	for {
		select {
		case <-c.done:
			return
		case <-c.tick.C:
			if !c.dirty {
				// don't draw if no updates
				continue
			}
			err := termbox.Flush()
			if err != nil {
				log.Printf("couldn't flush termbox: %v", err)
			}
			c.cells = termbox.CellBuffer()
		}
	}
}

func (c *Canvas) pollEvents(resizer []ResizeHandler, inputer []InputHandler) error {
	defer c.Close()
	for {
		select {
		case <-c.done:
			return nil
		default:
		}
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			more := c.handleKey(ev.Key, ev.Mod)
			if !more {
				return nil
			}
			for _, input := range inputer {
				input.KeyPress(ev.Key, ev.Mod)
			}
		case termbox.EventMouse:
			for _, input := range inputer {
				input.Mouse(ev)
			}
		case termbox.EventError:
			return ev.Err
		case termbox.EventResize:
			for _, resize := range resizer {
				// TODO(antoine): resize the proper size each children
				resize.Resize(0, 0, ev.Width, ev.Height)
			}
		}
	}
}

func (c *Canvas) handleKey(key termbox.Key, mod termbox.Modifier) bool {
	switch key {
	case termbox.KeyEsc, termbox.KeyCtrlC, termbox.KeyCtrlD:
		return false
	default:
		return true
	}
}

func (c *Canvas) computeIndex(x, y int) int {
	w, _ := termbox.Size()
	return w*y + x
}
