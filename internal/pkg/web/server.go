package web

import (
	"bytes"
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

// Init ... init the server with the relay config
func Init(relays []gpio.IRelay) {
	relaysInfo = relays
	http.HandleFunc("/", statusHandler)
	http.HandleFunc("/water", setHandler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func watchTimers() {

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

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	key := r.URL.RequestURI()
	b := box.Get(key)
	if b != nil {
		w.Write(b)
		return
	}
	_ = b
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<html><head>
	<title>WaterPi Irrigation System</title>
	<meta http-equiv="refresh" content="5">

	<script type="text/javascript">
		function setRelay(relaynum, mode ){
			const xhttp = new XMLHttpRequest();
			url = "water?relay="+relaynum+"&mode="+mode;
			//alert(url);
			xhttp.open("GET", url, true);
			xhttp.send();
			location.reload();

		}
	</script>

	</head><body>`)
	fmt.Fprintf(w, fmt.Sprintf("<p>Time is : %s</p>", string(time.Now().Format(time.RFC3339))))

	if !isWindows() {
		fmt.Fprintf(w, fmt.Sprintf("<p>Temp is : %s</p>", getCPUTemp()))

	}
	// print status table
	fmt.Fprintf(w, "<table border =1>	")
	fmt.Fprintf(w, fmt.Sprintf("<tr><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th></th><th></th><th>%s</th></tr>",
		"#", "id", "Name", "Status", "Seconds to off"))
	for key, relay := range relaysInfo {
		lineCode := fmt.Sprintf(`
		<tr><td>%d</td><td>%d</td><td>%s</td><td>%s</td>
		<td><button onclick='setRelay(%d, "on")'>On</button></td>
		<td><button onclick='setRelay(%d, "off")'>Off</button></td>
		<td>%s</td>
		</tr>`,
			key,
			relay.GetId(),
			relay.GetName(),
			relay.GetCurrentMode(),
			key,
			key,
			relay.GetSecondsOff(),
		)
		fmt.Fprintln(w, lineCode)
	}
	fmt.Fprintln(w, "</table>")

	fmt.Fprintln(w, "</body></html>")

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
	return tempOut
}

func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}
