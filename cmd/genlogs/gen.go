package main

import (
	"flag"
	"github.com/aybabtme/iocontrol"
	"github.com/aybabtme/logterm/testutil/genlogs"
	"github.com/dustin/go-humanize"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	lines := flag.Uint("lines", 100, "lines of log output to produce")
	output := flag.String("output", "", "where to write the output, default to stdout")
	onlyJSON := flag.Bool("only-json", false, "only log random JSON objects")
	onlyFmt := flag.Bool("only-logfmt", false, "only log random logfmt'd strings")
	onlyRaw := flag.Bool("only-raw", false, "only log random strings that aren't JSON or logfmt")
	flag.Parse()

	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	var out io.Writer
	if *output != "" {
		f, err := os.OpenFile(*output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("couldn't open or create %q, %v", *output, err)
		}
		defer f.Close()
		out = f
	} else {
		out = os.Stdout
	}
	measure := iocontrol.NewMeasuredWriter(out)

	var logger *genlogs.RandLogger
	switch {
	case *onlyJSON:
		logger = genlogs.NewJSONLogger(measure)
	case *onlyFmt:
		logger = genlogs.NewFmtLogger(measure)
	case *onlyRaw:
		logger = genlogs.NewRawLogger(measure)
	default:
		logger = genlogs.NewRandLogger(measure)
	}

	var mu sync.Mutex
	var i uint

	go func() {
		ticker := time.NewTicker(time.Second)
		for _ = range ticker.C {
			mu.Lock()
			log.Printf("lines=%d, rate=%sps", i, humanize.Bytes(measure.BytesPerSec()))
			mu.Unlock()
		}
	}()

	for {
		mu.Lock()
		if i >= *lines {
			log.Printf("lines=%d, rate=%sps", i, humanize.Bytes(measure.BytesPerSec()))
			mu.Unlock()
			return
		}
		i++
		mu.Unlock()

		logger.Log()
	}

}
