// +build darwin netbsd freebsd openbsd dragonfly

package netpoll

import (
	"knet/internal"
	"log"

	"golang.org/x/sys/unix"
)

const (
	EVFilterWrite = unix.EVFILT_WRITE
	EVFilterRead  = unix.EVFILT_READ
	EVFilterSock  = -0xd
)

type eventList struct {
	size   int
	events []unix.Kevent_t
}

func newEventList(size int) *eventList {
	return &eventList{size, make([]unix.Kevent_t, size)}
}

func (el *eventList) increase() {
	el.size <<= 1
	el.events = make([]unix.Kevent_t, el.size)
}

type Poller struct {
	fd       int
	asyncJos internal.AsyncJobQueue
}

func OpenPoller() (*Poller, error) {
	poller := new(Poller)
	kfd, err := unix.Kqueue()
	if err != nil {
		return nil, err
	}
	poller.fd = kfd
	_, err = unix.Kevent(poller.fd, []unix.Kevent_t{{
		Ident:  0,
		Filter: unix.EVFILT_USER,
		Flags:  unix.EV_ADD | unix.EV_CLEAR,
	}}, nil, nil)
	if err != nil {
		return nil, err
	}
	poller.asyncJos = internal.NewAsyncJobQueue()
	return poller, err
}

func (p *Poller) Close() error {
	return unix.Close(p.fd)
}

func (p *Poller) Trigger(job internal.Job) error {
	p.asyncJos.Push(job)
	_, err := unix.Kevent(p.fd, []unix.Kevent_t{{
		Ident:  0,
		Filter: unix.EVFILT_USER,
		Fflags: unix.NOTE_TRIGGER,
	}}, nil, nil)
	return err
}

func (p *Poller) Polling(callback func(fd int, filter int16) error) error {
	// events := make([]unix.Kevent_t, 128)
	el := newEventList(1042)
	var wakeup bool
	for {
		n, err := unix.Kevent(p.fd, nil, el.events, nil)
		if err != nil && err != unix.EINTR {
			log.Println(err)
			continue
		}
		var evFilter int16
		for i := 0; i < n; i++ {
			if fd := int(el.events[i].Ident); fd != 0 {
				evFilter = el.events[i].Filter
				if (el.events[i].Flags&unix.EV_EOF != 0) || (el.events[i].Flags&unix.EV_ERROR != 0) {
					evFilter = EVFilterSock
				}
				if err = callback(fd, evFilter); err != nil {
					return err
				}
			} else {
				wakeup = true
			}
		}
		if wakeup {
			wakeup = false
			if err := p.asyncJos.Walk(); err != nil {
				return err
			}
		}
		if n == el.size {
			el.increase()
		}
	}
}

// AddRead ...
func (p *Poller) AddRead(fd int) error {
	if _, err := unix.Kevent(p.fd, []unix.Kevent_t{
		{Ident: uint64(fd), Flags: unix.EV_ADD, Filter: unix.EVFILT_READ}}, nil, nil); err != nil {
		return err
	}
	return nil
}

// AddWrite ...
func (p *Poller) AddWrite(fd int) error {
	if _, err := unix.Kevent(p.fd, []unix.Kevent_t{
		{Ident: uint64(fd), Flags: unix.EV_ADD, Filter: unix.EVFILT_WRITE}}, nil, nil); err != nil {
		return err
	}
	return nil
}

// AddReadWrite ...
func (p *Poller) AddReadWrite(fd int) error {
	if _, err := unix.Kevent(p.fd, []unix.Kevent_t{
		{Ident: uint64(fd), Flags: unix.EV_ADD, Filter: unix.EVFILT_READ},
		{Ident: uint64(fd), Flags: unix.EV_ADD, Filter: unix.EVFILT_WRITE},
	}, nil, nil); err != nil {
		return err
	}
	return nil
}

// ModRead ...
func (p *Poller) ModRead(fd int) error {
	if _, err := unix.Kevent(p.fd, []unix.Kevent_t{
		{Ident: uint64(fd), Flags: unix.EV_DELETE, Filter: unix.EVFILT_WRITE}}, nil, nil); err != nil {
		return err
	}
	return nil
}

// ModReadWrite ...
func (p *Poller) ModReadWrite(fd int) error {
	if _, err := unix.Kevent(p.fd, []unix.Kevent_t{
		{Ident: uint64(fd), Flags: unix.EV_ADD, Filter: unix.EVFILT_WRITE}}, nil, nil); err != nil {
		return err
	}
	return nil
}

// Delete ...
func (p *Poller) Delete(fd int) error {
	return unix.Close(fd)
}
