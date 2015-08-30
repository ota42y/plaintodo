package command

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"../util"
)

func TestExitCommand(t *testing.T) {
	cmd := NewExit()
	config, _ := util.ReadTestConfigRelativePath("..")
	s := &State{
		Config: config,
	}

	Convey("correct", t, func() {
		terminate := cmd.Execute("", s)
		So(terminate, ShouldBeTrue)
	})
}
