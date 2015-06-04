package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
        "path/filepath"
        "errors"
	"regexp"
	"strconv"
	"time"
)

type Datapoint struct {
	At    time.Time `json:at`
	Value float32   `json:value`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Specify channel ID as an argument.")
	}
	channelId := os.Args[1]

	temperature, err := readTemperature()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Temperature: %f\n", temperature)

	// TODO: Pass channel ID as an argument.
	err = sendTemperature(channelId, temperature)
	if err != nil {
		log.Fatal(err)
	}
}

func readTemperature() (float32, error) {
	thermDevice, err := getThermDevice("/sys/bus/w1/devices")
	if err != nil {
		return 0, err
	}
	dat, err := ioutil.ReadFile(thermDevice)
	if err != nil {
		return 0, err
	}
	s := string(dat)

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

func getThermDevice(baseDir string) (devPath string, err error) {
        subDirInfos, err := ioutil.ReadDir(baseDir)
        if err != nil {
                return "", err
        }

        for _, fileInfo := range subDirInfos {
                var subDir = filepath.Join(baseDir, (fileInfo).Name())
                dirInfo, err := os.Stat(subDir)
                if err == nil && dirInfo.IsDir() {
                        var devfile = filepath.Join(subDir, "w1_slave")
                        _, err := os.Stat(devfile)
                        if !os.IsNotExist(err) {
                                return devfile, nil
                        }
                }
        }
        return "", errors.New("w1_slave is not found.")
}

