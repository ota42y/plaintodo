package main

func ReadTestTasks() []*Task{
	filename := "task.txt"
	return ReadTasks(filename)
}
