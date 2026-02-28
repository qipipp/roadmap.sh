# GitHub User Activity CLI (Go)

A small CLI tool that fetches a GitHub user’s **recent public activity** using the GitHub REST API and prints a simple, human-readable summary in the terminal.

## Project URL
`https://roadmap.sh/projects/github-user-activity`

## Requirements
- Go (any recent version)
- Internet connection

> This project uses only Go’s standard library (no external dependencies).

## What it does
- Accepts a GitHub username as an argument
- Calls the GitHub events endpoint: `https://api.github.com/users/<username>/events`
- Prints a short summary for common event types (Push, Issues, Watch/Star, Fork, PR, Release, etc.)
- Currently fetches **up to 10 events** (see `main.go`)

## Build

### Windows (CMD)
~~~bat
go build -o github-activity.exe
~~~

### macOS / Linux
~~~bash
go build -o github-activity
~~~

## Run

### Option 1) Run without building
~~~bash
go run . <username>
~~~

Example:
~~~bash
go run . kamranahmedse
~~~

### Option 2) Run the built binary

#### Windows (CMD)
~~~bat
github-activity.exe <username>
~~~

#### PowerShell
~~~powershell
.\github-activity.exe <username>
~~~

#### macOS / Linux
~~~bash
./github-activity <username>
~~~

## Usage
~~~text
github-activity <username>
~~~

If you forget the username:
~~~text
github-activity need: username
~~~

## Example output
~~~text
1): Pushed to owner/repo
2): Starred owner/another-repo
3): Interacted with an issue in owner/repo
~~~

## Notes
- GitHub returns **public events** for the user (so private activity won’t show up).
- GitHub API has rate limits (unauthenticated requests are limited). If you hit limits, try again later or extend the tool to support authentication.
- Output formatting is intentionally simple and best-effort: unknown event types are printed as `<Type> in <repo>`.

## Possible improvements (optional)
- Add a `--limit` flag to control the number of events instead of hard-coding 10
- Show commit counts for `PushEvent` by parsing the payload
- Support caching to avoid repeated API calls
