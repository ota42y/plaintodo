package main

func ReadTestTasks() []*Task {
	filename := "test_task.txt"
	return ReadTasks(filename)
}
