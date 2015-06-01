package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
)

var (
	// mcu pin 10 corresponds to physical pin 19.
	pin = rpio.Pin(10)
)

func main() {
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

	pin.Output()

	for x := 0; x < 20; x++ {
		fmt.Print(".")
		pin.Toggle()
		time.Sleep(time.Second / 5)
	}

	fmt.Println("")
	fmt.Println("Done")
}
