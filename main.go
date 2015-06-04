package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

type Datapoint struct {
	At    time.Time `json:at`
	Value float32   `json:value`
}

func main() {
	temperature, err := readTemperature()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Temperature: %f\n", temperature)

	// TODO: Pass channel ID as an argument.
	err = sendTemperature("556a961004957c0003000001", temperature)
	if err != nil {
		log.Fatal(err)
	}
}

func readTemperature() (float32, error) {
	// TODO: Get the device ID dynamically.
	out, err := exec.Command("cat", "/sys/bus/w1/devices/28-0414703e47ff/w1_slave").Output()
	if err != nil {
		return 0, err
	}
	s := bytes.NewBuffer(out).String()

	pattern := regexp.MustCompile("t=(\\d+)")
	matches := pattern.FindStringSubmatch(s)
	temperature, err := strconv.ParseFloat(matches[1], 32)
	if err != nil {
		return 0, err
	}

	return float32(temperature / 1000), nil
}

func sendTemperature(channelId string, temperature float32) error {
	url := "http://beergarden.herokuapp.com/channels/" + channelId + "/datapoints"
	log.Println(url)

	marshaled, err := json.Marshal(&Datapoint{time.Now(), temperature})
	if err != nil {
		return err
	}
	fmt.Printf("Request body: %s\n", marshaled)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshaled))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
	fmt.Println("Response headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response body:", string(body))

	return nil
}
