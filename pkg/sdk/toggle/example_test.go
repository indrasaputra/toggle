package toggle_test

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/pkg/sdk/toggle"
)

func ExampleNewClient() {
	ctx := context.Background()
	client, err := toggle.NewClient(ctx, "localhost:8080", grpc.WithInsecure())
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
}
