# PlainTodo

PlainTodo is plain text based todo list system inspired by [Todo.txt](http://todotxt.com/)  
Now newest version is v0.1.4

## Feature

- plain text task list
- sub task support
  - space num means subtask
- other feature support as option
  - base format contains task name and subtask

## Task Syntax

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
  related url

### :subTaskFile
  sub task file.
  If this option set in a task, read task from specified file and set as a sub tasks 
   
```
// ./tasks/home_tasks.txt have 3 tasks

home :subTaskFile ./tasks/home_tasks.txt
  task 1 in home_tasks.txt
  task 2 in home_tasks.txt
  task 3 in home_tasks.txt
```

### :lock
  Lock task.
  If this option set in a task, it's not change any command.
  So we can't complete the task, change attribute, and other.
  If you delete this option, you must edit task file by text editor.

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

## commands

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
`

If you set task id, the new task will be subtask.
`
> ls
ls hit
go to SSA :id 1 :start 2015-02-01
  create a set list :id 2 :start 2015-01-31 :important
    add music to player :id 3 :start 2015-01-30

rss :id 8
  my site :id 9 :start 2015-02-01 :important :repeat every 1 day :url http://ota42y.com
  
> subtask change volume :id 2
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
If not set :no-sub-tasks, show all sub tasks.
If not set :complete, show not completed task.

|option|example|description|
|:id|ls :id 1| show specific task|
|:no-sub-tasks|ls :id 1 :no-sub-tasks| show specific task|
|:level| ls :level 1| show only tasks which same or lower level|
|:complete| ls :complete | show completed task|
|:overdue| ls :overdue 2015-01-31| show overdue task|

If no options, overdue task, which check :start

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
Save tasks.
This command disable by default.

This is very dangerous, use update command.
Because it save not completed tasks only.
We well change to save all task.

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

### archive
Archive task.
If a task which completed, it archive other file.
This command disable by default.

This command not change task list.
So, when this command execute, execute save command by after.
But, it^s not good, update command will archive and save, so use this.

```
> complete 5
> archive
archive hit
append tasks to archives/2015-02-01.txt
> lsall
lsall hit
go to SSA :id 1 :start 2015-02-01
  create a set list :id 2 :start 2015-01-31 :important
    add music to player :id 3 :start 2015-01-30
  buy items :id 4
    buy battery :id 5 :complete 2015-02-01 00:00
    buy ultra orange :id 6
    buy king blade :id 7

rss :id 8
  my site :id 9 :start 2015-01-31 :important :repeat every 1 day :url http://ota42y.com
```

### update
Archive, save and reload tasks.
```
> complete 5
> update
update hit
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

#### set
Set attribute to task.
If already set, this command overwrite it.

```
> ls :id 8
ls hit
rss :id 8

> set :id 8 :repeat every 1 day
set hit
set attribute rss :id 8 :repeat every 1 day
> ls :id 8
ls hit
rss :id 8 :repeat every 1 day
```

### start
Set start attribute with now datetime to task.

```
> ls :id 8
ls hit
rss :id 8

> start :id 8
start hit
set attribute rss :id 8 :start 2015-02-01 14:00
> ls :id 8
ls hit
rss :id 8 :start 2015-02-01 14:00
```

### postpone
Set postpone task.
If task postpone, ls :overdue check postpone attribute

```
> ls
ls hit
go to SSA :id 1 :start 2015-02-01
  create a set list :id 2 :important :start 2015-01-31
    add music to player :id 3 :start 2015-01-30

rss :id 8
  my site :id 9 :important :repeat every 1 day :start 2015-02-01 :url http://ota42y.com

> postpone :id 9 :postpone 1 year
postpone hit
set attribute   my site :id 9 :important :postpone 2016-02-01 00:00 :repeat every 1 day :start 2015-02-01 :url http://ota42y.com
> ls
ls hit
go to SSA :id 1 :start 2015-02-01
  create a set list :id 2 :important :start 2015-01-31
    add music to player :id 3 :start 2015-01-30

> postpone :id 9 :postpone today 20:00
postpone hit
set attribute   my site :id 9 :important :postpone 2015-02-01 20:00 :repeat every 1 day :start 2015-02-01 :url http://ota42y.com

```

### move
Move task to other task's sub task.
If move to  :id=0 task's sub task, task move to top level task

```
> lsall
lsall hit
go to SSA :id 1 :start 2015-02-01
  buy items :id 4
    buy battery :id 5
    buy ultra orange :id 6
    buy king blade :id 7

rss :id 8
  my site :id 9 :important :repeat every 1 day :start 2015-02-01 :url http://ota42y.com

> move :from 4 :to 8
move hit
task moved to sub task
parent: rss :id 8
> lsall
lsall hit
go to SSA :id 1 :start 2015-02-01

rss :id 8
  my site :id 9 :important :repeat every 1 day :start 2015-02-01 :url http://ota42y.com
  buy items :id 4
    buy battery :id 5
    buy ultra orange :id 6
    buy king blade :id 7

> move :from 4 :to 0
move hit
task moved to top level task
> lsall
lsall hit
go to SSA :id 1 :start 2015-02-01

rss :id 8
  my site :id 9 :important :repeat every 1 day :start 2015-02-01 :url http://ota42y.com

buy items :id 4
  buy battery :id 5
  buy ultra orange :id 6
  buy king blade :id 7
```

#### open
Execute open command by :url attribute

```
> open :id 9
open hit
open: http://ota42y.com
# Probably your computer open this url by web browser
```

#### nice
Do nice.
If you don't select task, this command do all tasks.

This command execute these

- convert https://www.evernote.com/shard/... to evernote:///view/... in :url attribute
- convert Today to YYYY-MM-DD in :start attribute (sorry, not implemented)

#### alias
Show command aliases
See command alias example.

```
> alias
alias hit
lsalltest = ls :level 1 :no-sub-tasks
pt = postpone :postpone 1 day :id
```

#### exit
Exit this application.
This command don't save task.


## command alias
You can create other name with option for command.

For example, you can create `all` command, which execute `ls ;level 1`. 
Alias data write to config.toml like this.

```
[command]

[[command.aliases]]
name = "lsalltest"
command = "ls :level 1 :no-sub-tasks"

[[command.aliases]]
name = "pt"
command = "postpone :postpone 1 day :id"
```

When this case, if you type `lsalltest`, it will replace `ls :level 1 :no-sub-tasks`
Alias setting replace command word, so if you set option, it's not change.
So, if you type `pt 42`, it will replace `postpone :postpone 1 day :id 42`
