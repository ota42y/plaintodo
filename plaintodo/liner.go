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
	f, err := os.Open(l.automaton.Config.Paths.History)
	if err == nil {
		n, _ := l.ReadHistory(f)
		fmt.Fprintln(l.automaton.Config.Writer, "load", n, "history")
		f.Close()
	}

	for {
		cmd, err := l.Prompt("> ")
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "fatal: %s", err)
			os.Exit(1)
		}

		l.AppendHistory(cmd)
		if l.automaton.Execute(cmd) {
			break
		}
	}

	f, err = os.Create(l.automaton.Config.Paths.History)
	if err == nil {
		fmt.Fprintln(l.automaton.Config.Writer, "write history")
		l.WriteHistory(f)
		f.Close()
	} else {
		fmt.Fprintln(l.automaton.Config.Writer, "write history error:", err)
	}
}
