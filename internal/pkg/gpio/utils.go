package gpio

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	LOG_FILE = "water.txt"
)

func WriteToFile(filename string, text string) error {
	var f *os.File
	var err error
	if fileExists(filename) {
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		f, err = os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer f.Close()
	_, err = f.WriteString(text)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func printModeChange(br BaseRelay) {

	WriteToFile(LOG_FILE, relayToTxt(br))

}

func relayToTxt(br BaseRelay) string {
	out := fmt.Sprintf("%s,%d,%d,%s,%s\n",
		time.Now().Format(time.RFC3339),
		br.GetId(),
		br.GetPin(),
		br.GetName(),
		br.GetCurrentMode())
	return out
}

func LoadLogs(filename string) (string, error) {
	var f *os.File
	var err error
	var result string
	var gp BaseRelay
	if fileExists(filename) {
		f, err = os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(f)
		var relayList []BaseRelay = make([]BaseRelay, 0)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			result += scanner.Text() + "\n"
			gp, _ = txtToRelay(scanner.Text())
			relayList = append(relayList, gp)
		}
		bytes, _ := json.Marshal(relayList)
		result = string(bytes)
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	} else {
		return result, fmt.Errorf("file %s not found", filename)
	}
	defer f.Close()
	return result, nil
}
func txtToRelay(line string) (BaseRelay, error) {
	items := strings.Split(line, ",")

	updateTime, err := time.Parse(time.RFC3339, items[0])
	if err != nil {
		return BaseRelay{}, fmt.Errorf("date %v malformed", err)
	}
	id, err := strconv.Atoi(items[1])
	if err != nil {
		return BaseRelay{}, fmt.Errorf("id %v malformed", err)
	}
	pin, err := strconv.Atoi(items[2])
	if err != nil {
		return BaseRelay{}, fmt.Errorf("pin %v malformed", err)
	}

	relay := BaseRelay{}
	relay.ID = id
	relay.Pin = uint8(pin)
	relay.Name = items[3]
	relay.CurrentStatus = items[4]
	relay.UpdateTime = updateTime
	return relay, nil

}
