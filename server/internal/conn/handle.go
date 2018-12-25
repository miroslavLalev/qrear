package conn

import (
	"encoding/gob"
	"fmt"
	"net"
	"os/exec"

	"github.com/miroslavLalev/qrear/message"
)

func Handle(c net.Conn) {
	defer c.Close()

	var e message.Exec
	err := gob.NewDecoder(c).Decode(&e)
	if err != nil {
		gob.NewEncoder(c).Encode(message.NewResponse(err.Error()))
		return
	}

	err = gob.NewEncoder(c).Encode(message.Response{})
	if err != nil {
		gob.NewEncoder(c).Encode(message.NewResponse(err.Error()))
		return
	}

	cmd := exec.Command(e.Command, e.Args...)
	cmd.Stdout = c
	cmd.Stderr = c

	err = cmd.Start()
	if err != nil {
		fmt.Printf("err during start: %s\n", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("err during wait: %s\n", err)
		return
	}
}
