package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/sifes/docker-mtrpz/httptools"
	"github.com/sifes/docker-mtrpz/signal"
)

var port = flag.Int("port", 8080, "server port")

const confResponseDelaySec = "CONF_RESPONSE_DELAY_SEC"
const confHealthFailure = "CONF_HEALTH_FAILURE"

type Report map[string]interface{}

func (r Report) Process(req *http.Request) {
}

func (r Report) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(r)
}

func main() {
	flag.Parse()
	
	h := new(http.ServeMux)
	
	h.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Header().Set("Allow", "GET")
			return
		}
		
		rw.Header().Set("content-type", "text/plain")
		if failConfig := os.Getenv(confHealthFailure); failConfig == "true" {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte("FAILURE"))
		} else {
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte("OK"))
		}
	})
	
	report := make(Report)
	
	h.HandleFunc("/api/v1/some-data", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Header().Set("Allow", "GET")
			return
		}
		
		respDelayString := os.Getenv(confResponseDelaySec)
		if delaySec, parseErr := strconv.Atoi(respDelayString); parseErr == nil && delaySec > 0 && delaySec < 300 {
			time.Sleep(time.Duration(delaySec) * time.Second)
		}
		
		report.Process(r)
		rw.Header().Set("content-type", "application/json")
		rw.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(rw).Encode([]string{
			"1", "2",
		})
	})
	
	h.Handle("/report", report)
	
	server := httptools.CreateServer(*port, h)
	server.Start()
	signal.WaitForTerminationSignal()
}