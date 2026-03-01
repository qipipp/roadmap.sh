# Expense Tracker (CLI)

A simple CLI app to track expenses. Data is stored in `expenses.json`.

## Requirements
- Go 1.18+

## Run
```bash
go run . add "Lunch" 20
go run . list
go run . summary
go run . summary --month 8
go run . delete 1
