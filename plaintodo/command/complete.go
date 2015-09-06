package command

import (
	"fmt"
	"strings"
	"time"

	"../task"
	"../util"
)

// Complete complete task and new repeated task
type Complete struct {
	MaxTaskID int
}

type completeResult struct {
	completeTask *task.Task
	newTasks     []*task.Task
	completeNum  int
}

func (c *Complete) setNewRepeat(baseTime time.Time, task *task.Task) {
	repeatString, ok := task.Attributes["repeat"]
	if !ok {
		return
	}

	// every 1 day
	splits := strings.Split(repeatString, " ")
	if len(splits) != 3 {
		return
	}

	if splits[0] == "every" {
		startString, ok := task.Attributes["start"]
		if !ok {
			return
		}
		baseTime, ok = util.ParseTime(startString)
		if !ok {
			return
		}
	}

	newTime := util.AddDuration(baseTime, splits[1], splits[2])
	task.Attributes["start"] = newTime.Format(util.DateTimeFormat)
}

func (c *Complete) completeAllSubTask(completeDate time.Time, t *task.Task) (repeatTask *task.Task, completeNum int) {
	n := 0
	var newSubTasks []*(task.Task)

	for _, subTask := range t.SubTasks {
		repeatSubTask, num := c.completeAllSubTask(completeDate, subTask)
		n += num
		if repeatSubTask != nil {
			newSubTasks = append(newSubTasks, repeatSubTask)
		}
	}

	_, ok := t.Attributes["repeat"]
	if len(newSubTasks) != 0 || ok {
		repeatTask = t.Copy(c.MaxTaskID+1, false)
		delete(repeatTask.Attributes, "postpone")
		repeatTask.SubTasks = newSubTasks
		c.setNewRepeat(completeDate, repeatTask)
		c.MaxTaskID++
	}

	_, ok = t.Attributes["complete"]
	if !ok {
		// if not completed, set complete date
		t.Attributes["complete"] = completeDate.Format(util.DateTimeFormat)
		n++
	}

	return repeatTask, n
}

func (c *Complete) completeTask(taskID int, tasks []*task.Task) (*completeResult, error) {
	for _, task := range tasks {
		if task.ID == taskID {
			_, ok := task.Attributes["lock"]
			if ok {
				// if lock task, not change
				return nil, fmt.Errorf("task :id %d is locked", taskID)
			}

			repeatTask, n := c.completeAllSubTask(time.Now(), task)
			if repeatTask != nil {
				tasks = append(tasks, repeatTask)
			}

			result := &completeResult{
				completeNum:  n,
				completeTask: task,
				newTasks:     tasks,
			}

			return result, nil
		}

		result, err := c.completeTask(taskID, task.SubTasks)
		if err == nil {
			task.SubTasks = result.newTasks
			result.newTasks = tasks
			return result, err
		}
	}
	return nil, fmt.Errorf("There is no Task which have task id: %d\n", taskID)
}

// Execute complete task and if set repeat attribute, create new task
func (c *Complete) Execute(option string, s *State) (terminate bool) {
	c.MaxTaskID = s.MaxTaskID

	optionMap := task.ParseOptions(" " + option)

	id, err := util.GetIntAttribute("id", optionMap)
	if err != nil {
		s.Config.Writer.Write([]byte(err.Error()))
		return false
	}

	result, err := c.completeTask(id, s.Tasks)
	if err != nil {
		fmt.Fprintf(s.Config.Writer, err.Error())
		return false
	}

	s.Tasks = result.newTasks
	s.MaxTaskID = c.MaxTaskID
	fmt.Fprintf(s.Config.Writer, "Complete %s and %d sub tasks\n", result.completeTask.Name, result.completeNum-1)
	return false
}

// NewComplete return Complete
func NewComplete() *Complete {
	return &Complete{}
}
