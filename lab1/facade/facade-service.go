package facade

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
)

const (
	LoggingServiceURL  = "http://localhost:3001/log"
	MessagesServiceURL = "http://localhost:3002/message"
)

func handleFacade(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		message := string(body)
		uniqueID := uuid.New().String()

		logRequestBody, _ := json.Marshal(map[string]string{
			"uuid": uniqueID,
			"msg":  message,
		})

		logResp, err := http.Post(LoggingServiceURL, "application/json", strings.NewReader(string(logRequestBody)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer logResp.Body.Close()

		var logRespBody map[string]interface{}
		if err := json.NewDecoder(logResp.Body).Decode(&logRespBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"uuid":     uniqueID,
			"response": logRespBody,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	} else if r.Method == http.MethodGet {
		logResp, err := http.Get(LoggingServiceURL)
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

		response := map[string]interface{}{
			"logged_messages":  strings.Join(logRespBody, " "),
			"default_response": messagesRespBody,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/facade", handleFacade).Methods(http.MethodPost, http.MethodGet)
	log.Fatal(http.ListenAndServe(":3000", r))
}
