package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func LoadJSON[T any](path string) (T, error) {
	var zero T

	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return zero, nil
		}
		return zero, err
	}
	if len(bytes.TrimSpace(b)) == 0 {
		return zero, nil
	}

	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		return zero, err
	}
	return v, nil
}

func SaveJSON[T any](path string, v T) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func AddJSON[T any](path string, v T) error {
	arr, err := LoadJSON[[]T](path)
	if err != nil {
		return err
	}
	arr = append(arr, v)
	err = SaveJSON(path, arr)
	if err != nil {
		return err
	}
	return nil
}

func DelJSON[T any](path string, id int, getID func(T) int) error {
	arr, err := LoadJSON[[]T](path)
	if err != nil {
		return err
	}
	idx := -1
	for i, t := range arr {
		if getID(t) == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("id %d not found", id)
	}
	arr = append(arr[:idx], arr[idx+1:]...)
	return SaveJSON(path, arr)
}

func UpdateJson[T any](path string, id int, getID func(T) int, apply func(*T) error) error {
	arr, err := LoadJSON[[]T](path)
	if err != nil {
		return err
	}

	for i := range arr {
		if getID(arr[i]) == id {
			if err := apply(&arr[i]); err != nil {
				return err
			}
			return SaveJSON(path, arr)
		}
	}
	return fmt.Errorf("id %d not found", id)
}

func GetDate() string {
	return time.Now().Format("2006-01-02")
}

func GetTime() string {
	return time.Now().Format(time.RFC3339)
}
