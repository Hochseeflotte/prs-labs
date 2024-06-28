package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hazelcast/hazelcast-go-client"
)

var (
	nodeIP []string
	nodeID string
	port   int
	client hazelcast.Client
	msgMap *hazelcast.Map
	ctx    = context.Background()
)

func init() {
	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1])
	} else {
		port = 3001
	}

	switch port {
	case 3001:
		nodeIP = []string{"192.168.131.102:5701"}
		nodeID = "7473be79d30b68b31907764fde0d2cd8168852e788fabc58c197113b268aa756"
	case 3002:
		nodeIP = []string{"192.168.131.102:5702"}
		nodeID = "50ee4384d4eae646d450113b4346e44037769653aa0ea266b5e6e255b22411d3"
	case 3003:
		nodeIP = []string{"192.168.131.102:5703"}
		nodeID = "cdff39a5a895634e848324df1838f0df14d1acfadc73b336b8d66e3709363dd9"
	}
}

func stopNode(node string) {
	cmd := exec.Command("docker", "stop", node)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("Error stopping node: %v", err)
	}
	fmt.Printf("Stopped node: %s\n", out.String())
}

func logMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		uuid, exists := data["uuid"]
		if !exists {
			http.Error(w, "UUID not provided in the request", http.StatusBadRequest)
			return
		}

		fmt.Printf("id: %s\n", uuid)
		fmt.Printf("msg: %s\n", data["msg"])
		if _, err := msgMap.Put(ctx, uuid, data["msg"]); err != nil {
			http.Error(w, "Failed to log message", http.StatusInternalServerError)
			return
		}

		response := map[string]string{"status": "Message logged"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	} else if r.Method == http.MethodGet {
		keys, err := msgMap.GetKeySet(ctx)
		if err != nil {
			http.Error(w, "Failed to retrieve keys", http.StatusInternalServerError)
			return
		}

		var values []string
		for _, key := range keys {
			value, err := msgMap.Get(ctx, key)
			if err != nil {
				http.Error(w, "Failed to retrieve values", http.StatusInternalServerError)
				return
			}
			values = append(values, fmt.Sprintf("%v", value))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(values)
	}
}

func main() {
	config := hazelcast.NewConfig()
	config.Cluster.Name = "lab2"
	config.Cluster.Network.Addresses = nodeIP
	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Failed to start Hazelcast client: %v", err)
	}
	defer client.Shutdown(ctx)

	msgMap, err = client.GetMap(ctx, "messages-map")
	if err != nil {
		log.Fatalf("Failed to get map: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/log", logMessage).Methods(http.MethodPost, http.MethodGet)
	http.Handle("/", r)

	go func() {
		fmt.Printf("Server running on port %d\n", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	stopNode(nodeID)
	client.Shutdown(ctx)
	fmt.Println("Server stopped gracefully")
}
