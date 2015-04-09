package main

import (
	"fmt"
	"github.com/peterh/liner"
	"io"
	"os"
)

type Liner struct {
	*liner.State
	automaton *Automaton
}

func NewLiner(config *Config, commands map[string]Command) *Liner {
	s := liner.NewLiner()
	a := NewAutomaton(config, commands)
	return &Liner{
		State:     s,
		automaton: a,
	}
}

func (l *Liner) Start() {
	for {
		cmd, err := l.Prompt("> ")
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Fprintf(os.Stderr, "fatal: %s", err)
			os.Exit(1)
		}

		if l.automaton.Execute(cmd) {
			return
		}
	}
}
