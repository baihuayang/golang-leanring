package cyclelist

// CycleList represents a doubly ring linked list.
// The zero value for List is an empty list ready to use.
type CycleList struct {
	head, tail *Element
	len        int // current list length excluding (this) sentinel element
	eLen       int
	// full       bool
}

// Element is an element of a linked list.
type Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element

	// The list to which this element belongs.
	// list *CycleList

	// The value stored with this element.
	Value interface{}
}

//New a CycleList
func New() *CycleList {
	return new(CycleList).Init()
}

//Len return the len of CycleList
func (l *CycleList) Len() int {
	return l.len
}

//Init init the new CycleList
func (l *CycleList) Init() *CycleList {
	e := new(Element)
	l.head = e
	l.tail = e
	e.next = e
	e.prev = e
	l.len = 0
	l.eLen = 1
	// l.full = false
	return l
}

func (l *CycleList) insert(e, at *Element) *Element {
	e.next = at.next
	at.next = e
	e.prev = at
	e.next.prev = e
	// l.len++
	l.eLen++
	return e
}

func (l *CycleList) insertValue(at *Element) *Element {
	e := new(Element)
	e.next = at.next
	at.next = e
	e.prev = at
	e.next.prev = e
	// l.len++
	l.eLen++
	return e
}

func (l *CycleList) insertBefore(at *Element) *Element {
	e := new(Element)
	e.prev = at.prev
	at.prev = e
	e.next = at
	e.prev.next = e
	// l.len++
	l.eLen++
	return e
}

//Push a new value
func (l *CycleList) Push(v interface{}) {
	if l.tail == l.head {
		if l.len == 0 {
			l.tail.Value = v
			l.tail = l.tail.next
		} else {
			e := l.insertBefore(l.head)
			e.Value = v
			// l.tail = e
		}
	} else {
		l.tail.Value = v
		l.tail = l.tail.next
	}
	l.len++
}

//Pop a value
func (l *CycleList) Pop() interface{} {
	if l.len == 0 {
		return nil
	}
	v := l.head.Value
	l.head.Value = nil
	l.head = l.head.next
	l.len--
	// l.full = false
	return v
}
