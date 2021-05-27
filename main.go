package main

import (
	"GoMicroservices2/Activity/calc"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	err := http.ListenAndServe(":8080", handlers())
	if err != nil {
		log.Fatal(err)
	}
}

func handlers() http.Handler {
	h := http.NewServeMux()
	h.HandleFunc("/", incHandler)
	return h
}

func incHandler(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("action")
	value1 := r.FormValue("val1")
	value2 := r.FormValue("val2")

	if action == "" {
		http.Error(w, "missing value", http.StatusBadRequest)
		return
	}

	number1, err := strconv.Atoi(value1)
	if err != nil {
		http.Error(w, "not a number: "+value1, http.StatusBadRequest)
		return
	}

	number2, err := strconv.Atoi(value2)
	if err != nil {
		http.Error(w, "not a number: "+value2, http.StatusBadRequest)
		return
	}
	calculation := calc.Add(number1, number2)
	if action == "Add" {
		fmt.Fprintln(w, calculation)
	}
}
