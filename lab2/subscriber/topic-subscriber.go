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
	messageListener := func(event *hazelcast.MessagePublished) {
		currentTime := time.Now().Format("2006-01-02 15:04:05.000")
		fmt.Printf("Received: %v at %s\n", event.Value, currentTime)
	}

	subscriptionID, err := hzTopic.AddMessageListener(ctx, messageListener)
	if err != nil {
		log.Fatalf("Failed to add message listener: %v", err)
	}

	fmt.Println("Listening. Press Enter to stop...")
	fmt.Scanln()

	if err := hzTopic.RemoveListener(ctx, subscriptionID); err != nil {
		log.Printf("Failed to remove message listener: %v", err)
	}
}
