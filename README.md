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

### :due
  task deadline

### :completed
  task completed time

### :repeat
  (not support yet)
  when task completed, next deadline set

### :url
  (not support yet)
  releated url


## Usage

### sample text
task.txt
```
go to SSA :due 2015-02-01
  create a set list :due 2015-01-31
    add music to player
  buy items
    buy battery
    buy ultra orange
    buy king blade
rss
  my site :url http://ota42y.com :due 2015-02-01 :repeat every 1 day
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

#### ls
Show all overdate task, which check :due
(sorry, this isn't support query yet)

```
> ls
ls hit
go to SSA :id 1 :due 2015-02-01
  create a set list :id 2 :due 2015-01-31 :important
    add music to player :id 3 :due 2015-01-30

rss :id 8
  my site :id 9 :due 2015-02-01 :important :repeat every 1 day :url http://ota42y.com
```

#### lsall
Show all tasks
```
> lsall
go to SSA :due 2015-02-01
  create a set list :important :due 2015-01-31
    add music to player :due 2015-01-30
  buy items
    buy battery
    buy ultra orange
    buy king blade

rss
  my site :due 2015-02-01 :important :repeat every 1 day :url http://ota42y.com
```

#### complete
Complete selected task.
This command add :complete attribute to selected task.
```
> complete 5
complete hit
> lsall
lsall hit
go to SSA :id 1 :due 2015-02-01
  create a set list :id 2 :due 2015-01-31 :important
    add music to player :id 3 :due 2015-01-30
  buy items :id 4
    buy battery :id 5 :complete 2015-02-01 13:48
    buy ultra orange :id 6
    buy king blade :id 7

rss :id 8
  my site :id 9 :due 2015-02-01 :important :repeat every 1 day :url http://ota42y.com
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
go to SSA :id 1 :due 2015-02-01
  create a set list :id 2 :due 2015-01-31 :important
    add music to player :id 3 :due 2015-01-30
  buy items :id 4
    buy ultra orange :id 5
    buy king blade :id 6

rss :id 7
  my site :id 8 :due 2015-01-31 :important :repeat every 1 day :url http://ota42y.com
```

#### exit
Exit this application.
This command don't save task.
