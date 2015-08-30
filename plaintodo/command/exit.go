package command

// Exit exit this application
type Exit struct {
}

// Execute do nothing, return true
func (t *Exit) Execute(option string, s *State) (terminate bool) {
	return true
}

// NewExit return Exit
func NewExit() *Exit {
	return &Exit{}
}
