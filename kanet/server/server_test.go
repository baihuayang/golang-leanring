package main

import (
	"kanet/actor"
	"kanet/base/rpc"
	"kanet/service"
	"reflect"
	"testing"
)

func TestServer(t *testing.T) {
	ids := make([]int64, 10000)
	ch := make(chan bool, 10000)
	for i := 0; i < len(ids); i++ {
		ts := new(service.TestService)
		ts.Init(10)
		ids[i] = ts.ActorID()
	}
	tts := new(service.TestService)
	tts.Init(10000)
	tts.SetCh(ch)

	m := make(map[int]string)
	m[1000] = "1000"
	m[2000] = "2000"
	m[3000] = "1000"
	m[4000] = "2000"
	m[5000] = "1000"
	m[6000] = "2000"
	m[7000] = "1000"
	m[8000] = "2000"
	m[9000] = "1000"
	m[10000] = "2000"
	m[11000] = "1000"
	m[12000] = "2000"
	sl := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17}
	// println(actor.MGR.GetActor(ids[0]).ActorID())

	for i := 0; i < len(ids); i++ {
		actor.MGR.GetActor(ids[i]).SendByID(tts.ActorID(), "test_cb", 42, "kanonlee", m, sl)
	}
	// actor.MGR.GetActor(ids[0]).SendByID(ids[1], "test_cb", 42, "kanonlee", m, sl)

	for i := 0; i < len(ch); i++ {
		<-ch
	}
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	// <-c
}

func TestServer2(t *testing.T) {
	ids := make([]int64, 10000)
	ch := make(chan bool, 10000)
	for i := 0; i < len(ids); i++ {
		ts := new(service.TestService)
		ts.Init(10)
		ids[i] = ts.ActorID()
	}
	tts := new(service.TestService)
	tts.Init(10000)
	tts.SetCh(ch)

	m := make(map[int]string)
	m[1000] = "1000"
	m[2000] = "2000"
	m[3000] = "1000"
	m[4000] = "2000"
	m[5000] = "1000"
	m[6000] = "2000"
	m[7000] = "1000"
	m[8000] = "2000"
	m[9000] = "1000"
	m[10000] = "2000"
	m[11000] = "1000"
	m[12000] = "2000"
	sl := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17}
	// println(actor.MGR.GetActor(ids[0]).ActorID())

	for i := 0; i < len(ids); i++ {
		var ci actor.CallIO
		ci.Buff = rpc.Pack("test_cb", 42, "kanonlee", m, sl)
		ci.Source = ids[i]
		tts.Send(ci)
		// actor.MGR.GetActor(ids[i]).SendByID(ids[100], "test_cb", 42, "kanonlee", m, sl)
	}
	// actor.MGR.GetActor(ids[0]).SendByID(ids[1], "test_cb", 42, "kanonlee", m, sl)

	for i := 0; i < len(ch); i++ {
		<-ch
	}
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	// <-c
}

func f1(s string, i int, m map[int]string, sl []int) {
	s = s + "lee"
	i++
	m[42] = "GZ"
	sl[0] = 99
}

func f2(args ...interface{}) {
	s := args[0].(string)
	i := args[1].(int)
	m := args[2].(map[int]string)
	sl := args[3].([]int)

	s = s + "lee"
	i++
	m[42] = "GZ"
	sl[0] = 99
}

func TestReflect1(t *testing.T) {
	s := "kanon"
	i := 4444
	m := make(map[int]string)
	m[1000] = "1000"
	m[2000] = "2000"
	sl := []int{0, 1, 2, 3, 4, 5, 6, 7}
	args := []interface{}{s, i, m, sl}

	for i := 0; i < 1000000; i++ {
		fv := reflect.ValueOf(f1)
		ft := reflect.TypeOf(f1)
		in := make([]reflect.Value, len(args))
		for i, arg := range args {
			in[i] = reflect.ValueOf(arg)
			if ft.In(i).Kind() != in[i].Kind() {

			}
		}
		fv.Call(in)
	}
}

func TestReflect2(t *testing.T) {
	s := "kanon"
	i := 4444
	m := make(map[int]string)
	m[1000] = "1000"
	m[2000] = "2000"
	sl := []int{0, 1, 2, 3, 4, 5, 6, 7}
	args := []interface{}{s, i, m, sl}

	for i := 0; i < 1000000; i++ {
		f2(args...)
	}
}
