package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hazelcast/hazelcast-go-client"
)

func main() {
	config := hazelcast.NewConfig()
	config.Cluster.Name = "lab2"
	ctx := context.Background()

	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Failed to start Hazelcast client: %v", err)
	}
	defer client.Shutdown(ctx)

	hzMap, err := client.GetMap(ctx, "my-map")
	if err != nil {
		log.Fatalf("Failed to get map: %v", err)
	}

	for i := 0; i < 1000; i++ {
		key := i
		value := "value" + strconv.Itoa(i)
		if _, err := hzMap.Put(ctx, key, value); err != nil {
			log.Printf("An error occurred: %v", err)
		}
	}

	fmt.Println("Finished putting values into the map.")
}
