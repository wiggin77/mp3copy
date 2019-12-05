package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	PROGRESS_BAR     = "********************"
	PROGRESS_BAR_LEN = len(PROGRESS_BAR)
)

var (
	Term = &Terminal{sb: strings.Builder{}}
)

// Terminal provides a simple way to write to the terminal, providing
// a progress bar and error output.
type Terminal struct {
	progress float32 // percent complete
	sb       strings.Builder
}

func (t *Terminal) Printf(format string, a ...interface{}) {
	t.printf(os.Stdout, format, a)
}

func (t *Terminal) Errorf(format string, a ...interface{}) {
	t.printf(os.Stderr, format, a)
}

func (t *Terminal) printf(w io.Writer, format string, a ...interface{}) {
	t.clearLine()
	fmt.Fprintf(os.Stdout, format, a)
	t.displayProgress()
}

func (t *Terminal) clearLine() {
	fmt.Fprintf(os.Stdout, "\033[2K\r")
}

func (t *Terminal) displayProgress() {
	prog :=

		t.sb.Reset()
	for i := 0; i < PROGRESS_WIDTH; i++ {
		if float32(i) <= PROGRESS_WIDTH*t.progress {
			t.sb.WriteString("*")
		} else {
			t.sb.WriteString(" ")
		}
	}
	fmt.Fprintf(os.Stdout, "\r[%s]%d%%", t.sb.String(), t.progress)
}
