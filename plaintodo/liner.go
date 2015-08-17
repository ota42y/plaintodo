package main

import (
	"fmt"
	"github.com/peterh/liner"
	"io"
	"os"

	"./command"
	"./config"
	"./executor"
)

type Liner struct {
	*liner.State
	e *executor.Executor
}

func NewLiner(config *config.Config, commands map[string]command.Command) *Liner {
	s := liner.NewLiner()
	e := executor.NewExecutor(config, commands)
	return &Liner{
		State: s,
		e:     e,
	}
}

func (l *Liner) Start() {
	f, err := os.Open(l.e.S.Config.Liner.History)
	if err == nil {
		n, _ := l.ReadHistory(f)
		fmt.Fprintln(l.e.S.Config.Writer, "load", n, "history")
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
		if l.e.Execute(cmd) {
			break
		}
	}

	f, err = os.Create(l.e.S.Config.Liner.History)
	if err == nil {
		fmt.Fprintln(l.e.S.Config.Writer, "write history")
		l.WriteHistory(f)
		f.Close()
	} else {
		fmt.Fprintln(l.e.S.Config.Writer, "write history error:", err)
	}
}
