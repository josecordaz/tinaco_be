package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	rpio "github.com/stianeikeland/go-rpio"
)

var (
	// Use mcu pin 10, corresponds to physical pin 19 on the pi
	bombPin = rpio.Pin(17)
	trigPin = rpio.Pin(23)
	echoPin = rpio.Pin(24)
)

var bomba_state = false

var cont int64

type Conn struct {
	Status bool `json:"status,omitempty"`
}
type Level struct {
	Level int `json:"level,omitempty"`
}

func Swap(arrayzor []float64, i, j int) {
	tmp := arrayzor[j]
	arrayzor[j] = arrayzor[i]
	arrayzor[i] = tmp
}

func bubbleSort(arrayzor []float64) {
	swapped := true
	for swapped {
		swapped = false
		for i := 0; i < len(arrayzor)-1; i++ {
			if arrayzor[i+1] < arrayzor[i] {
				Swap(arrayzor, i, i+1)
				swapped = true
			}
		}
	}
}

func getMeasurement() float64 {
	// GPIO.setup(TRIG,GPIO.OUT)
	trigPin.Output()

	// GPIO.output(TRIG,0)
	trigPin.Low()

	// GPIO.setup(ECHO,GPIO.IN)
	echoPin.Input()

	// time.sleep(0.1)
	time.Sleep(time.Millisecond * 100)

	// GPIO.output(TRIG, 1)
	trigPin.High()

	// time.sleep(0.00001)
	time.Sleep(time.Microsecond * 10)

	// GPIO.output(TRIG, 0)
	trigPin.Low()

	// while GPIO.input(ECHO) == 0:
	//     pass
	// start = time.time()
	for echoPin.Read() == rpio.Low {
	}
	start := time.Now()

	// while GPIO.input(ECHO) == 1:
	//     pass
	// stop = time.time()
	for echoPin.Read() == rpio.High {
	}
	d := time.Since(start)

	res := d.Seconds() * 17000

	// fmt.Println("Distancia ", res)
	return res
}

func getDistance() float64 {
	ds := make([]float64, 0)
	var avg, sum float64
	for i := 0; i < 10; i++ {
		ds = append(ds, getMeasurement())
	}

	bubbleSort(ds)

	ds = ds[2:7]

	for _, v := range ds {
		sum += v
	}

	avg = sum / 5

	return avg
}

func convertToPCT(d float64) int {
	low := 5.5
	hight := 22.5

	if d < low {
		return 100
	} else if d > hight {
		return 0
	}

	rdistance := d - low
	total := hight - low

	tmp := ((rdistance * 100) / total)
	fmt.Println("tmp", tmp)

	return int(100 - ((rdistance * 100) / total))
}

// Display all from the people var
func GetConn(w http.ResponseWriter, r *http.Request) {
	conn := Conn{rand.Intn(2) == 1}
	json.NewEncoder(w).Encode(conn)
}

func GetLevel(w http.ResponseWriter, r *http.Request) {
	d := getDistance()
	level := convertToPCT(d)
	json.NewEncoder(w).Encode(Level{level})
	cont++
	fmt.Println(cont, " Level: ", level)
}

func GetBombStatus(w http.ResponseWriter, r *http.Request) {
	// conn := Conn{rand.Intn(2) == 1}
	conn := Conn{bomba_state}
	json.NewEncoder(w).Encode(conn)
}

func TurnBombOn(w http.ResponseWriter, r *http.Request) {
	bombPin.Low()
	bomba_state = true
	json.NewEncoder(w).Encode(Conn{true})
}

func TurnBombOff(w http.ResponseWriter, r *http.Request) {
	bombPin.High()
	bomba_state = false
	json.NewEncoder(w).Encode(Conn{true})
}

// main function to boot up everything
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/con", GetConn).Methods("GET")
	router.HandleFunc("/level", GetLevel).Methods("GET")
	router.HandleFunc("/b_on", TurnBombOn).Methods("GET")
	router.HandleFunc("/b_off", TurnBombOff).Methods("GET")
	router.HandleFunc("/b_status", GetBombStatus).Methods("GET")
	router.Headers("Access-Control-Allow-Origin", "*")

	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	bombPin.Output()

	fmt.Println("Server ready!!")
	err := http.ListenAndServe(":8000", handlers.CORS()(router))
	if err != nil {
		fmt.Println("Error", err)
	}

}
