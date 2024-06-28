package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	MessagesServiceURL = "http://localhost:3005/message"
)

var LoggingServiceURLs = []string{
	"http://localhost:3001/log",
	"http://localhost:3002/log",
	"http://localhost:3003/log",
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		loggingServiceURL := LoggingServiceURLs[rand.Intn(len(LoggingServiceURLs))]
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		message := string(body)
		uniqueID := uuid.New().String()
		payload := map[string]string{"uuid": uniqueID, "msg": message}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := http.Post(loggingServiceURL, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var respBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"uuid":     uniqueID,
			"response": respBody,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else if r.Method == http.MethodGet {
		loggingServiceURL := LoggingServiceURLs[rand.Intn(len(LoggingServiceURLs))]

		messagesResp, err := http.Get(MessagesServiceURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer messagesResp.Body.Close()

		var messagesRespBody []string
		if err := json.NewDecoder(messagesResp.Body).Decode(&messagesRespBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logResp, err := http.Get(loggingServiceURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer logResp.Body.Close()

		var logRespBody []string
		if err := json.NewDecoder(logResp.Body).Decode(&logRespBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"logged_messages":  strings.Join(logRespBody, " "),
			"default_response": messagesRespBody,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	r := mux.NewRouter()
	r.HandleFunc("/facade", handleMessage).Methods(http.MethodPost, http.MethodGet)
	log.Fatal(http.ListenAndServe(":3000", r))
}
