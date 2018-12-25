package conn

import (
	"fmt"
	"log"
	"net"
)

func Handle(c net.Conn) {
	defer c.Close()

	buf := make([]byte, 1024)
	for {
		n, err := c.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}
		t := buf[:n]
		fmt.Fprintf(c, "client send: %s", t)
	}
}
