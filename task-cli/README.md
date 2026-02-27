# Task Tracker CLI (Go)

A simple CLI task tracker written in Go.  
Tasks are stored in a local JSON file (`tasks.json`) in the current directory.

## Requirements
- Go (any recent version)

## Setup (first run)
This project reads tasks from `tasks.json`. Create an empty file like this on first run:

~~~json
[]
~~~

## Build
### Windows (CMD)
~~~bat
go build -o task-cli.exe
~~~

## Run
### Windows (CMD)
~~~bat
task-cli add "Buy groceries"
task-cli list
~~~

> PowerShell users:
~~~powershell
.\task-cli.exe add "Buy groceries"
.\task-cli.exe list
~~~

## Commands

### Add a task
~~~bat
task-cli add "Buy groceries"
~~~

### Update a task
~~~bat
task-cli update 1 "Buy groceries and cook dinner"
~~~

### Delete a task
~~~bat
task-cli delete 1
~~~

### Mark status
~~~bat
task-cli mark-in-progress 1
task-cli mark-in-done 1
~~~

### List tasks
~~~bat
task-cli list
task-cli list todo
task-cli list in-progress
task-cli list done
~~~

## Data format
`tasks.json` is a JSON array of tasks.

Example:
~~~json
[
  {
    "id": 1,
    "description": "Buy groceries",
    "status": "todo",
    "createdAt": "2026-02-28T12:00:00+09:00",
    "updatedAt": "2026-02-28T12:00:00+09:00"
  }
]
~~~

## Notes
- This is a CLI app (no server). Each command reads/writes `tasks.json` and exits.
- If you want to match roadmap.sh spec exactly, you may want to rename `mark-in-done` to `mark-done`.