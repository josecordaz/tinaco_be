package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
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

var level Level

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

	trigPin.Low()

	time.Sleep(time.Millisecond * 200)

	trigPin.High()

	time.Sleep(time.Microsecond * 20)

	trigPin.Low()

	var wg sync.WaitGroup
	var start time.Time
	var d time.Duration

	wg.Add(1)
	go func() {
		for echoPin.Read() == rpio.Low {
		}
		start = time.Now()
		for echoPin.Read() == rpio.High {
		}
		d = time.Since(start)
		wg.Done()
	}()

	wg.Wait()

	res := d.Seconds() * 17000

	return res
}

func getDistance() float64 {
	ds := make([]float64, 0)
	var avg, sum float64
	for i := 0; i < 10; i++ {
		ds = append(ds, getMeasurement())
		time.Sleep(time.Millisecond * 100)
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
	low := 20.0
	hight := 110.0

	if d < low {
		return 100
	} else if d > hight {
		return 0
	}

	rdistance := d - low
	total := hight - low

	return int(100 - ((rdistance * 100) / total))
}

func GetLevel(w http.ResponseWriter, r *http.Request) {
	d := getDistance()
	level := convertToPCT(d)
	// level.mu.Lock()
	json.NewEncoder(w).Encode(Level{level})
	// level.mu.Unlock()
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
	json.NewEncoder(w).Encode(Conn{false})
}

// main function to boot up everything
func main() {
	router := mux.NewRouter()
	// router.HandleFunc("/con", GetConn).Methods("GET")
	router.HandleFunc("/level", GetLevel).Methods("GET")
	router.HandleFunc("/b_on", TurnBombOn).Methods("GET")
	router.HandleFunc("/b_off", TurnBombOff).Methods("GET")
	router.HandleFunc("/b_status", GetBombStatus).Methods("GET")
	router.Headers("Access-Control-Allow-Origin", "*")

	// Set pin to output mode

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	trigPin.Output()

	echoPin.Input()
	bombPin.Output()

	// Open and map memory to access gpio, check for errors

	// var wg sync.WaitGroup

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	for {

	// 		d := getDistance()
	// 		level.mu.Lock()
	// 		level.Level = convertToPCT(d)
	// 		level.mu.Unlock()
	// 		fmt.Println(time.Now().Format(time.RFC3339), " level Readed ", level)
	// 		time.Sleep(time.Second)
	// 	}
	// }()

	fmt.Println("Server ready!!")
	err := http.ListenAndServe("0.0.0.0:80", handlers.CORS()(router))
	if err != nil {
		fmt.Println("Error", err)
	}

	// wg.Wait()

	rpio.Close()

}
