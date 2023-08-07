// Package cli runs tool as command line util
package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "treco",
	Short: "Test Report Collector",
}

func init() {
	rootCmd.AddCommand(collectCmd)
	rootCmd.AddCommand(serveCmd)
}

// Execute ...
func Execute() error {
	return rootCmd.Execute()
}
