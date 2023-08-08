package cli

import (
	"treco/server"

	"github.com/spf13/cobra"
)

// newServeCommand
func newServeCommand() *cobra.Command {
	var port int
	var cfgFile string

	serveCmd := cobra.Command{
		Use:   "serve",
		Short: "Runs as a web server",
		Run: func(cmd *cobra.Command, args []string) {
			server.Start(cfgFile, port)
		},
	}

	flags := serveCmd.Flags()
	flags.IntVarP(&port, "port", "p", 8080, "port for server to run")
	flags.StringVarP(&cfgFile, "config", "c", "", "config file")

	return &serveCmd
}
