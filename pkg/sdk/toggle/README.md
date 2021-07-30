## Example Usage

[snip]:#
```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/indrasaputra/toggle/entity"

	"github.com/indrasaputra/toggle/pkg/sdk/toggle"
	"google.golang.org/grpc"
)

func main() {
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
	get(client, key, "after create")
	_ = client.Enable(ctx, key)
	get(client, key, "after enable")
	_ = client.Disable(ctx, key)
	get(client, key, "after disable")
	_ = client.Delete(ctx, key)
	get(client, key, "after delete")
}

func get(client *toggle.Client, key, msg string) {
	resp, _ := client.Get(context.Background(), key)
	fmt.Println(msg)
	fmt.Println(resp)
}
```