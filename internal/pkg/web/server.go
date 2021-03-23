package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	box "waterberry/internal/generator"
	"waterberry/internal/pkg/config"
	"waterberry/internal/pkg/gpio"
	"waterberry/internal/pkg/weather"
)

var relaysInfo []gpio.IRelay
var BuildTime string
var GitCommit string
var globalConfig config.GlobalConfig

// Init ... init the server with the relay config
func Init(relays []gpio.IRelay, config config.GlobalConfig) {
	relaysInfo = relays
	globalConfig = config
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/relays", relaysHandler)
	http.HandleFunc("/water", setHandler)
	http.HandleFunc("/log", logHandler)
	http.HandleFunc("/status", systemStatusHandler)
	http.HandleFunc("/config", configHandler)
	http.HandleFunc("/weather", weatherHandler)
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
	for _, relay := range relaysInfo {
		singleMap := relay.GetPropertiesMap()
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
		delayMinute := 1
		currentRelay.SetOn(delayMinute)

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

func configHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(globalConfig, "", "\t")
	if err == nil {
		w.Write(bytes)
	} else {
		w.Write([]byte(fmt.Sprintf("Error in logs : %v", err)))
	}
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/forecast?q=Pardesiyya,il&APPID=1047fac797d01a4424c8b01b1642e907&units=metric")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error in calling weather api : %v", err)))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	// Convert to weatherinfo
	weatherInfo, err := weather.NewWeatherInfo(body)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(weatherInfo.PrettyPrint()))
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

// func RunTimmgs() {

// 	for _, relay := range relaysInfo {
// 		br := relay.GetBaseRelay()
// 		br.GetCurrentMode()
// 	}
// }
