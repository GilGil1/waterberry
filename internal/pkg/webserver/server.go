package webserver

import (
	"fmt"
	"log"
	"net/http"
)

func Init() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":9090", nil))

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}
