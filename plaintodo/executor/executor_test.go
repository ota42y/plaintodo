package executor

import (
	"bytes"
	"testing"

	"../command"
	"../config"
)

type CommandTest struct {
	T         *testing.T
	Option    string
	Called    bool
	Terminate bool
}

func (t *CommandTest) Execute(option string, s *command.State) (terminate bool) {
	t.Called = true

	if option != t.Option {
		t.T.Errorf("option shud be %s but %s", t.Option, option)
		t.T.FailNow()
	}

	return t.Terminate
}

func TestAutomatonWithOption(t *testing.T) {
	cmd := &CommandTest{
		T:         t,
		Option:    "option test",
		Called:    false,
		Terminate: false,
	}

	cmds := make(map[string]command.Command)
	cmds["test"] = cmd

	buf := &bytes.Buffer{}

	config := config.ReadConfig("../test_config.toml")
	config.Writer = buf
	e := NewExecutor(config, cmds)

	terminate := e.Execute("test " + cmd.Option)

	if !cmd.Called {
		t.Errorf("command not called")
		t.FailNow()
	}

	if terminate != cmd.Terminate {
		t.Errorf("Automation.Execute shud be return %v but %v", terminate, cmd.Terminate)
		t.FailNow()
	}

	output := buf.String()
	if output != "test hit\n" {
		t.Errorf("Invalid output '%s'", output)
		t.FailNow()
	}
}

func TestAutomaton(t *testing.T) {
	cmd := &CommandTest{
		T:         t,
		Option:    "",
		Called:    false,
		Terminate: true,
	}

	cmds := make(map[string]command.Command)
	cmds["test"] = cmd

	e := NewExecutor(nil, cmds)

	terminate := e.Execute("test " + cmd.Option)

	if !cmd.Called {
		t.Errorf("command not called")
		t.FailNow()
	}

	if terminate != cmd.Terminate {
		t.Errorf("Automation.Execute shud be return %v but %v", terminate, cmd.Terminate)
		t.FailNow()
	}
}

func TestAutomatonWithInvalidCommand(t *testing.T) {
	cmd := &CommandTest{
		T:         t,
		Option:    "",
		Called:    false,
		Terminate: true,
	}

	cmds := make(map[string]command.Command)
	cmds["test"] = cmd

	buf := &bytes.Buffer{}

	config := config.ReadConfig("../test_config.toml")
	config.Writer = buf
	e := NewExecutor(config, cmds)

	terminate := e.Execute("no " + cmd.Option)

	if cmd.Called {
		t.Errorf("command called")
		t.FailNow()
	}

	if terminate {
		t.Errorf("If no command hit, shuld return false but true")
		t.FailNow()
	}

	output := buf.String()
	if output != "no not hit\n" {
		t.Errorf("Invalid output '%s'", output)
		t.FailNow()
	}
}
