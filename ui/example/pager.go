package main

import (
	"flag"
	"github.com/aybabtme/logterm/ui"
	"github.com/aybabtme/tailf"
	"io"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	f, err := os.OpenFile("canvas.log1", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("can't create log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Print("starting")

	follow := flag.String("f", "", "file to follow")
	flag.Parse()

	if *follow == "" {
		log.Print("need a file to follow")
		flag.PrintDefaults()
		return
	}

	src, err := tailf.Follow(*follow, true)
	if err != nil {
		log.Fatalf("couldn't follow %q: %v", *follow, err)
	}
	defer src.Close()

	c, err := ui.NewCanvas(60)
	if err != nil {
		log.Fatalf("couldn't create canvas: %v", err)
	}
	defer c.Close()

	win := c.Fullscreen()
	pager := ui.NewPagerBox(win)

	go func() {
		n, err := io.Copy(pager, src)
		if err != nil {
			log.Fatalf("writing to pager: %v", err)
		}
		log.Printf("%d bytes written", n)
	}()

	err = c.Run([]ui.ResizeHandler{win}, nil)
	if err != nil {
		log.Printf("error running canvas: %v", err)
	}

}
