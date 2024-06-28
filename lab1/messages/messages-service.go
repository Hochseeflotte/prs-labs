package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func defaultResponse(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Not implemented"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/message", defaultResponse).Methods(http.MethodGet)
	http.ListenAndServe(":3002", r)
}
