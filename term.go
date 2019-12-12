package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

const (
	PROGRESS_BAR     = "******************************"
	PROGRESS_SPACE   = "                              "
	PROGRESS_BAR_LEN = len(PROGRESS_BAR)
)

var (
	Term = &Terminal{sb: strings.Builder{}, out: ioutil.Discard, err: ioutil.Discard}
)

// Terminal provides a simple way to write to the terminal, providing
// a progress bar and error output.
type Terminal struct {
	progress float32 // percent complete
	sb       strings.Builder
	out      io.Writer
	err      io.Writer
}

func (t *Terminal) SetOut(w io.Writer) {
	if w == nil {
		w = ioutil.Discard
	}
	t.out = w
}

func (t *Terminal) SetErr(w io.Writer) {
	if w == nil {
		w = ioutil.Discard
	}
	t.err = w
}

func (t *Terminal) Printf(format string, a ...interface{}) {
	t.printf(t.out, format, a...)
}

func (t *Terminal) Errorf(format string, a ...interface{}) {
	t.printf(t.err, format, a...)
}

func (t *Terminal) Progress(current int64, total int64) {
	if total == 0 {
		t.progress = 0
		return
	}
	t.progress = float32(current) / float32(total)
	t.displayProgress()
}

func (t *Terminal) printf(w io.Writer, format string, a ...interface{}) {
	t.clearLine()
	fmt.Fprintf(t.out, format, a...)
	t.displayProgress()
}

func (t *Terminal) clearLine() {
	fmt.Fprintf(t.out, "\033[2K\r")
}

func (t *Terminal) displayProgress() {
	prog := int(float32(PROGRESS_BAR_LEN)*t.progress + 0.5)
	if prog > 0 {
		bar := PROGRESS_BAR[:prog]
		space := PROGRESS_SPACE[:PROGRESS_BAR_LEN-prog]
		fmt.Fprintf(t.out, "\r[%s%s] %.1f%%  ", bar, space, t.progress*100)
	}
}
