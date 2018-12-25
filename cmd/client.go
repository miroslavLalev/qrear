package cmd

import (
	"log"

	"github.com/miroslavLalev/qrear/cli"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clientCmd)
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client options",
	Run: func(cmd *cobra.Command, args []string) {
		// log.Println(cli.Connect())
		log.Println(cli.Start())
	},
}
