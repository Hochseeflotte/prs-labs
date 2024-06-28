package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func defaultResponse(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Not implemented yet"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/message", defaultResponse).Methods(http.MethodGet)

	fmt.Println("Server running on port 3005")
	log.Fatal(http.ListenAndServe(":3005", r))
}
