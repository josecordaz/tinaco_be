// rpio "github.com/stianeikeland/go-rpio"
package main

import (
	"fmt"
	"os"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

var (
	bombPin = rpio.Pin(17)
)

func main() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bombPin.Output()
	bombPin.Mode(rpio.Output)

	bombPin.High()
	// bombPin.Toggle()

	bombPin.High()
	fmt.Println("High")

	time.Sleep(time.Second * 5)

	bombPin.Low()
	bombPin.Input()
	fmt.Println("Low")

	// for i := 0; bombPin.Read() == rpio.High; i++ {
	// 	fmt.Println("i", i)
	// 	bombPin.PullDown()
	// 	bombPin.Low()
	// }
	// fmt.Println("Low")

	rpio.Close()
}
