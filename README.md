# Task Tracker CLI

A lightweight, terminal-based task manager written in Go. Built for the [roadmap.sh Task Tracker project](https://roadmap.sh/projects/task-tracker), it keeps your tasks in a local JSON file—no database or third-party dependencies required.

## Features

- **Create tasks** with a description and an automatically assigned ID.
- **Update descriptions** and **delete tasks** by ID.
- **Track progress** with `todo`, `in-progress`, and `done` statuses.
- **List tasks** all at once or filter by status.
- **Keep tasks between runs** in a human-readable JSON file.
- **Never reuse deleted IDs**: the next ID is stored separately from the task list.
- **Track timestamps** for task creation and modification.
- **Handle invalid input** with clear errors and nonzero exit codes.

## Getting started

**Requirement:** Go 1.27 or newer (as declared in `go.mod`; the project uses `encoding/json/v2`).

```bash
git clone https://github.com/Waraidako/task-tracker-cli.git
cd task-tracker-cli
go build .
```

This produces `task-cli` on macOS/Linux or `task-cli.exe` on Windows. From the project directory, run `./task-cli` as shown below. You can also use `go run .` instead of building.

## Usage

| Command | Action |
| --- | --- |
| `./task-cli add "Buy groceries"` | Create a task |
| `./task-cli update 1 "Buy groceries and milk"` | Update task 1's description |
| `./task-cli delete 1` | Delete task 1 |
| `./task-cli mark-in-progress 1` | Mark task 1 as in progress |
| `./task-cli mark-done 1` | Mark task 1 as done |
| `./task-cli list` | Show every task |
| `./task-cli list todo` | Show only to-do tasks |
| `./task-cli list in-progress` | Show only in-progress tasks |
| `./task-cli list done` | Show only completed tasks |
| `./task-cli --help` | Display command help |

### Example workflow

```bash
# Create a task
./task-cli add "Write the README"
# Added task Write the README: ID = 1

# Start working on it
./task-cli mark-in-progress 1

# Change its description
./task-cli update 1 "Finish the project README"

# Mark it completed
./task-cli mark-done 1

# Show completed tasks
./task-cli list done

# Remove it
./task-cli delete 1
```

**Notes:** Put descriptions containing spaces in quotes. Task IDs must be positive integers and must identify an existing task. The `list` command takes at most one optional filter, and only `todo`, `in-progress`, and `done` are accepted.

### Example task output

```text
Task Finish the project README:
  - id: 1
  - status: done
  - created at: 2026-10-26 12:00:00
  - updated at: 2026-10-26 12:05:00
```

## Data storage

The application saves data to `task-list.json` in the **current working directory**. If the file doesn't exist or is empty, it starts with an empty task list; data is saved when a task change succeeds.

```json
{
    "next_id": 3,
    "tasks": [
        {
            "id": 1,
            "description": "Buy groceries",
            "status": "todo",
            "created_at": "2026-10-26 12:00:00",
            "updated_at": "2026-10-26 12:00:00"
        }
    ]
}
```

The `next_id` field ensures IDs won't be reassigned when tasks are deleted. Saves use a temporary file followed by a rename to reduce the risk of incomplete JSON writes. `task-list.json` is excluded from Git tracking.

This is a simple local CLI, not a multi-user database: concurrent writes from separate processes are not locked.

## Project structure

```text
task-tracker-cli/
├── main.go                     # CLI parsing, commands, and output
├── internal/
│   ├── task/task.go            # Task model and task operations
│   └── storage/storage.go      # JSON persistence
├── go.mod
└── .gitignore
```

## Project reference

[Task Tracker — roadmap.sh](https://roadmap.sh/projects/task-tracker)
