package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var configDir string

var rootCmd = &cobra.Command{
	Use:   "qrear",
	Short: "Shell Multiplexer",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("qrear")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configDir, "config", "c", "", "Directory with config files")
}
