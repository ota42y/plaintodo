package task

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// I don't know how to math by one regexp...
var blankLineRegexp, _ = regexp.Compile("^( *)$")
var baseRegExpWithAttributes, _ = regexp.Compile("^( *)([^:]+)( :.+)")
var baseRegExpNoAttributes, _ = regexp.Compile("^( *)([^:]+)")

var attributeSplit = " :"
var attributeKeyValueSeparator = " "

var spaceNum = 2 // The Task.Level is task's top space num divide this.

// Task is task struct
// one task line convert to one task struct
type Task struct {
	Level      int
	ID         int
	Name       string
	Attributes map[string]string
	SubTasks   []*Task
}

// Equal is check same task
// This method not check task.ID
func (t *Task) Equal(task *Task) bool {
	if t.Name != task.Name {
		return false
	}

	if t.Level != task.Level {
		return false
	}

	if len(t.Attributes) != len(task.Attributes) {
		return false
	}

	for key, value := range task.Attributes {
		if t.Attributes[key] != value {
			return false
		}
	}

	if len(t.SubTasks) != len(task.SubTasks) {
		return false
	}

	for index, subTask := range task.SubTasks {
		if !t.SubTasks[index].Equal(subTask) {
			return false
		}
	}

	return true
}

// Copy methon copy task, if copySuubTask is true, copy all sub tasks.
// Thin method set taskID to new task's task.ID
func (t *Task) Copy(taskID int, copySubTask bool) *Task {
	attributes := make(map[string]string)
	for key, value := range t.Attributes {
		attributes[key] = value
	}

	var subTasks []*Task
	if copySubTask {
		for index, subTask := range t.SubTasks {
			subTasks = append(subTasks, subTask.Copy(taskID+index+1, copySubTask))
		}
	}

	return &Task{
		Level:      t.Level,
		ID:         taskID,
		Name:       t.Name,
		Attributes: attributes,
		SubTasks:   subTasks,
	}
}

func (t *Task) String(showID bool) string {
	spaces := strings.Repeat(" ", t.Level*spaceNum)

	taskString := make([]string, 1)
	taskString[0] = t.Name

	if showID {
		taskString = append(taskString, fmt.Sprint(":id ", t.ID))
	}

	attributesArray := make([]string, len(t.Attributes))
	i := 0
	for k := range t.Attributes {
		attributesArray[i] = k
		i++
	}
	sort.Strings(attributesArray)

	for _, key := range attributesArray {
		value := t.Attributes[key]

		str := ":" + key
		if value != "" {
			str += " " + value
		}

		taskString = append(taskString, str)
	}

	return spaces + strings.Join(taskString, " ")
}

// RemoveSubTask remove subtask which have specific id
// If removed task have subtask, all subtask removed too.
// This function not delete object
func (t *Task) RemoveSubTask(id int) bool {
	index := -1
	for i, task := range t.SubTasks {
		if task.ID == id {
			index = i
		}
	}

	if index == -1 {
		// task isn't exist
		return false
	}

	t.SubTasks = append(t.SubTasks[:index], t.SubTasks[index+1:]...)
	return true
}

// LoadResult contain task and error line for file loading
type LoadResult struct {
	Tasks     []*Task
	FailLines []string
}

// ParseOptions convert option string to map
// Option string linke ' :key1 value1 :key2 value2' (need first space)
func ParseOptions(raw string) map[string]string {
	options := make(map[string]string)

	// split :key1 value1 :key2 value2 to ["key1 value1", "key2 value2"]
	splits := strings.Split(raw, attributeSplit)
	for _, attribute := range splits {
		if len(attribute) != 0 {
			// split "key1 value1" to "key1", "value1"
			fields := strings.SplitAfterN(attribute, attributeKeyValueSeparator, 2)

			if 0 < len(fields) {
				key := strings.TrimSpace(fields[0])
				value := ""

				if 1 < len(fields) {
					// attribute with value
					value = fields[1]
				}
				options[key] = value
			}
		}
	}
	return options
}

// NewTask create task from task string.
// Returned task seted taskID
func NewTask(line string, taskID int) (*Task, error) {
	b := []byte(line)

	match := blankLineRegexp.FindSubmatch(b)
	if len(match) != 0 {
		return nil, errors.New("blank line")
	}

	match = baseRegExpWithAttributes.FindSubmatch(b)
	if len(match) != 4 {
		match = baseRegExpNoAttributes.FindSubmatch(b)
		if len(match) != 3 {
			return nil, errors.New("parse error")
		}
	}

	task := Task{
		ID: taskID,
	}

	spaces := match[1]
	task.Level = len(spaces) / spaceNum

	task.Name = string(match[2])

	if 3 < len(match) {
		task.Attributes = ParseOptions(string(match[3]))
	} else {
		task.Attributes = make(map[string]string)
	}
	return &task, nil
}

// create subtask under the level.
// return subtasks and next Task (which Task.Level is greater than or same level)
// if nextTask in null, all task read.
func createSubTasks(level int, nowTaskID int, s *bufio.Scanner) (subTasks []*Task, nextTask *Task, maxTaskID int, err error) {
	subTasks = make([]*Task, 0)
	var nowTask *Task
	maxTaskID = nowTaskID

	// read next task or end input
	for s.Scan() {
		line := s.Text()
		nowTask, err = NewTask(line, maxTaskID+1)

		if nowTask != nil {
			maxTaskID = nowTask.ID
			break
		}

		// if blank line, skip this line
		// if not blank line end parse
		if err.Error() != "blank line" {
			return subTasks, nowTask, maxTaskID, err
		}
	}

	for nowTask != nil && level <= nowTask.Level {
		subTasks = append(subTasks, nowTask)

		// get subTasks
		nowTask.SubTasks, nextTask, nowTaskID, err = createSubTasks(nowTask.Level+1, maxTaskID, s)
		if err != nil {
			return subTasks, nowTask, nowTaskID, err
		}

		// if get smaller level task, createSubTasks end
		// if get same level task, create next subtask
		// createSubTasks don't return greater level task
		if nextTask != nil {
			if nextTask.Level < level {
				return subTasks, nextTask, nowTaskID, nil
			}
		}

		maxTaskID = nowTaskID
		nowTask = nextTask
	}
	return subTasks, nowTask, maxTaskID, nil
}

func createTasks(s *bufio.Scanner) ([]*Task, int) {
	taskID := 0
	topLevelTasks, nextTask, maxTaskID, err := createSubTasks(0, taskID, s)

	if err != nil {
		panic(err)
	}

	if nextTask != nil {
		panic("create Task error, there is -1 or smaller task level exist")
	}

	return topLevelTasks, maxTaskID
}

// ReadTasks read tasks from file.
func ReadTasks(filename string) ([]*Task, int) {
	var fp *os.File
	var err error

	fp, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	tasks, maxTaskID := createTasks(scanner)
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return tasks, maxTaskID
}

// GetTask return specific id's task from task array
func GetTask(id int, tasks []*Task) (parent *Task, task *Task) {
	for _, task := range tasks {
		if task.ID == id {
			return nil, task
		}

		p, t := GetTask(id, task.SubTasks)
		if t != nil {
			if p == nil {
				p = task
			}
			return p, t
		}
	}

	return nil, nil
}
