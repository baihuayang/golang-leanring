package cyclelist

import "testing"

func TestList(t *testing.T) {
	// println("cyclelist test")
	// t.Error("=======cycly list test========")
	// println("=======cycly list test========")
	cl := New()
	cl.Push(1)
	cl.Push(2)
	cl.Push(3)
	cl.Push(4)
	cl.Push(5)

	println("elen 111", cl.eLen)
	for {
		// v := cl.Pop()
		if v := cl.Pop(); v != nil {
			// println(v.(int))
		} else {
			break
		}
	}
	println("cl.len = ", cl.len)
	cl.Push(1)
	cl.Push(2)
	cl.Push(3)
	cl.Push(4)
	cl.Push(5)
	println(cl.Pop().(int))
	cl.Push(6)

	println("cl.elen = ", cl.eLen)
}
