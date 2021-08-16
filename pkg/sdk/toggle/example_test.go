package toggle_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/pkg/sdk/toggle"
)

func ExampleNewClient() {
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

	key := "toggle-test"
	req := &entity.Toggle{
		Key:         key,
		Description: "first try using toggle",
	}

	_ = client.Create(ctx, req)

	resp, err := client.Get(context.Background(), key)
	fmt.Println(err)
	fmt.Println(resp)

	_ = client.Enable(ctx, key)
	_ = client.Disable(ctx, key)
	_ = client.Delete(ctx, key)

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
