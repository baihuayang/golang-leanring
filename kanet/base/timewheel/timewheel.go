package timewheel

import (
	// "container/list"

	"time"
)

//时间精度
const timePrecision int64 = 10 * 1000 * 1000

// type

// type callback func(...interface{})

type task struct {
	cycle int
	once  bool
	// id    interface{}
	// params   []interface{}
	interval time.Duration
	cb       func()
}

type TimeWheel struct {
	interval   time.Duration
	slots      []*taskList
	slotNum    int
	curSlotIdx int
	// addTaskChan chan task
	// removeTaskChan chan interface{}
	// stopChan       chan bool
	// cb             callback
	// timerMap map[interface{}]int

	// Ticker *time.Ticker
}

type taskList struct {
	tasks []task
	head  int
	tail  int
	full  bool
}

func (tl *taskList) push(t task) {
	if tl.head != tl.tail || !tl.full {
		tl.tasks[tl.tail] = t
	} else {
		l := len(tl.tasks)
		newList := make([]task, 2*l)
		index := tl.tail
		for i := len(newList) - l; i < len(newList); i++ {
			newList[i] = tl.tasks[index]
			index++
			if index >= l {
				index = 0
			}
		}
		tl.tasks = newList
		tl.tail = 0
		tl.head = len(newList) - l
		tl.tasks[tl.tail] = t
	}
	tl.tail++
	if tl.tail >= len(tl.tasks) {
		tl.tail = 0
	}
	if tl.tail == tl.head {
		tl.full = true
	}
}

func (tl *taskList) pop() (task, bool) {
	if tl.head != tl.tail || tl.full {
		t := tl.tasks[tl.head]
		var nt task
		tl.tasks[tl.head] = nt
		tl.head++
		if tl.head >= len(tl.tasks) {
			tl.head = 0
		}
		tl.full = false
		return t, true
	}
	return task{}, false
}

// type ITimeWheel interface {
// 	Start()
// 	loop()
// 	Stop()
// 	AddTimer(time.Duration, interface{}, ...interface{})
// 	AddOnceTimer(time.Duration, interface{}, ...interface{})
// 	RemoveTimer(interface{})
// 	tickHandler()
// 	// SetCallback(callback)
// }

func Second(sec int) int {
	return int(int64(sec*1000000000) / timePrecision)
}

func MilliSecond(ms int) int {
	return int(int64(ms*1000000) / timePrecision)
}

//New 新建时间轮实例
func New(interval time.Duration, slotNum int) *TimeWheel {
	if interval <= 0 || slotNum <= 0 {
		return nil
	}

	tw := new(TimeWheel)
	tw.interval = interval
	tw.slots = make([]*taskList, slotNum, slotNum)
	tw.curSlotIdx = 0
	// tw.addTaskChan = make(chan task, 100)
	// tw.removeTaskChan = make(chan interface{})
	// tw.stopChan = make(chan bool)
	tw.slotNum = slotNum
	// tw.Ticker = nil
	// tw.timerMap = make(map[interface{}]int)

	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = &taskList{
			tasks: make([]task, 10, 20),
			head:  0,
			tail:  0,
			full:  false,
		}
	}

	return tw
}

func (tw *TimeWheel) Start() {
	// println("TimeWheel start", tw.interval)
	// tw.Ticker = time.NewTicker(tw.interval)
	// go tw.loop()
}

// func (tw *TimeWheel) SetCallback(cb callback) {
// 	// tw.cb = cb
// }

func (tw *TimeWheel) loop() {
	for {
		select {
		// case <-tw.Ticker.C:
		// println(time.Now().UnixNano() / 1000000)
		// tw.TickHandler()
		// case t := <-tw.addTaskChan:
		// 	tw.addTask(t)
		// case id := <-tw.removeTaskChan:
		// 	tw.removeTask(id)
		// case <-tw.stopChan:
		// 	tw.Ticker.Stop()
		// 	return
		}
	}
}

func (tw *TimeWheel) TickHandler() {
	tw.curSlotIdx++
	if tw.curSlotIdx == tw.slotNum {
		tw.curSlotIdx = 0
	}
	l := tw.slots[tw.curSlotIdx]
	for t, ok := l.pop(); ok; {
		t.cycle--
		// println("tick cricle", t.cycle)
		if t.cycle <= 0 {
			t.cb()
			// tw.cb(t.params...)
			if !t.once {
				tw.addTask(t)
			}
			// if !t.once {
			// 	tw.AddTimer(t.interval, t.id, t.once, t.params...)
			// }
			// next := e.Next()
			// l.Remove(e)
			// e = next
		}
		t, ok = l.pop()
	}
	// for e := l.Pop(); e != nil; {
	// 	t := e.(task)
	// 	t.cycle--
	// 	// println("tick cricle", t.cycle)
	// 	if t.cycle <= 0 {
	// 		// tw.cb(t.params...)
	// 		if !t.once {
	// 			tw.addTask(t)
	// 		}
	// 		// if !t.once {
	// 		// 	tw.AddTimer(t.interval, t.id, t.once, t.params...)
	// 		// }
	// 		// next := e.Next()
	// 		// l.Remove(e)
	// 		// e = next
	// 	}
	// 	e = l.Pop()
	// }
}

func (tw *TimeWheel) add(interval time.Duration, once bool, cb func()) {
	t := task{
		interval: interval,
		once:     once,
		cb:       cb,
	}

	tw.addTask(t)
}

func (tw *TimeWheel) AddTimer(interval time.Duration, cb func()) {
	tw.add(interval, false, cb)
}

func (tw *TimeWheel) AddOnceTimer(interval time.Duration, cb func()) {
	tw.add(interval, true, cb)
}

func (tw *TimeWheel) addTask(t task) {
	// tms := int(t.interval.Seconds() * 1000)
	// twms := (int)(tw.interval.Seconds() * 1000)
	idx := 0
	if t.interval <= 0 {
		t.cycle = 0
		idx = tw.curSlotIdx + 1
	} else {
		t.cycle = int(t.interval / (tw.interval * time.Duration(tw.slotNum)))
		idx = tw.curSlotIdx + int((t.interval%(tw.interval*time.Duration(tw.slotNum)))/tw.interval)
	}
	// idx--
	// if idx < 0 {
	// 	idx = tw.slotNum - 1
	// }
	idx = idx % tw.slotNum
	// println("cycle : ", t.cycle, "idx : ", idx)
	tw.slots[idx].push(t)
}

func (tw *TimeWheel) RemoveTimer(id interface{}) {

}

func (tw *TimeWheel) removeTask(id interface{}) {

}
