package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

// The person Type (more like an object)
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

type Conn struct {
	Status bool `json:"status,omitempty"`
}
type Level struct {
	Level int `json:"level,omitempty"`
}

var people []Person

// Display all from the people var
func GetConn(w http.ResponseWriter, r *http.Request) {
	conn := Conn{rand.Intn(2) == 1}
	json.NewEncoder(w).Encode(conn)
}

func GetLevel(w http.ResponseWriter, r *http.Request) {
	level := Level{rand.Intn(100)}
	json.NewEncoder(w).Encode(level)
}

// main function to boot up everything
func main() {
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	router.HandleFunc("/con", GetConn).Methods("GET")
	router.HandleFunc("/level", GetLevel).Methods("GET")
	fmt.Println("Server ready!!")
	log.Fatal(http.ListenAndServe(":8000", router))
}
