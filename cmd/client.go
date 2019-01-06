package cmd

import (
	"log"

	"github.com/miroslavLalev/qrear/client"
	"github.com/miroslavLalev/qrear/client/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clientCmd)
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client options",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &config.ClientConfig{}
		err := NewConfigParser(configDir).Parse(cfg)
		if err != nil {
			panic(err)
		}
		log.Println(client.Start(cfg))
	},
}
