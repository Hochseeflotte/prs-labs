package main

import (
	"context"
	"fmt"
	"log"

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

	queue, err := client.GetQueue(ctx, "my-queue")
	if err != nil {
		log.Fatalf("Failed to get queue: %v", err)
	}

	for i := 0; i < 100; i++ {
		if err := queue.Put(ctx, i); err != nil {
			log.Printf("An error occurred: %v", err)
		} else {
			fmt.Printf("Producing %d\n", i)
		}
	}

	fmt.Println("Finished producing values to the queue.")
}
