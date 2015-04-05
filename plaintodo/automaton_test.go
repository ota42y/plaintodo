package main

import "testing"

type CommandTest struct {
	T         *testing.T
	Option    string
	Called    bool
	Terminate bool
}

func (t CommandTest) Execute(option string, automaton *Automaton) (terminate bool) {
	if option != t.Option {
		t.T.Errorf("option shud be %s but %s", t.Option, option)
		t.T.FailNow()
	}

	return t.Terminate
}

func TestAutomatonWithOption(t *testing.T) {
	cmd := CommandTest{
		T:         t,
		Option:    "option test",
		Called:    false,
		Terminate: false,
	}

	cmds := make(map[string]Command)
	cmds["test"] = cmd

	a := NewAutomaton(nil, cmds)

	terminate := a.Execute("test " + cmd.Option)

	if !cmd.Called {
		t.Errorf("command not called")
		t.FailNow()
	}

	if terminate != cmd.Terminate {
		t.Errorf("Automation.Execute shud be return %v but %v", terminate, cmd.Terminate)
		t.FailNow()
	}
}

func TestAutomaton(t *testing.T) {
	cmd := CommandTest{
		T:         t,
		Option:    "",
		Called:    false,
		Terminate: true,
	}

	cmds := make(map[string]Command)
	cmds["test"] = cmd

	a := NewAutomaton(nil, cmds)

	terminate := a.Execute("test " + cmd.Option)

	if !cmd.Called {
		t.Errorf("command not called")
		t.FailNow()
	}

	if terminate != cmd.Terminate {
		t.Errorf("Automation.Execute shud be return %v but %v", terminate, cmd.Terminate)
		t.FailNow()
	}
}
