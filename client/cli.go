package client

import (
	"log"
	"os"
	"time"

	"github.com/miroslavLalev/qrear/client/internal/screen"
)

func Start() error {
	logFile, err := os.OpenFile("dev.log", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}

	logger := log.New(logFile, "", log.Ltime|log.Ldate)

	waitCh := make(chan time.Duration, 1)
	s := screen.New(logger, waitCh)

	go func() {
		select {
		case d := <-waitCh:
			time.Sleep(d)

			wrFn, err := s.AddTab("sometab")
			if err != nil {
				logger.Printf("error adding tab: %s", err)
				return
			}

			s.AddTab("test")

			err = NewClient(func(b []byte) { wrFn(b) }).Connect()
			if err != nil {
				logger.Println(err)
				return
			}
		}
	}()
	return s.Init()
}
