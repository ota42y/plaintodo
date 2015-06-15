# PlainTodo

PlainTodo is plain text based todo list system inspired by [Todo.txt](http://todotxt.com/)


## Feature

- plain text task list
- subtask support
  - space num means subtask
- other feature support as option
  - base format contains task name and subtask

## Task Syntacks

A task composed by one line.

```
task = <space><task name><attributes>
<attributes> = ( :<attribute name> <attribute value>) | <attributes>
```

## Options

### :start
  task start time

### :completed
  task completed time

### :repeat
  `:repeat [every|after] num [minutes|hour|day|week|year]`  
  when task completed, next deadline set
  
```
:start 2015-01-30 :repeat every 1 day 
// when complete in 2015-02-14, set :start 2015-01-31
:start 2015-01-30 :repeat after 1 day 
// when complete in 2015-02-14, set :start 2015-02-15
```
  

### :url
  (not support yet)
  releated url


## Usage

### sample text
task.txt
```
go to SSA :start 2015-02-01
  create a set list :start 2015-01-31
    add music to player
  buy items
    buy battery
    buy ultra orange
    buy king blade
rss
  my site :url http://ota42y.com :start 2015-02-01 :repeat every 1 day
```

config.toml
```
[paths]
task = "./task.txt"

[archive]
directory = "archives"
nameFormat = "2006-01-02"
```

### commands

### task
Create new task.

```
> ls
ls hit
go to SSA :id 1 :start 2015-02-01
  create a set list :id 2 :start 2015-01-31 :important
    add music to player :id 3 :start 2015-01-30

rss :id 8
  my site :id 9 :start 2015-02-01 :important :repeat every 1 day :url http://ota42y.com
  
> task write reply mail :start 2015-02-01
task hit
Create task: write reply mail :id 10 :start 2015-02-01
> ls
ls hit
go to SSA :id 1 :start 2015-02-01
  create a set list :id 2 :start 2015-01-31 :important
    add music to player :id 3 :start 2015-01-30

rss :id 8
  my site :id 9 :start 2015-02-01 :important :repeat every 1 day :url http://ota42y.com
  
write reply mail :id 10 :start 2015-02-01
```

### subtask
Create sub task.

```
> ls
ls hit
go to SSA :id 1 :start 2015-02-01
  create a set list :id 2 :start 2015-01-31 :important
    add music to player :id 3 :start 2015-01-30

rss :id 8
  my site :id 9 :start 2015-02-01 :important :repeat every 1 day :url http://ota42y.com
  
> subtask 2 change volume
subtask hit
Create SubTask:
Parent: create a set list :id 2 :start 2015-01-31 :important
SubTask: change volume :id 10
> ls
ls hit
go to SSA :id 1 :start 2015-02-01
  create a set list :id 2 :start 2015-01-31 :important
    add music to player :id 3 :start 2015-01-30
    change volume :id 10

rss :id 8
  my site :id 9 :start 2015-02-01 :important :repeat every 1 day :url http://ota42y.com
```

#### ls
Show all tasks.

ls command take options.
If not set :subtask, show specific task.
If not set :complete, show not completed task.

|option|example|description|
|:id|ls :id 1| show specific task|
|:subtask|ls :id 1 :subtask| show all subtasks|
|:level| ls :level 1| show only tasks which same or lower level|
|:complete| ls :complete | show completed task|

If no options, overdate task, which check :start

```
> ls
ls hit
go to SSA :id 1 :start 2015-02-01
  create a set list :id 2 :start 2015-01-31 :important
    add music to player :id 3 :start 2015-01-30

rss :id 8
  my site :id 9 :start 2015-02-01 :important :repeat every 1 day :url http://ota42y.com
```

#### lsall
Show all tasks
```
> lsall
go to SSA :start 2015-02-01
  create a set list :important :start 2015-01-31
    add music to player :start 2015-01-30
  buy items
    buy battery
    buy ultra orange
    buy king blade

rss
  my site :start 2015-02-01 :important :repeat every 1 day :url http://ota42y.com
```

#### complete
Complete selected task.
This command add :complete attribute to selected task.
```
> complete 5
complete hit
> lsall
lsall hit
go to SSA :id 1 :start 2015-02-01
  create a set list :id 2 :start 2015-01-31 :important
    add music to player :id 3 :start 2015-01-30
  buy items :id 4
    buy battery :id 5 :complete 2015-02-01 13:48
    buy ultra orange :id 6
    buy king blade :id 7

rss :id 8
  my site :id 9 :start 2015-02-01 :important :repeat every 1 day :url http://ota42y.com
```

#### save
Save tasks and reload.
If a task which completed yesterday or before, it archive other file.

```
> complete 5
> save
save hit
append tasks to archives/2015-02-01.txt
> lsall
lsall hit
go to SSA :id 1 :start 2015-02-01
  create a set list :id 2 :start 2015-01-31 :important
    add music to player :id 3 :start 2015-01-30
  buy items :id 4
    buy ultra orange :id 5
    buy king blade :id 6

rss :id 7
  my site :id 8 :start 2015-01-31 :important :repeat every 1 day :url http://ota42y.com
```

#### exit
Exit this application.
This command don't save task.
