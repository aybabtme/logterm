package main

import (
	"code.google.com/p/go.crypto/ssh/terminal"
	"errors"
	"flag"
	"fmt"
	"github.com/aybabtme/iocontrol"
	"github.com/aybabtme/tailf"
	"github.com/dustin/go-humanize"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const prompt = "humanlog> "

var debugComplete = func(line string, pos int, key rune) (string, int, bool) {
	var out string
	var handled bool
	switch key {
	case '\t':
		log.Printf("tab!")
		return out, pos, true
	case ETX:
		log.Printf("quittin!")
		os.Exit(0)
	default:
	}
	log.Printf("autocomplete (%d) %q: %q", pos, key, line)
	return out, pos, handled
}

func main() {
	log.SetFlags(0)
	tui := flag.Bool("tui", false, "run as an interactive terminal interface")
	follow := flag.String("f", "", "file to follow")
	tail := flag.Bool("tail", false, "when following a file, don't first read the whole file's content (similar to `tail -f`)")
	flag.Parse()

	var src io.Reader
	var err error
	if *follow != "" {
		fsrc, err := tailf.Follow(*follow, !*tail)
		if err != nil {
			log.Fatalf("can't follow file %q, %v", *follow, err)
		}
		defer fsrc.Close()
		src = fsrc
	} else if *tui {
		src, err = followCommand(flag.Args())
		if err != nil {
			log.Printf("can't read output of command %q, %v", strings.Join(flag.Args(), " "), err)
			log.Fatal("no file to follow, need a command or a file to follow when in interactive mode")
		}
	} else {
		src = os.Stdin
	}

	var out io.Writer
	if *tui {
		term, err := startTUI(debugComplete, func(line string) error {
			return nil
		})
		if err != nil {
			log.Fatalf("error with interactive mode: %v", err)
		}
		measured := iocontrol.NewMeasuredReader(src)
		src = measured
		out = term
		go func() {
			for _ = range time.Tick(time.Second) {
				persec := measured.BytesPerSec()
				term.SetPrompt(fmt.Sprintf("%vps: %s", humanize.Bytes(persec), prompt))
			}
		}()

	} else {
		out = os.Stdout
	}

	_, err = io.Copy(out, src)
	if err != nil {
		log.Fatalf("error with input source: %v", err)
	}
}

func followCommand(args []string) (io.Reader, error) {
	if len(args) < 1 {
		return nil, errors.New("need a command to run in interactive mode")
	}
	rd, combinedOut := io.Pipe()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = combinedOut
	cmd.Stderr = combinedOut
	if err := cmd.Start(); err != nil {
		return rd, err
	}
	go func() {
		if err := cmd.Wait(); err != nil {
			combinedOut.CloseWithError(err)
		} else {
			combinedOut.Close()
		}
	}()
	return rd, nil
}

type OnAutocomplete func(line string, pos int, key rune) (string, int, bool)
type OnReadline func(line string) error

func startTUI(tabComplete OnAutocomplete, onReadline OnReadline) (*terminal.Terminal, error) {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return nil, err
	}

	term := terminal.NewTerminal(os.Stdin, prompt)
	term.AutoCompleteCallback = tabComplete

	log.SetOutput(term)

	go func() {
		defer terminal.Restore(0, oldState)
		for {
			line, err := term.ReadLine()
			switch err {
			case io.EOF:
				return
			default:
				panic(err)
			case nil:
				err = onReadline(line)
				log.Fatalf("couldn't use callback: %v", err)
			}
		}
	}()

	return term, err
}

var (
	NUL = rune(0x00)
	SOH = rune(0x01)
	STX = rune(0x02)
	ETX = rune(0x03)
	EOT = rune(0x04)
	ENQ = rune(0x05)
	ACK = rune(0x06)
	BEL = rune(0x07)
	BS  = rune(0x08)
	TAB = rune(0x09)
	LF  = rune(0x0A)
	VT  = rune(0x0B)
	FF  = rune(0x0C)
	CR  = rune(0x0D)
	SO  = rune(0x0E)
	SI  = rune(0x0F)
	DEL = rune(0x7F)
)
