package main

import (
	"bufio"
	"errors"
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

type Task struct {
	Level      int
	Id         int
	Name       string
	Attributes map[string]string
	SubTasks   []*Task
}

func (t *Task) String(showId bool) string {
	spaces := strings.Repeat(" ", t.Level*spaceNum)

	taskString := make([]string, 1)
	taskString[0] = t.Name

	attributesArray := make([]string, len(t.Attributes))
	i := 0
	for k, _ := range t.Attributes {
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

type LoadResult struct {
	Tasks     []*Task
	FailLines []string
}

func getAttributes(raw string) map[string]string {
	attributes := make(map[string]string)

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
				attributes[key] = value
			}
		}
	}
	return attributes
}

func NewTask(line string, taskId int) (*Task, error) {
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
		Id: taskId,
	}

	spaces := match[1]
	task.Level = len(spaces) / spaceNum

	task.Name = string(match[2])

	if 3 < len(match) {
		task.Attributes = getAttributes(string(match[3]))
	}
	return &task, nil
}

// create subtask under the level.
// return subtasks and next Task (which Task.Level is greater than or same level)
// if nextTask in null, all task read.
func createSubTasks(level int, taskId int, s *bufio.Scanner) (subTasks []*Task, nextTask *Task, err error) {
	subTasks = make([]*Task, 0)
	var nowTask *Task = nil

	// read next task or end input
	for s.Scan() {
		line := s.Text()
		nowTask, err = NewTask(line, taskId)

		if nowTask != nil {
			break
		}

		// if blank line, skip this line
		// if not blank line end parse
		if err.Error() != "blank line" {
			return subTasks, nowTask, err
		}
	}

	for nowTask != nil && level <= nowTask.Level {
		subTasks = append(subTasks, nowTask)

		// get subTasks
		nowTask.SubTasks, nextTask, err = createSubTasks(nowTask.Level+1, nowTask.Id+1, s)
		if err != nil {
			return subTasks, nowTask, err
		}

		// if get smaller level task, createSubTasks end
		// if get same level task, create next subtask
		// createSubTasks don't return greater level task
		if nextTask != nil {
			if nextTask.Level < level {
				return subTasks, nextTask, nil
			}
		}

		nowTask = nextTask
	}

	return subTasks, nowTask, nil
}

func createTasks(s *bufio.Scanner) []*Task {
	taskId := 1
	topLevelTasks, nextTask, err := createSubTasks(0, taskId, s)

	if err != nil {
		panic(err)
	}

	if nextTask != nil {
		panic("create Task error, there is -1 or smaller task level exist")
	}

	return topLevelTasks
}

func ReadTasks(filename string) []*Task {
	var fp *os.File
	var err error

	fp, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	tasks := createTasks(scanner)
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return tasks
}
