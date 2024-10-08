package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

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
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Service is healthy"))
	})
	http.HandleFunc("/startup", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("The startup probe is successful")
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/hello-there", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain")
		val := os.Getenv("response")
		w.Write([]byte(val))
	})

	//imitate long startup
	timer := time.NewTimer(15 * time.Second)
	<-timer.C
	logger.Info("Service started")
	http.ListenAndServe(":8080", nil)
}
