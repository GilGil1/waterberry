package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	box "waterberry/internal/generator"
	"waterberry/internal/pkg/gpio"
)

var relaysInfo []gpio.IRelay
var BuildTime string
var GitCommit string

// Init ... init the server with the relay config
func Init(relays []gpio.IRelay) {
	relaysInfo = relays
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/relays", relaysHandler)
	http.HandleFunc("/water", setHandler)
	http.HandleFunc("/log", logHandler)
	http.HandleFunc("/status", systemStatusHandler)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}

}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
func relaysHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	jsonData := make([]map[string]interface{}, 0)
	for _, v := range relaysInfo {
		singleMap := v.GetPropertiesMap()
		jsonData = append(jsonData, singleMap)
	}
	resp, err := json.MarshalIndent(jsonData, "", "\t")
	if err != nil {
		fmt.Printf("%v", err)
	}
	w.Write(resp)
}

type SystemStatus struct {
	SystemTime string
	CpuTemp    string
	GitCommit  string
	BuildTime  string
}

func systemStatusHandler(w http.ResponseWriter, r *http.Request) {
	status := SystemStatus{
		SystemTime: string(time.Now().Format(time.RFC3339)),
		GitCommit:  GitCommit,
		BuildTime:  BuildTime,
	}
	if !isWindows() {
		status.CpuTemp = getCPUTemp()
	}
	bytes, _ := json.MarshalIndent(status, "", "\t")
	w.Write(bytes)
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	res, err := gpio.LoadLogs("water.txt")
	if err == nil {
		w.Write([]byte(res))
	} else {
		w.Write([]byte(fmt.Sprintf("Error in logs : %v", err)))
	}

}
func setHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)

	relay, ok := r.URL.Query()["relay"]
	if !ok || len(relay[0]) < 1 {
		log.Println("Url Param 'relay' is missing")
		return
	}
	relayNum, err := strconv.Atoi(relay[0])
	if err != nil {
		log.Printf("relay %s \n cannot  be converted to num", relay)
		return
	}

	mode, ok := r.URL.Query()["mode"]
	if !ok || len(mode[0]) < 1 {
		log.Println("Url Param 'mode' is missing")
		return
	}

	// Sety Relay mode:
	currentRelay := relaysInfo[relayNum]
	switch mode[0] {

	case gpio.RelayOn:
		currentRelay.SetOn()

	case gpio.RelayOff:
		currentRelay.SetOff()

	default:
		fmt.Printf("Mode unrecognized : %s\n", mode)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	key := r.URL.RequestURI()
	if key == "/" || key == "" {
		key = "/index.html"
	}
	b := box.Get(key)
	if b != nil {
		w.Write(b)
		return
	}
}

func runCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return out.String(), nil
}

func getCPUTemp() string {
	var err error

	var tempOut string
	tempOut, err = runCommand("vcgencmd", []string{"measure_temp"})
	if err != nil {
		tempOut = "Could not read temperature"
	}
	return strings.Replace(tempOut, "temp=", "", 1)
}

func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}

func RunTimmgs() {

	for _, relay := range relaysInfo {
		relay.GetCurrentMode()
	}
}

// func RelayExpectedModeNow(gp) {

// }
