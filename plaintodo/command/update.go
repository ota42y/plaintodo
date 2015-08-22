package command

// Update archive and save and reload tasks
type Update struct {
	archive *Archive
	save    *Save
	reload  *Reload
}

// Execute execute update
func (u *Update) Execute(option string, s *State) (terminate bool) {
	u.archive.Execute(option, s)
	u.save.Execute(option, s)
	u.reload.Execute(option, s)
	return false
}

// NewUpdate return Update
func NewUpdate() *Update {
	return &Update{
		archive: NewArchive(),
		save:    NewSave(),
		reload:  NewReload(),
	}
}
