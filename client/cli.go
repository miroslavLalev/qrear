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
			wFn1, err := s.AddTab("sometab")
			if err != nil {
				logger.Printf("error adding tab: %s", err)
				return
			}
			wFn2, err := s.AddTab("othertab")
			if err != nil {
				logger.Printf("error adding tab: %s", err)
				return
			}
			for {
				wFn1([]byte("text\n"))
				wFn2([]byte("moreText\n"))
				time.Sleep(time.Second)
			}
		}
	}()

	err = s.Init()
	if err != nil {
		return err
	}

	// if err := s.Init(); err != nil {
	// 	return err
	// }

	return nil
}
