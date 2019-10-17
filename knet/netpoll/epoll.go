// +build linux

package netpoll

import (
	"kanet/network/knet/internal"

	"golang.org/x/sys/unix"
)

type Poller struct {
	fd       int
	wfd      int
	asyncJos internal.AsyncJobQueue
	wfdbuf   []byte
}

const (
	ErrEvents = unix.EPOLLERR | unix.EPOLLHUP
	OutEvents = unix.EPOLLPOUT
	InEvents  = ErrEvents | unix.EPOLLRDHUP | unix.EPOLLIN
)

type eventList struct {
	size   int
	events []unix.EpollEvent
}

func newEventList(size int) *eventList {
	return &eventList{size, make([]unix.EpollEvent, size)}
}

func (el *eventList) increase() {
	el.size <<= 1
	el.events = make([]unix.EpollEvent, el.size)
}

func OpenPoller() (*Poller, error) {
	poller := new(Poller)
	epollFD, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	poller.fd = epollFD
	r0, _, errno := unix.Syscall(unix.SYS_EVENTFD2, 0, 0, 0)
	if errno != 0 {
		unix.Close(epollFD)
		return nil, errno
	}
	poller.wfd = int(r0)

	if err = poller.AddRead(poller.wfd); err != nil {
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
	_, err := unix.Write(p.wfd, []byte{0, 0, 0, 0, 0, 0, 0, 1})
	return err
}

func (p *Poller) Polling(callback func(fd int, filter int16) error) error {
	el := newEventList(1024)
	events := make([]unix.EpollEvent, 1024)
	var wakeup bool
	for {
		n, err := unix.EpollWait(p.fd, el.events, -1)
		if err != nil && err != unix.EINTR {
			return err
		}
		for i := 0; i < n; i++ {
			if fd := int(el.events[i].Fd); fd != p.wfd {
				if err := callback(fd, el.events[i].Events); err != nil {
					return err
				}
			} else {
				wakeup = true
				_, _ = unix.Read(p.wfd, p.wfdbuf)
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
	return nil
}

// AddReadWrite ...
func (p *Poller) AddReadWrite(fd int) error {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd),
		Events: unix.EPOLLIN | unix.EPOLLOUT})
}

// AddRead ...
func (p *Poller) AddRead(fd int) error {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: unix.EPOLLIN})
}

// AddWrite ...
func (p *Poller) AddWrite(fd int) error {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: unix.EPOLLOUT})
}

// ModRead ...
func (p *Poller) ModRead(fd int) error {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_MOD, fd, &unix.EpollEvent{Fd: int32(fd), Events: unix.EPOLLIN})
}

// ModReadWrite ...
func (p *Poller) ModReadWrite(fd int) error {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_MOD, fd, &unix.EpollEvent{Fd: int32(fd),
		Events: unix.EPOLLIN | unix.EPOLLOUT})
}

// Delete ...
func (p *Poller) Delete(fd int) error {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_DEL, fd, nil)
}
