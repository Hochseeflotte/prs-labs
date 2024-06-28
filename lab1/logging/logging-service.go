package logging

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var messages = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func logMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		uuid, exists := data["uuid"]
		if !exists {
			http.Error(w, "UUID not provided in the request", http.StatusBadRequest)
			return
		}

		messages.Lock()
		messages.m[uuid] = data["msg"]
		messages.Unlock()

		response := map[string]string{"status": "Message logged"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	} else if r.Method == http.MethodGet {
		messages.RLock()
		msgValues := make([]string, 0, len(messages.m))
		for _, msg := range messages.m {
			msgValues = append(msgValues, msg)
		}
		messages.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msgValues)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/log", logMessage).Methods(http.MethodPost, http.MethodGet)
	http.ListenAndServe(":3001", r)
}
