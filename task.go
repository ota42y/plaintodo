package main

import(
	"bufio"
	"os"
	"regexp"
	"errors"
)

// I don't know how to math by one regexp...
var blankLineRegxpp, _ = regexp.Compile("^( *)$")
var baseRegExpWithAttributes, _ = regexp.Compile("^( *)([^:]+)( :.+)")
var baseRegExpNoAttributes, _ = regexp.Compile("^( *)([^:]+)")

var spaceNum = 2 // The Task.Level is task's top space num divide this.

type Task struct {
	Level int
	Id int
	Name string
	Attribute map[string]string
	SubTasks []*Task
}

type LoadResult struct{
	Tasks []*Task
	FailLines []string
}

func NewTask(line string) (*Task, error){
	b := []byte(line)

	match := baseRegExpWithAttributes.FindSubmatch(b)
	if len(match) != 4 {
		match = baseRegExpNoAttributes.FindSubmatch(b)
		if len(match) != 3{
			match = blankLineRegxpp.FindSubmatch(b)
			if len(match) != 0{
				return nil, errors.New("blank line")
			}else{
				return nil, errors.New("parse error")
			}
		}
	}

	spaces := match[1]
	level := len(spaces) / spaceNum

	name := string(match[2])
	//attributes := match[3]

	return &Task{
		Name: name,
		Level: level,
	}, nil
}

// create subtask under the level.
// return subtasks and next Task (which Task.Level is greater than or same level)
// if nextTask in null, all task read.
func createSubTasks(level int, s *bufio.Scanner) (subTasks []*Task , nextTask *Task, err error){
	subTasks = make([]*Task, 0)
	var nowTask *Task = nil

	// read next task or end input
	for s.Scan() {
		line := s.Text()
		nowTask, err = NewTask(line)

		if nowTask != nil {
			break
		}

		// if blank line, skip this line
		// if not blank line end parse
		if err.Error() != "blank line"{
			return subTasks, nowTask, err
		}
	}

	for nowTask != nil && level <= nowTask.Level{
		subTasks = append(subTasks, nowTask)

		// get subTasks
		nowTask.SubTasks, nextTask, err = createSubTasks(nowTask.Level + 1, s)
		if err != nil{
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

func createTasks(s *bufio.Scanner) []*Task{
	topLevelTasks, nextTask, err:= createSubTasks(0, s)

	if err != nil{
		panic(err)
	}

	if nextTask != nil{
		panic("create Task error, there is -1 or smaller task level exist")
	}

	return topLevelTasks
}

func ReadTasks(filename string) []*Task{
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
