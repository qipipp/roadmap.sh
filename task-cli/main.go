package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

const file = "tasks.json"

func load() ([]Task, error) {
	b, err := os.ReadFile(file)
	if os.IsNotExist(err) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(b)) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(b, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func save(tasks []Task) error {
	b, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, b, 0644)
}

func getTime() string {
	return time.Now().Format(time.RFC3339)
}

func get_id(tasks []Task) int {
	m := 0
	for _, t := range tasks {
		if t.ID > m {
			m = t.ID
		}
	}
	return m + 1
}

func get_task(id int, s string) Task {
	cur := getTime()
	t := Task{id, s, "todo", cur, cur}
	return t
}

func (t Task) print_task() {
	fmt.Printf("{id: %v, description: %v, status: %v, createdAt: %v, updatedAt: %v}\n", t.ID, t.Description, t.Status, t.CreatedAt, t.UpdatedAt)
}

func add(s string) {
	tasks, err := load()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	t := get_task(get_id(tasks), s)
	tasks = append(tasks, t)
	if err := save(tasks); err != nil {
		fmt.Println("save error:", err)
		return
	}
	fmt.Printf("Task added successfully (Id: %v)", t.ID)
}

func update(id int, s string) {
	tasks, err := load()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	found := false
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Description = s
			tasks[i].UpdatedAt = getTime()
			found = true
		}
	}
	if !found {
		fmt.Printf("Task not found (Id: %v)", id)
		return
	}
	if err := save(tasks); err != nil {
		fmt.Println("save error:", err)
		return
	}
	fmt.Printf("Task updated succesfully\n")
}

func del(id int) {
	tasks, err := load()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	idx := -1
	for i, t := range tasks {
		if t.ID == id {
			idx = i
		}
	}
	if idx == -1 {
		fmt.Printf("Task not found (Id: %v)", id)
		return
	}
	tasks = append(tasks[:idx], tasks[idx+1:]...)
	if err := save(tasks); err != nil {
		fmt.Println("save error:", err)
		return
	}
	fmt.Printf("Task deleted succesfully\n")
}

func mark(id int, s string) {
	tasks, err := load()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	found := false
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Status = s
			tasks[i].UpdatedAt = getTime()
			found = true
		}
	}
	if !found {
		fmt.Printf("Task not found (Id: %v)", id)
		return
	}
	if err := save(tasks); err != nil {
		fmt.Println("save error:", err)
		return
	}
	fmt.Printf("Task marked succesfully\n")
}

func listing(s string) {
	tasks, err := load()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, t := range tasks {
		if s == "" || s == t.Status {
			t.print_task()
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("need command")
		return
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	switch cmd {
	case "add":
		if len(args) < 1 {
			fmt.Println("add need: description")
			return
		}
		add(args[0])
	case "update":
		if len(args) < 2 {
			fmt.Println("update needs: id desc")
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("invalid id:", args[0])
			return
		}
		update(id, args[1])
	case "delete":
		if len(args) < 1 {
			fmt.Println("delete need: id")
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("invalid id:", args[0])
			return
		}
		del(id)
	case "mark-in-progress":
		if len(args) < 1 {
			fmt.Println("mark-in-progress need: id")
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("invalid id:", args[0])
			return
		}
		mark(id, "in-progress")
	case "mark-done":
		if len(args) < 1 {
			fmt.Println("mark-done need: id")
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("invalid id:", args[0])
			return
		}
		mark(id, "done")
	case "list":
		if len(args) < 1 {
			listing("")
			return
		}
		if args[0] != "done" && args[0] != "todo" && args[0] != "in-progress" {
			fmt.Println("invalid args")
		}
		listing(args[0])
	default:
		fmt.Println("invalid cmd")
	}
}
