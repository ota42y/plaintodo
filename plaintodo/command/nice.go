package command

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"../task"
	"../util"
)

var evernoteRegexp, _ = regexp.Compile("^https://www.evernote.com/shard/(.+)/nl/(.+)/([^/]+)")

// Nice change some effect to task
type Nice struct {
}

func (c *Nice) fixDateInKey(task *task.Task, key string, today time.Time) bool {
	tomorrow := today.AddDate(0, 0, 1)

	var isReplaced = false
	if strings.Contains(task.Attributes[key], "now") {
		task.Attributes[key] = strings.Replace(task.Attributes[key], "now", today.Format(util.DateTimeFormat), -1)
		isReplaced = true
	}

	if strings.Contains(task.Attributes[key], "today") {
		task.Attributes[key] = strings.Replace(task.Attributes[key], "today", today.Format(util.DateFormat), -1)
		isReplaced = true
	}

	if strings.Contains(task.Attributes[key], "tomorrow") {
		task.Attributes[key] = strings.Replace(task.Attributes[key], "tomorrow", tomorrow.Format(util.DateFormat), -1)
		isReplaced = true
	}

	return isReplaced
}

func (c *Nice) fixDate(tasks []*task.Task) int {
	today := time.Now()

	count := 0
	for _, task := range tasks {
		// if not locked do nice
		_, ok := task.Attributes["lock"]
		if !ok {
			if c.fixDateInKey(task, "start", today) || c.fixDateInKey(task, "postpone", today) {
				count++
			}
		}
		count += c.fixDate(task.SubTasks)
	}
	return count
}

func (c *Nice) fixEvernoteURL(tasks []*task.Task) int {
	count := 0
	for _, task := range tasks {

		// if not locked do nice
		_, ok := task.Attributes["lock"]
		if !ok {
			match := evernoteRegexp.FindSubmatch([]byte(task.Attributes["url"]))
			if len(match) == 4 {
				task.Attributes["url"] = fmt.Sprintf("evernote:///view/%s/%s/%s/%s/", match[2], match[1], match[3], match[3])
				count++
			}
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

	num = c.fixDate(tasks)
	fmt.Fprintf(s.Config.Writer, "change %d tasks date\n", num)

	return false
}

// NewNice return Nice
func NewNice() *Nice {
	return &Nice{}
}
