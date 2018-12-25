package conn

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"

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

	for {
		fmt.Fprintln(c, "some output...")
		time.Sleep(500 * time.Millisecond)
	}
}
