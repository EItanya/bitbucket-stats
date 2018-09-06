package gui

import (
	"fmt"
	"time"
)

var (
	sequence = []string{"", ".", "..", "..."}
	text     = "Loading"
)

// Loader basic gui loader
type Loader struct {
	Text     string
	Sequence []string
	done     bool
}

func (l *Loader) loading() {
	i := 0
	for {
		if l.done {
			break
		} else {
			fmt.Printf("\r%s%s ", l.Text, l.Sequence[i])
			if i == len(l.Sequence)-1 {
				i = 0
			} else {
				i++
			}
		}
		time.Sleep(time.Millisecond * 125)
	}
}

// Start the loader
func (l *Loader) Start() {
	if l.Sequence == nil {
		l.Sequence = sequence
	}
	if l.Text == "" {
		l.Text = text
	}
	l.done = false
	go l.loading()
}

// Stop the loader
func (l *Loader) Stop() {
	l.done = true
	fmt.Printf("\rAll done\n")
}
