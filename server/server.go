package server

import (
	"net"

	"github.com/miroslavLalev/qrear/server/config"
	"github.com/miroslavLalev/qrear/server/internal/conn"
)

type Listener struct {
	netLis net.Listener

	cfg *config.ServerConfig
}

func NewListener(cfg *config.ServerConfig) *Listener {
	return &Listener{cfg: cfg}
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
