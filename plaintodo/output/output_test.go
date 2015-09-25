package output

import (
	"bytes"
	"strings"
	"testing"

	"../ls"
	"../util"
)

func TestPrint(t *testing.T) {
	tasks := util.ReadTestTaskRelativePath("..")
	showTasks := ls.Ls(tasks, nil)

	// show first task
	outputTasks := make([]*ls.ShowTask, 1)
	outputTasks[0] = showTasks[0]

	buf := &bytes.Buffer{}
	var omitStrings []string
	Output(buf, outputTasks, true, 0, omitStrings)

	correctString := `go to SSA :id 1 :start 2015-02-01
  create a set list :id 2 :important :start 2015-01-31
    add music to player :id 3 :start 2015-01-30
  buy items :id 4
    buy battery :id 5
    buy ultra orange :id 6
    buy king blade :id 7

`
	results := strings.Split(buf.String(), "\n")

	corrects := strings.Split(correctString, "\n")

	if len(corrects) != len(results) {
		t.Errorf("return %d strings bud %d", len(corrects), len(results))
		t.FailNow()
	}

	for index, str := range corrects {
		if results[index] != str {
			t.Errorf("return shuld be '%s', but '%s'", str, results[index])
			t.FailNow()
		}
	}

}

func TestAllTask(t *testing.T) {
	tasks := util.ReadTestTaskRelativePath("..")
	showTasks := ls.Ls(tasks, nil)

	// show all task
	buf := &bytes.Buffer{}
	var omitStrings []string
	Output(buf, showTasks, true, 0, omitStrings)

	correctString := `go to SSA :id 1 :start 2015-02-01
  create a set list :id 2 :important :start 2015-01-31
    add music to player :id 3 :start 2015-01-30
  buy items :id 4
    buy battery :id 5
    buy ultra orange :id 6
    buy king blade :id 7

rss :id 8
  my site :id 9 :important :repeat every 1 day :start 2015-02-01 :url http://ota42y.com

`
	results := strings.Split(buf.String(), "\n")

	corrects := strings.Split(correctString, "\n")

	if len(corrects) != len(results) {
		t.Errorf("return %d strings bud %d", len(corrects), len(results))
		t.FailNow()
	}

	for index, str := range corrects {
		if results[index] != str {
			t.Errorf("return shuld be '%s', but '%s'", str, results[index])
			t.FailNow()
		}
	}
}

func TestAllTaskWithLevel(t *testing.T) {
	tasks := util.ReadTestTaskRelativePath("..")
	showTasks := ls.Ls(tasks, nil)

	// show all task
	buf := &bytes.Buffer{}
	var omitStrings []string
	Output(buf, showTasks, true, 1, omitStrings)

	correctString := `  go to SSA :id 1 :start 2015-02-01
    create a set list :id 2 :important :start 2015-01-31
      add music to player :id 3 :start 2015-01-30
    buy items :id 4
      buy battery :id 5
      buy ultra orange :id 6
      buy king blade :id 7

  rss :id 8
    my site :id 9 :important :repeat every 1 day :start 2015-02-01 :url http://ota42y.com

`
	results := strings.Split(buf.String(), "\n")

	corrects := strings.Split(correctString, "\n")

	if len(corrects) != len(results) {
		t.Errorf("return %d strings bud %d", len(corrects), len(results))
		t.FailNow()
	}

	for index, str := range corrects {
		if results[index] != str {
			t.Errorf("return shuld be '%s', but '%s'", str, results[index])
			t.FailNow()
		}
	}
}

func TestOmit(t *testing.T) {
	tasks := util.ReadTestTaskRelativePath("..")
	showTasks := ls.Ls(tasks, nil)

	correctString := "  my site :id 9 :important :repeat ... :start ... :url http://ota42y.com\n\n"

	// show first task
	outputTasks := make([]*ls.ShowTask, 1)
	outputTasks[0] = showTasks[1].SubTasks[0]

	var omitStrings []string
	omitStrings = append(omitStrings, "start")
	omitStrings = append(omitStrings, "repeat")

	buf := &bytes.Buffer{}
	Output(buf, outputTasks, true, 1, omitStrings)

	if buf.String() != correctString {
		t.Errorf("return '%s' strings but '%s'", correctString, buf.String())
		t.FailNow()
	}
}
