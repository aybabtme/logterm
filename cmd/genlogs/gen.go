package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/aybabtme/gypsum"
	"github.com/aybabtme/iocontrol"
	"github.com/dustin/go-humanize"
	"github.com/dustin/randbo"
	"io"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

func main() {
	lines := flag.Uint("lines", 100, "lines of log output to produce")
	output := flag.String("output", "", "where to write the output, default to stdout")
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
	logger := NewRandLogger(measure)

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

type RandLogger struct {
	logjson *logrus.Logger
	logfmt  *logrus.Logger
	lograw  *log.Logger
}

func NewRandLogger(out io.Writer) *RandLogger {
	logfmt := logrus.New()
	logfmt.Formatter = &logrus.TextFormatter{
		DisableColors: true,
	}
	logfmt.Out = out

	logjson := logrus.New()
	logjson.Formatter = &logrus.JSONFormatter{}
	logjson.Out = out

	return &RandLogger{
		logjson: logjson,
		logfmt:  logfmt,
		lograw:  log.New(out, "[raw log]", log.LstdFlags),
	}
}

func (r *RandLogger) Log() {
	switch i := rand.Intn(3); i {
	case 0:
		r.logrusRandLvl(r.logfmt, randomFields(false))
	case 1:
		r.logrusRandLvl(r.logjson, randomFields(true))
	case 2:
		switch i := rand.Intn(5); i {
		case 0:
			r.lograw.Printf("debug: %#v", randomFields(true))
		case 1:
			r.lograw.Printf("info: %#v", randomFields(true))
		case 2:
			r.lograw.Printf("warn: %#v", randomFields(true))
		case 3:
			r.lograw.Printf("error: %#v", randomFields(true))
		case 4:
			r.lograw.Printf("fatal: %#v", randomFields(true))
		default:
			panic(i)
		}
	default:
		panic(i)
	}
}

func (r *RandLogger) logrusRandLvl(l *logrus.Logger, fields map[string]interface{}) {
	e := l.WithFields(logrus.Fields(fields))
	word := gypsum.WordLorem(rand.Intn(10))
	switch i := rand.Intn(4); i {
	case 0:
		e.Debug(word)
	case 1:
		e.Info(word)
	case 2:
		e.Warn(word)
	case 3:
		e.Error(word)
	default:
		panic(i)
	}
}

var randbytes = randbo.New()

func randomField(includeBytes bool) (string, interface{}) {
	var val interface{}
	// this is pretty shitty
	switch i := rand.Intn(7); i {
	case 0: // field []byte
		if includeBytes {
			b := make([]byte, rand.Intn(1<<8)+1)
			_, _ = randbytes.Read(b)
			val = b
			break
		}
		fallthrough
	case 1: // field string
		val = gypsum.Lorem()
	case 2: // field number
		val = rand.Float64()
	case 3: // field duration
		val = time.Duration(rand.Float64()) * time.Second * 10
	case 4: // field time
		val = time.Date(
			rand.Intn(2000)+1,
			time.Month(rand.Intn(12)+1),
			rand.Intn(28)+1,
			rand.Intn(24),
			rand.Intn(60),
			rand.Intn(60),
			rand.Intn(1000),
			time.UTC,
		)
	case 5: // field boolean
		val = rand.Float64() > 0.5
	case 6:
		val = nil
	default:
		panic(i)
	}
	return gypsum.WordLorem(1), val
}

func randomFields(includeBytes bool) map[string]interface{} {
	count := rand.Intn(20)
	fields := make(map[string]interface{}, count)
	for i := 0; i < count; i++ {
		name, val := randomField(includeBytes)
		fields[name] = val
	}
	return fields
}
