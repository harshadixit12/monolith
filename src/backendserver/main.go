package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port int

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Pong from server: %d", port)))
}

func doWork(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Response from server: %d", port)))
}

func main() {
	portPtr := flag.Int("port", 8080, "Port the HTTP server binds to.")

	flag.Parse()
	port = *portPtr

	fmt.Println(port, portPtr)

	http.HandleFunc("GET /ping", healthCheck)
	http.HandleFunc("GET /work", doWork)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		log.Fatal("Unable to start server")
	}
}
