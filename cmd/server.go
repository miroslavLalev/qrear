package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/miroslavLalev/qrear/server"
	"github.com/miroslavLalev/qrear/server/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server options",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &config.ServerConfig{}
		err := NewConfigParser(configDir).Parse(cfg)
		if err != nil {
			panic(err)
		}
		lis := server.NewListener(cfg)
		// TODO: find better place
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			lis.Stop()
		}()
		log.Println(lis.Start())
	},
}
