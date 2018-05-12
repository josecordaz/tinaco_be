// rpio "github.com/stianeikeland/go-rpio"
package main

import (
	"fmt"
	"os"

	rpio "github.com/stianeikeland/go-rpio"
)

var (
	bombPin = rpio.Pin(27)
)

func main() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bombPin.Output()

	// bombPin.High()

	// time.Sleep(time.Second * 5)

	bombPin.Low()

	rpio.Close()
}