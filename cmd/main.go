package main

import (
	"log"
	"time"

	"github.com/yuanyu90221/go-snowflake-id-service/snowflake"
)

func main() {
	currentDate := time.Now().UTC().UnixMilli() - 30000
	txShardId := 0
	orderShardId := 1
	snowFlakeTx, _ := snowflake.NewShard(int64(txShardId), int64(currentDate))
	snowFlakeOrder, _ := snowflake.NewShard(int64(orderShardId), int64(currentDate))
	for i := 0; i < 100000; i++ {
		txId := snowFlakeTx.NextID()
		orderId := snowFlakeOrder.NextID()
		log.Printf("txId: %v\t", txId)
		log.Printf("orderId: %v\n", orderId)
	}
}
