package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

const (
	seqIdLength   uint8 = 12
	shardIdLength uint8 = 10
	seqIdMax      int64 = (1 << seqIdLength) - 1
	shardIdMax    int64 = (1 << shardIdLength) - 1
	shardShift    uint8 = seqIdLength
	timeShift     uint8 = shardIdLength + seqIdLength
)

type Shard struct {
	mulock   sync.Mutex
	epoch    int64
	lastTime int64
	shardId  int64
	seqId    int64
}

func NewShard(shardId int64, epoch int64) (*Shard, error) {
	if shardId < 0 || shardId > shardIdMax {
		return nil, errors.New("shardId is over maximum : " + strconv.FormatInt(shardIdMax, 10))
	}
	return &Shard{
		lastTime: 0,
		shardId:  shardId,
		epoch:    epoch,
	}, nil
}

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

func currentMillisecond() int64 {
	return time.Now().UTC().UnixMilli()
}
