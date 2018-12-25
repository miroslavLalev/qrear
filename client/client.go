package client

import (
	"fmt"
	"log"
	"net"
	"time"
)

func Connect() error {
	c, err := net.Dial("unix", "shmux.sock")
	if err != nil {
		return err
	}
	defer c.Close()

	buf := make([]byte, 1024)
	for {
		_, err := c.Write([]byte("Hi from client"))
		if err != nil {
			log.Printf("error during client write: %s\n", err)
			return err
		}

		n, err := c.Read(buf)
		if err != nil {
			log.Printf("error during client read: %s\n", err)
			return err
		}
		t := buf[:n]
		fmt.Printf("server response: %s\n", t)
		time.Sleep(time.Second)
	}
}
