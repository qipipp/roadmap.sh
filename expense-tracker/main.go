package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

type Expense struct {
	ID          int    `json:"id"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
}

func listExpenses() error {
	expenses, err := LoadJSON[[]Expense](path)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "# ID\tDate\tDescription\tAmount")

	for _, e := range expenses {
		fmt.Fprintf(w, "# %d\t%s\t%s\t$%d\n", e.ID, e.Date, e.Description, e.Amount)
	}

	return w.Flush()
}

const path = "expenses.json"

func getId(expenses []Expense) int {
	m := 0
	for _, e := range expenses {
		m = max(m, e.ID)
	}
	return m + 1
}

func getExpense(description string, amount int) (Expense, error) {
	expenses, err := LoadJSON[[]Expense](path)
	if err != nil {
		return Expense{}, err
	}
	id := getId(expenses)
	return Expense{id, GetDate(), description, amount}, nil
}

func summaryExpenses(month int) error {
	expenses, err := LoadJSON[[]Expense](path)
	if err != nil {
		return err
	}

	total := 0

	if month == 0 {
		for _, e := range expenses {
			total += e.Amount
		}
		fmt.Printf("Total expenses: $%d\n", total)
		return nil
	}

	if month < 1 || month > 12 {
		return fmt.Errorf("invalid month: %d (use 1~12)", month)
	}

	year := time.Now().Year()
	for _, e := range expenses {
		t, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			continue
		}
		if t.Year() == year && int(t.Month()) == month {
			total += e.Amount
		}
	}

	monthName := time.Month(month).String()
	fmt.Printf("Total expenses for %s: $%d\n", monthName, total)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("need command")
		return
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	switch cmd {
	case "add":
		if len(args) < 2 {
			fmt.Printf("add need: description, amount")
			return
		}
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("invalid amount:", args[1])
			return
		}
		e, err := getExpense(args[0], amount)
		if err != nil {
			fmt.Printf("failed to make expense: %s", err)
			return
		}
		err = AddJSON(path, e)
		if err != nil {
			fmt.Printf("failed to add expense: %s", err)
		}

	case "delete":
		if len(args) < 1 {
			fmt.Printf("delete need: id")
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("invalid id: %s", args[0])
			return
		}
		err = DelJSON[Expense](path, id, func(e Expense) int { return e.ID })
		if err != nil {
			fmt.Printf("failed to delete expense: %s", err)
			return
		}
	case "update":
		if len(args) < 1 {
			fmt.Println("update need: id")
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("invalid id: %s\n", args[0])
			return
		}
		fs := flag.NewFlagSet("update", flag.ContinueOnError)
		description := fs.String("description", "", "new description")
		amount := fs.Int("amount", -1, "new amount")
		if err := fs.Parse(args[1:]); err != nil {
			return
		}
		if *description == "" && *amount == -1 {
			fmt.Println("nothing to update (use --description and/or --amount)")
			return
		}
		err = UpdateJson[Expense](
			path,
			id,
			func(e Expense) int { return e.ID },
			func(e *Expense) error {
				if *description != "" {
					e.Description = *description
				}
				if *amount != -1 {
					e.Amount = *amount
				}
				return nil
			},
		)
	case "list":
		if err := listExpenses(); err != nil {
			fmt.Printf("failed to list: %s\n", err)
			return
		}
	case "summary":
		fs := flag.NewFlagSet("summary", flag.ContinueOnError)
		month := fs.Int("month", 0, "month (1-12)")
		if err := fs.Parse(args); err != nil {
			return
		}
		if err := summaryExpenses(*month); err != nil {
			fmt.Printf("failed to summary: %s\n", err)
		}
	default:
		fmt.Println("invalid cmd")
	}
}
