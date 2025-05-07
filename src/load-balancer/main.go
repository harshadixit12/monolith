package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type BackendServer struct {
	host        url.URL
	numRequests int
}

func (b *BackendServer) GetProxyServer() *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(&b.host)
}

var ConnectedBackends map[string]BackendServer = make(map[string]BackendServer)

// Simple routing logic - pick backend server based on lowest requests served
func getNextBackend() (*BackendServer, error) {
	if len(ConnectedBackends) == 0 {
		return nil, fmt.Errorf("No backend servers connected :(")
	}

	minReqs := math.MaxInt
	selectedHost := ""

	for host, server := range ConnectedBackends {
		if server.numRequests < minReqs {
			minReqs = server.numRequests
			selectedHost = host
		}
	}

	b, ok := ConnectedBackends[selectedHost]
	if ok {
		b.numRequests++
	}

	return &b, nil
}

func handleRequests(w http.ResponseWriter, r *http.Request) {
	server, err := getNextBackend()
	if err != nil {
		w.WriteHeader(503)
		return
	}
	w.Header().Add("X-Forwarded-Server", server.host.String())
	server.GetProxyServer().ServeHTTP(w, r)
}

func registerBackend(w http.ResponseWriter, r *http.Request) {
	host := r.Header.Get("backend")

	parsedUrl, err := url.Parse(host)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Host is not a valid URL"))
	}
	ConnectedBackends[host] = BackendServer{host: *parsedUrl, numRequests: 0}
}

func main() {
	http.HandleFunc("POST /register", registerBackend)
	http.HandleFunc("/", handleRequests)
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatal("Server did not start :(")
	}
}
