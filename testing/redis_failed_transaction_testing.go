package testing

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

func CreateJsonTx(client *redis.Client) {
	ping := client.Ping(context.Background())
	if ping.Err() != nil {
		panic(ping.Err())
	}
	ctx := context.Background()
	jsondata, _ := json.Marshal(map[string]any{
		"status":     "completed",
		"dataFormat": "application/json",
	})
	client.Set(ctx, "1", string(jsondata), time.Minute*10)
}

func CreateXmlTx(client *redis.Client) {
	ping := client.Ping(context.Background())
	if ping.Err() != nil {
		panic(ping.Err())
	}
	ctx := context.Background()
	jsondata, _ := json.Marshal(map[string]any{
		"status":     "completed",
		"dataFormat": "application/xml",
	})
	client.Set(ctx, "2", string(jsondata), time.Minute*10)
}
