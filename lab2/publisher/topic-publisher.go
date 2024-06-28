package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	hzTopic, err := client.GetTopic(ctx, "my-topic")
	if err != nil {
		log.Fatalf("Failed to get topic: %v", err)
	}

	for i := 1; i <= 100; i++ {
		msg := fmt.Sprintf("Msg %d", i)
		if err := hzTopic.Publish(ctx, msg); err != nil {
			log.Printf("An error occurred: %v", err)
		} else {
			fmt.Printf("Published: %d\n", i)
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Finished publishing messages to the topic.")
}
