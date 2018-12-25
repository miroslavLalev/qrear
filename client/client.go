package client

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/miroslavLalev/qrear/message"
)

type Client struct {
	readCb func([]byte)
}

func NewClient(cb func([]byte)) *Client {
	return &Client{readCb: cb}
}

func (cl *Client) Connect() error {
	c, err := net.Dial("unix", "qrear.sock")
	if err != nil {
		return err
	}
	defer c.Close()

	if err := gob.NewEncoder(c).Encode(message.NewExec("ls", []string{"-l"})); err != nil {
		return err
	}

	var res message.Response
	if err := gob.NewDecoder(c).Decode(&res); err != nil {
		return err
	}

	if res.ErrMsg != nil {
		return fmt.Errorf("exec failed: %s", *res.ErrMsg)
	}

	buf := make([]byte, 1024)
	for {
		n, err := c.Read(buf)
		if err != nil {
			return err
		}
		cl.readCb(buf[:n])
	}
}
