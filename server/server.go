package server

import (
	"net"

	"github.com/miroslavLalev/qrear/server/internal/conn"
)

type Listener struct {
	netLis net.Listener
}

func NewListener() *Listener {
	return &Listener{}
}

func (l *Listener) Start() error {
	list, err := net.Listen("unix", "qrear.sock")
	if err != nil {
		return err
	}
	l.netLis = list
	defer list.Close()

	for {
		c, err := list.Accept()
		if err != nil {
			return err
		}
		go conn.Handle(c)
	}
}

func (l *Listener) Stop() {
	l.netLis.Close()
}
