package genlogs

import (
	"github.com/Sirupsen/logrus"
	"github.com/aybabtme/gypsum"
	"github.com/dustin/randbo"
	"io"
	"log"
	"math/rand"
	"time"
)

type RandLogger struct {
	logjson *logrus.Logger
	logfmt  *logrus.Logger
	lograw  *log.Logger

	onlyJSON bool
	onlyFmt  bool
	onlyRaw  bool
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

func NewJSONLogger(out io.Writer) *RandLogger {
	rl := NewRandLogger(out)
	rl.onlyJSON = true
	return rl
}

func NewFmtLogger(out io.Writer) *RandLogger {
	rl := NewRandLogger(out)
	rl.onlyFmt = true
	return rl
}

func NewRawLogger(out io.Writer) *RandLogger {
	rl := NewRandLogger(out)
	rl.onlyRaw = true
	return rl
}

func (r *RandLogger) Log() {
	switch {
	case r.onlyJSON:
		r.logrusRandLvl(r.logjson, randomFields(false))
		return
	case r.onlyFmt:
		r.logrusRandLvl(r.logfmt, randomFields(false))
		return
	case r.onlyRaw:
		r.logRawRandLvl(randomFields(true))
		return
	}

	switch i := rand.Intn(3); i {
	case 0:
		r.logrusRandLvl(r.logfmt, randomFields(false))
	case 1:
		r.logrusRandLvl(r.logjson, randomFields(true))
	case 2:
		r.logRawRandLvl(randomFields(true))
	default:
		panic(i)
	}
}

func (r *RandLogger) logRawRandLvl(fields map[string]interface{}) {
	switch i := rand.Intn(5); i {
	case 0:
		r.lograw.Printf("debug: %#v", fields)
	case 1:
		r.lograw.Printf("info: %#v", fields)
	case 2:
		r.lograw.Printf("warn: %#v", fields)
	case 3:
		r.lograw.Printf("error: %#v", fields)
	case 4:
		r.lograw.Printf("fatal: %#v", fields)
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
