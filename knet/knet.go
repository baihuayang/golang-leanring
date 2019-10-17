package knet

import (
	"errors"
	"log"
	"net"
)

var ErrClosing = errors.New("closing event loop")

var inst *server

func Server(loopNum int) error {
	return serve(loopNum)
}

func Listen(address string, ev eventHandler) error {
	var ln listener
	ln.network = "tcp"
	ln.handler = ev
	var err error
	ln.ln, err = net.Listen(ln.network, address)
	if err != nil {
		return err
	}
	if err = ln.system(); err != nil {
		return err
	}
	log.Println("listen:", ln.lnAddr.String())
	return inst.listen(&ln, ev)
}

func Stop() {
	inst.stop()
}

func printIfError(err error) {
	if err != nil {
		log.Println(err)
	}
}
