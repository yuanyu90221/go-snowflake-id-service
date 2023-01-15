# go-snowflake-id-service

This is a snowflake-id service implemented by golang

## [snowflake id](https://zh.wikipedia.org/zh-tw/%E9%9B%AA%E8%8A%B1%E7%AE%97%E6%B3%95)

A distributed monotone incremented Identifier generate algorithm

## format
| timestamp       | shardId        |  seqId     |               
|-----------------|----------------|------------|
| timestamp difference in millisecond(41bit) | shardId(10bit)| seqId(12bit)|

## functionality

1. in 1 millisecond could generate 2^12 different id 

2. if more than 2^12 request in 1 millisecond, this service will block this request


## for atomic operation

use golang mutex lock for generate seqId operation

```go
type Shard struct {
	mulock   sync.Mutex
	epoch    int64
	lastTime int64
	shardId  int64
	seqId    int64
}
```

first get current timestamp

check current == lastTime then mean in same millisecond, then increase seqId

if current == lastTime and seqId > seqMax then mean need to wait until next millisecond

current > lastTime mean recount the seqId to 0

```go
func (s *Shard) NextID() int64 {
	s.mulock.Lock()
	defer s.mulock.Unlock()
	now := currentMillisecond()
	if s.lastTime == now {
		s.seqId++
		// seqId over the max , then loop until next millisecond
		if s.seqId > seqIdMax {
			for now <= s.lastTime {
				now = currentMillisecond()
			}
		}
	} else {
		s.seqId = 0
		s.lastTime = now
	}
	id := int64((now-s.epoch)<<int64(timeShift) | (s.shardId << int64(shardShift)) | s.seqId)
	return id
}
```

## test with --race

```shell
go run test -v ./... --race
```

## benchmark test

```shell
go test -benchmem -run=none -bench ^BenchmarkTest$ github.com/yuanyu90221/go-snowflake-id-service/test --race
```