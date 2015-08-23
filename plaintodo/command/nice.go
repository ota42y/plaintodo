package command

import (
	"fmt"
	"regexp"

	"../task"
	"../util"
)

var evernoteRegexp, _ = regexp.Compile("^https://www.evernote.com/shard/(.+)/nl/(.+)/([^/]+)")

// Nice change some effect to task
type Nice struct {
}

func (c *Nice) fixEvernoteURL(tasks []*task.Task) int {
	count := 0
	for _, task := range tasks {
		match := evernoteRegexp.FindSubmatch([]byte(task.Attributes["url"]))
		if len(match) == 4 {
			task.Attributes["url"] = fmt.Sprintf("evernote:///view/%s/%s/%s/%s/", match[2], match[1], match[3], match[3])
			count++
		}
		count += c.fixEvernoteURL(task.SubTasks)
	}
	return count
}

// Execute do nice
func (c *Nice) Execute(option string, s *State) (terminate bool) {
	var tasks []*task.Task

	optionMap := task.ParseOptions(" " + option)
	id, err := util.GetIntAttribute("id", optionMap)
	if err != nil {
		// do all tasks
		tasks = s.Tasks
	} else {
		// do selected task
		_, t := task.GetTask(id, s.Tasks)
		tasks = make([]*task.Task, 1)
		tasks[0] = t
	}

	fmt.Fprintf(s.Config.Writer, "Done nice\n")

	num := c.fixEvernoteURL(tasks)
	fmt.Fprintf(s.Config.Writer, "evernote url change %d tasks\n", num)

	return false
}

// NewNice return Nice
func NewNice() *Nice {
	return &Nice{}
}
