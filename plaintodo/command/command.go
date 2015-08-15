package command

// Command is command interface
type Command interface {
	Execute(option string, s *State) (terminate bool)
}
