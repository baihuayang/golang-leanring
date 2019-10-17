package snowflake

import (
	"log"
	"sync"
	"time"
)

var snowFlake *SnowFlake

const (
	InitalTimeStamp int64 = 1483228800000 //2017-01-01
	WokerIDBits     uint  = 10
	SequenceBits    uint  = 12
)

type SnowFlake struct {
	lastTimeStamp int64
	sequence      int64
	workerID      int64
	sync.Mutex
}

func ParseUUID(uuid int64) (ts, workerID, sequence int64) {
	shift := WokerIDBits + SequenceBits
	ts = (uuid >> (shift)) + InitalTimeStamp
	workerID = (uuid & (1<<(shift) - 1)) >> SequenceBits
	sequence = uuid & ((1 << SequenceBits) - 1)
	return
}

func (s *SnowFlake) Init(workerID int) {
	if workerID < 0 || workerID > ((1<<WokerIDBits)-1) {
		log.Fatalln("wokerID must between 0 and 1023")
		return
	}
	s.workerID = int64(workerID)
}

func (s *SnowFlake) UUID() int64 {
	s.Lock()
	defer s.Unlock()
	now := time.Now().UnixNano() / 1000000
	if s.lastTimeStamp == now {
		s.sequence = (s.sequence + 1)
		if s.sequence > ((1 << SequenceBits) - 1) {
			s.sequence = 0
			for now <= s.lastTimeStamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		s.sequence = 0
	}
	s.lastTimeStamp = now
	return (s.lastTimeStamp-InitalTimeStamp)<<(WokerIDBits+SequenceBits) |
		(s.workerID << SequenceBits) | s.sequence
}

func Init(workerID int) {
	snowFlake = new(SnowFlake)
	snowFlake.Init(workerID)
}

func UUID() int64 {
	return snowFlake.UUID()
}
