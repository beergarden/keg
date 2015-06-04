package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

func main() {
	// TODO: Get the device ID dynamically.
	out, err := exec.Command("cat", "/sys/bus/w1/devices/28-0414703e47ff/w1_slave").Output()
	if err != nil {
		log.Fatal(err)
	}

	s := bytes.NewBuffer(out).String()

	pattern := regexp.MustCompile("t=(\\d+)")
	matches := pattern.FindStringSubmatch(s)
	log.Println(matches[1])
	temperature, err := strconv.ParseFloat(matches[1], 32)
	if err != nil {
		log.Fatal(err)
	}
	temperature = temperature / 1000

	fmt.Printf("Temperature: %f", temperature)
}
