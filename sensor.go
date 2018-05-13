// // import RPi.GPIO as GPIO
// // import time

// // GPIO.setmode(GPIO.BCM)

// // TRIG = 23
// // ECHO = 24

// // print "Distance Measurement In Progress"

// // GPIO.setup(TRIG,GPIO.OUT)
// // GPIO.output(TRIG,0)

// // GPIO.setup(ECHO,GPIO.IN)

// // time.sleep(0.1)

// // print "Stargin measurement"

// // GPIO.output(TRIG, 1)
// // time.sleep(0.00001)
// // GPIO.output(TRIG, 0)

// // while GPIO.input(ECHO) == 0:
// //     pass
// // start = time.time()

// // while GPIO.input(ECHO) == 1:
// //     pass
// // stop = time.time()

// // print (stop - start) * 17000

package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

var (
	// Use mcu pin 10, corresponds to physical pin 19 on the pi
	trigPin = rpio.Pin(23)
	echoPin = rpio.Pin(24)
)

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

func main() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	trigPin.Output()

	echoPin.Input()

	fmt.Println("Test", int(getDistance()))

	rpio.Close()
}

// import (
// 	"fmt"
// 	"os"
// 	"time"

// 	rpio "github.com/stianeikeland/go-rpio"
// )

// var (
// 	// Use mcu pin 10, corresponds to physical pin 19 on the pi
// 	trigPin = rpio.Pin(23)
// 	echoPin = rpio.Pin(24)
// )

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

// func getMeasurement() float64 {
// 	// GPIO.setup(TRIG,GPIO.OUT)
// 	trigPin.Output()

// 	// GPIO.output(TRIG,0)
// 	trigPin.Low()

// 	// GPIO.setup(ECHO,GPIO.IN)
// 	echoPin.Input()

// 	// time.sleep(0.1)
// 	time.Sleep(time.Millisecond * 100)

// 	// GPIO.output(TRIG, 1)
// 	trigPin.High()

// 	// time.sleep(0.00001)
// 	time.Sleep(time.Microsecond * 10)

// 	// GPIO.output(TRIG, 0)
// 	trigPin.Low()

// 	// while GPIO.input(ECHO) == 0:
// 	//     pass
// 	// start = time.time()
// 	for echoPin.Read() == rpio.Low {
// 	}
// 	start := time.Now()

// 	// while GPIO.input(ECHO) == 1:
// 	//     pass
// 	// stop = time.time()
// 	for echoPin.Read() == rpio.High {
// 	}
// 	d := time.Since(start)

// 	res := d.Seconds() * 17000

// 	// fmt.Println("Distancia ", res)
// 	return res
// }

// func getDistance() float64 {
// 	ds := make([]float64, 0)
// 	var avg, sum float64
// 	for i := 0; i < 10; i++ {
// 		ds = append(ds, getMeasurement())
// 	}

// 	bubbleSort(ds)

// 	ds = ds[2:7]

// 	for _, v := range ds {
// 		sum += v
// 	}

// 	avg = sum / 5

// 	return avg
// }

// func convertToPCT(d float64) int {
// 	low := 3.5
// 	hight := 22.5

// 	if d < low {
// 		return 100
// 	} else if d > hight {
// 		return 0
// 	}

// 	rdistance := d - low
// 	total := hight - low

// 	return int(100 - ((rdistance * 100) / total))
// }

// func main() {
// 	// Open and map memory to access gpio, check for errors
// 	if err := rpio.Open(); err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	d := getDistance()

// 	fmt.Println("Distance", d)
// 	fmt.Println("Percentage", convertToPCT(d))

// 	// Unmap gpio memory when done
// 	defer rpio.Close()

// }
