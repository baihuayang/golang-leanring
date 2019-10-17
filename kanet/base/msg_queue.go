package base

import "sync"

type (
	MsgData struct {
		Local  bool
		Source int64
		Type   uint8
		Buff   []byte
		fname  string
		args   []interface{}
	}

	MsgQueue struct {
		q    []MsgData
		head int
		tail int
		full bool
		l    sync.Locker
	}
)

func (mq *MsgQueue) expand() {
	l := len(mq.q)
	if l > 1024 {
		l += 32
	} else {
		l *= 2
	}
	newList := make([]MsgData, l)
	index := mq.tail
	for i := len(newList) - l; i < len(newList); i++ {
		newList[i] = mq.q[index]
		index++
		if index >= l {
			index = 0
		}
	}
	mq.q = newList
	mq.tail = 0
	mq.head = len(newList) - l
}

func (mq *MsgQueue) Push(msg MsgData) {
	mq.l.Lock()
	if mq.head != mq.tail || !mq.full {
		mq.q[mq.tail] = msg
	} else {
		mq.expand()
		mq.q[mq.tail] = msg
	}
	mq.tail++
	if mq.tail >= len(mq.q) {
		mq.tail = 0
	}
	if mq.tail == mq.head {
		mq.full = true
	}
	mq.l.Unlock()
}

func (mq *MsgQueue) Pop() (msg MsgData, ok bool) {
	mq.l.Lock()
	if mq.head != mq.tail || mq.full {
		msg = mq.q[mq.head]
		var m MsgData
		mq.q[mq.head] = m
		mq.head++
		if mq.head >= len(mq.q) {
			mq.head = 0
		}
		mq.full = false
		mq.l.Unlock()
		ok = true
		return
	}
	mq.l.Unlock()
	msg = MsgData{}
	ok = false
	return
}
