# PlainTodo

PlainTodo is plain text based todo list system inspired by [Todo.txt](http://todotxt.com/)


## Feature

- plain text task list
- subtask support
  - space num means subtask
- other feature support as option
  - base system contains task name and subtask

## Example

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

## Task Syntacks

A task composed by one line.

```
task = <space><task name><attributes>
<attributes> = ( :<attribute name> <attribute value>) | <attributes>
```

## Options

### :due
  task deadline

### :repeat
  when task completed, next deadline set

### :url
  releated url