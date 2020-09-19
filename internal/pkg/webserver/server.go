package webserver

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"waterberry/internal/pkg/gpio"
)

func Init() {
	http.HandleFunc("/water", handler)
	log.Fatal(http.ListenAndServe(":9090", nil))

}

func handler(w http.ResponseWriter, r *http.Request) {
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
	gpio.SetRelayMode(relayNum, mode[0])
}
