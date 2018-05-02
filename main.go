package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

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

	trigPin.Low()

	time.Sleep(time.Millisecond * 100)

	trigPin.High()

	time.Sleep(time.Microsecond * 10)

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
		time.Sleep(time.Millisecond * 200)
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
	// d := getDistance()
	// level := convertToPCT(d)
	json.NewEncoder(w).Encode(Level{100})
	cont++
	fmt.Println(cont, " Level: ", 100)
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
	// router := mux.NewRouter()
	// router.HandleFunc("/con", GetConn).Methods("GET")
	// router.HandleFunc("/level", GetLevel).Methods("GET")
	// router.HandleFunc("/b_on", TurnBombOn).Methods("GET")
	// router.HandleFunc("/b_off", TurnBombOff).Methods("GET")
	// router.HandleFunc("/b_status", GetBombStatus).Methods("GET")
	// router.Headers("Access-Control-Allow-Origin", "*")

	// Open and map memory to access gpio, check for errors

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := rpio.Open(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Unmap gpio memory when done
			// defer rpio.Close()

			// Set pin to output mode
			bombPin.Output()

			trigPin.Output()

			echoPin.Input()
			d := getDistance()
			level := convertToPCT(d)
			rpio.Close()
			fmt.Println(time.Now(), "level Readed ", level)
			time.Sleep(time.Millisecond * 100)
		}
	}()

	wg.Wait()

	// fmt.Println("Server ready!!")
	// err := http.ListenAndServe(":8000", handlers.CORS()(router))
	// if err != nil {
	// 	fmt.Println("Error", err)
	// }

}
