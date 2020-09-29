package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type GlobalConfig struct {
	Relays    []RelayConfig    `json:"relays"`
	Opentimes []OpenTimeConfig `json:"timings"`
}

type RelayConfig struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Pin  uint8  `json:"pin"`
}

type OpenTimeConfig struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	RelayID         int    `json:"relay_id"`
	Day             []int  `json:"days"`
	StartHour24H    int    `json:"start_hour"`
	OpenTimeMinutes int    `json:"open_time_min"`
}

func LoadConfig(filname string) GlobalConfig {
	config := GlobalConfig{}
	// Open our jsonFile
	jsonFile, err := os.Open(filname)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	fmt.Println(string(byteValue))

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	if err := json.Unmarshal(byteValue, &config); err != nil {
		panic(err)
	}
	return config
}
