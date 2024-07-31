package main

import (
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"time"
)

var failed int

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered in main", "trace", r)
		}
	}()

	logger.Info("Service starting...")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("The kubelet is checking my liveness/readiness")
		//fail randomly several times
		i := rand.Int()
		if i%2 == 0 && failed < 5 {
			failed++
			logger.Info("liveness probe failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info("Service is healthy")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Service is healthy"))
	})
	http.HandleFunc("/startup", func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("The startup probe is successful")
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/hello-there", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("General Kenobi"))
	})

	//imitate long startup
	timer := time.NewTimer(15 * time.Second)
	<-timer.C
	logger.Info("Service started")
	http.ListenAndServe(":8080", nil)
}
