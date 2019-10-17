package main

import (
	"kanet"
	"kanet/actor"
	"kanet/lib/snowflake"
	network "knet"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func init() {
	kanet.Init()
	actor.Init()
	snowflake.Init(1)
	// game.InitProtoMap()
	network.Server(runtime.NumCPU() + 1)
}

func destroy() {
	network.Stop()
}

func main() {
	println("=====server start======")
	// Init()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
	destroy()
	println("========server stop==========")
}
