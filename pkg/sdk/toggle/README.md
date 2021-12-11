## Example Usage

[snip]:#
```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"github.com/sony/gobreaker"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/internal/messaging"
	"github.com/indrasaputra/toggle/pkg/sdk/toggle"
)

func main() {
	// start of non-circuit breaker client
	ctx := context.Background()
	config := &toggle.DialConfig{
		Host:    "localhost:8080",
		Options: []grpc.DialOption{grpc.WithInsecure()},
	}
	client, err := toggle.NewClient(config, nil)
	if err != nil {
		log.Printf("err init client: %v\n", err)
		return
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "toggle",
	})
	subscriber := messaging.NewKafkaSubscriber(reader)
	go client.Subscribe(ctx, subscriber, []string{"toggle-test-1", "toggle-test-2", "toggle-test-3"})

	key := "toggle-test-1"
	req := &entity.Toggle{
		Key:         key,
		Description: "first try using toggle",
	}

	_ = client.Create(ctx, req)
	get(client, key, "after create")
	_ = client.Enable(ctx, key)
	get(client, key, "after enable")
	_ = client.Disable(ctx, key)
	get(client, key, "after disable")
	_ = client.Delete(ctx, key)
	get(client, key, "after delete")

	// end of non-circuit breaker client
	//
	// start of circuit-breaker client

	setting := gobreaker.Settings{
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 100 && failureRatio >= 0.6
		},
		Timeout: 2 * time.Second,
	}
	breaker := gobreaker.NewCircuitBreaker(setting)

	client, err = toggle.NewClient(config, breaker)
	if err != nil {
		log.Printf("err init client: %v\n", err)
		return
	}
	// end of circuit-breaker client
}

func get(client *toggle.Client, key, msg string) {
	resp, _ := client.Get(context.Background(), key)
	fmt.Println(msg)
	fmt.Println(resp)
}
```
