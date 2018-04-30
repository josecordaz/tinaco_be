package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"os"

// 	"github.com/gorilla/handlers"
// 	"github.com/gorilla/mux"
// 	rpio "github.com/stianeikeland/go-rpio"
// )

// var (
// 	// Use mcu pin 10, corresponds to physical pin 19 on the pi
// 	bombPin = rpio.Pin(17)
// )

// var bomba_state = false

// // The person Type (more like an object)
// type Person struct {
// 	ID        string   `json:"id,omitempty"`
// 	Firstname string   `json:"firstname,omitempty"`
// 	Lastname  string   `json:"lastname,omitempty"`
// 	Address   *Address `json:"address,omitempty"`
// }

// type Address struct {
// 	City  string `json:"city,omitempty"`
// 	State string `json:"state,omitempty"`
// }

// type Conn struct {
// 	Status bool `json:"status,omitempty"`
// }
// type Level struct {
// 	Level int `json:"level,omitempty"`
// }

// // Display all from the people var
// func GetConn(w http.ResponseWriter, r *http.Request) {
// 	conn := Conn{rand.Intn(2) == 1}
// 	json.NewEncoder(w).Encode(conn)
// }

// func GetLevel(w http.ResponseWriter, r *http.Request) {
// 	level := Level{rand.Intn(100)}
// 	json.NewEncoder(w).Encode(level)
// }

// func GetBombStatus(w http.ResponseWriter, r *http.Request) {
// 	// conn := Conn{rand.Intn(2) == 1}
// 	conn := Conn{bomba_state}
// 	json.NewEncoder(w).Encode(conn)
// }

// func TurnBombOn(w http.ResponseWriter, r *http.Request) {
// 	bombPin.Low()
// 	bomba_state = true
// 	json.NewEncoder(w).Encode(Conn{true})
// }

// func TurnBombOff(w http.ResponseWriter, r *http.Request) {
// 	bombPin.High()
// 	bomba_state = false
// 	json.NewEncoder(w).Encode(Conn{true})
// }

// func checkLevelWater() {
// 	// Here we will have to read from ultrasensor
// 	// and it will update all the variables and states
// 	// and if the level is lower than limit it will turn the bomb on
// }

// // main function to boot up everything
// func main() {
// 	router := mux.NewRouter()
// 	router.HandleFunc("/con", GetConn).Methods("GET")
// 	router.HandleFunc("/level", GetLevel).Methods("GET")
// 	router.HandleFunc("/b_on", TurnBombOn).Methods("GET")
// 	router.HandleFunc("/b_off", TurnBombOff).Methods("GET")
// 	router.HandleFunc("/b_status", GetBombStatus).Methods("GET")
// 	router.Headers("Access-Control-Allow-Origin", "*")

// 	// Open and map memory to access gpio, check for errors
// 	if err := rpio.Open(); err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	// Unmap gpio memory when done
// 	defer rpio.Close()

// 	// Set pin to output mode
// 	bombPin.Output()

// 	fmt.Println("Server ready!!")
// 	go checkLevelWater()
// 	log.Fatal(http.ListenAndServe(":8000", handlers.CORS()(router)))
// }
