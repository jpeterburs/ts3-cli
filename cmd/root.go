package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ts3",
	Short: "",
	Long:  ``,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrln(err)
		os.Exit(1)
	}
}

func init() {
}
