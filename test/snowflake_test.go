package test

import (
	"testing"
	"time"

	"github.com/yuanyu90221/go-snowflake-id-service/snowflake"
)

func Test_NextID(t *testing.T) {
	// create checkMap for check duplicated
	checkMap := make(map[int64]struct{})
	epoch := time.Now().UTC().UnixMilli() - 300000
	shardTransactionId := 0
	snowflakeTx, err := snowflake.NewShard(int64(shardTransactionId), epoch)
	if err != nil {
		t.Error(err)
	}
	// buffer equals 10 means time concurrent users
	idChannel := make(chan int64, 10)
	for count := 1; count <= 10000; count++ {
		go func(ch chan int64) {
			// get snowflakeTx and push to ch
			ch <- snowflakeTx.NextID()
		}(idChannel)
	}
	for count := 1; count <= 10000; count++ {
		// receive from idChannel
		id := <-idChannel
		if _, isExist := checkMap[id]; isExist {
			t.Errorf("already exist id : %v", id)
		} else {
			checkMap[id] = struct{}{}
		}
	}
}

func BenchmarkTest(b *testing.B) {
	// create checkMap for check duplicated
	checkMap := make(map[int64]struct{})
	epoch := time.Now().UTC().UnixMilli() - 300000
	shardTransactionId := 0
	snowflakeTx, err := snowflake.NewShard(int64(shardTransactionId), epoch)
	if err != nil {
		b.Error(err)
	}
	// buffer equals 10 means time concurrent users
	idChannel := make(chan int64, 10)
	for count := 1; count <= b.N; count++ {
		go func(ch chan int64) {
			// get snowflakeTx and push to ch
			ch <- snowflakeTx.NextID()
		}(idChannel)
	}
	for count := 1; count <= b.N; count++ {
		// receive from idChannel
		id := <-idChannel
		if _, isExist := checkMap[id]; isExist {
			b.Errorf("already exist id : %v", id)
		} else {
			checkMap[id] = struct{}{}
		}
	}
}
