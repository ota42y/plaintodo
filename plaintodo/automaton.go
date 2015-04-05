package main

type Command interface {
	Execute(option string, automaton *Automaton) (terminate bool)
}

type Automaton struct {
	tasks    []*Task
	commands map[string]Command
	config   *Config // main.go
}

func NewAutomaton(config *Config, commands map[string]Command) *Automaton {
	return &Automaton{
		tasks:    make([]*Task, 0),
		commands: commands,
		config:   config,
	}
}

// cmd shuld be "cmd options"
func (a *Automaton) Execute(cmd string) (terminate bool) {

	return true
}
