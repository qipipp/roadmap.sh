package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var lengthUnits = []string{"millimeter", "centimeter", "meter", "kilometer", "inch", "foot", "yard", "mile"}
var lengthToCm = []float64{0.1, 1, 100, 100000, 2.54, 30.48, 91.44, 160934.4}

var weightUnits = []string{"milligram", "gram", "kilogram", "ounce", "pound"}
var weightToKg = []float64{0.000001, 0.001, 1, 0.02834952, 0.45359237}

var temperatureUnits = []string{"celsius", "fahrenheit", "kelvin"}

func getLengthUnitIdx(s string) int {
	for i, t := range lengthUnits {
		if s == t {
			return i
		}
	}
	return -1
}

func convertLength(unit1 string, v1 float64, unit2 string) (float64, error) {
	idx1 := getLengthUnitIdx(unit1)
	idx2 := getLengthUnitIdx(unit2)
	if idx1 == -1 || idx2 == -1 {
		return 0, fmt.Errorf("invalid unit: %q -> %q", unit1, unit2)
	}
	return v1 * lengthToCm[idx1] / lengthToCm[idx2], nil
}

func getWeightUnitIdx(s string) int {
	for i, t := range weightUnits {
		if s == t {
			return i
		}
	}
	return -1
}

func convertWeight(unit1 string, v1 float64, unit2 string) (float64, error) {
	idx1 := getWeightUnitIdx(unit1)
	idx2 := getWeightUnitIdx(unit2)
	if idx1 == -1 || idx2 == -1 {
		return 0, fmt.Errorf("invalid unit: %q -> %q", unit1, unit2)
	}
	return v1 * weightToKg[idx1] / weightToKg[idx2], nil
}

func toCelsius(v float64, from string) (float64, error) {
	switch from {
	case "celsius":
		return v, nil
	case "fahrenheit":
		return (v - 32) * 5.0 / 9.0, nil
	case "kelvin":
		return v - 273.15, nil
	default:
		return 0, fmt.Errorf("invalid temperature unit: %q", from)
	}
}

func fromCelsius(c float64, to string) (float64, error) {
	switch to {
	case "celsius":
		return c, nil
	case "fahrenheit":
		return c*9.0/5.0 + 32, nil
	case "kelvin":
		return c + 273.15, nil
	default:
		return 0, fmt.Errorf("invalid temperature unit: %q", to)
	}
}

func convertTemp(v float64, from, to string) (float64, error) {
	c, err := toCelsius(v, from)
	if err != nil {
		return 0, err
	}
	return fromCelsius(c, to)
}

type ConvertRequest struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
	From  string  `json:"from"`
	To    string  `json:"to"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ConvertResponse struct {
	Result float64 `json:"result"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	var req ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
		return
	}

	var (
		result float64
		err    error
	)

	switch req.Type {
	case "length":
		result, err = convertLength(req.From, req.Value, req.To)
	case "weight":
		result, err = convertWeight(req.From, req.Value, req.To)
	case "temperature":
		result, err = convertTemp(req.Value, req.From, req.To)
	default:
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid type"})
		return
	}

	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, ConvertResponse{Result: result})
}

func main() {
	http.HandleFunc("/convert", convertHandler)
	log.Println("listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
