package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

// use map instead array.
var people map[string]Person

// this function doesn't need to be exported, so, can be lowercase.
// also, i've renamed to people
func all(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

// this function doesn't need to be exported, so, can be lowercase.
// also, I've renamed to person
func findByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	found, ok := people[params["id"]]
	if ok {
		json.NewEncoder(w).Encode(found)
	} else {
		http.Error(w, "Not found.", http.StatusNotFound)
	}
}

// this function doesn't need to be exported, so, can be lowercase.
// also, I've renamed to create
func store(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, "Error deserializing request body. I'm a teapot!", http.StatusTeapot)
		return
	}
	person.ID = params["id"]
	people[params["id"]] = person
	w.WriteHeader(http.StatusCreated)
}

// this function doesn't need to be exported, so, can be lowercase.
// also, I've renamed to delete
func remove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	delete(people, id)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Let's populate the map with initial values
	people = map[string]Person{
		"1": Person{ID: "1", Firstname: "Bono", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}},
		"2": Person{ID: "2", Firstname: "The Edge", Lastname: "Doe", Address: &Address{City: "City Y", State: "State Y"}},
		"3": Person{ID: "3", Firstname: "Adam", Lastname: "Doe"},
		"4": Person{ID: "4", Firstname: "Larry", Lastname: "Doe"},
	}

	router := mux.NewRouter()
	router.HandleFunc("/people", all).Methods("GET")
	router.HandleFunc("/people/{id}", findByID).Methods("GET")
	router.HandleFunc("/people/{id}", store).Methods("PUT")
	router.HandleFunc("/people/{id}", remove).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
