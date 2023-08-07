package cli

import (
	"log"
	"treco/conf"
	"treco/server"

	"github.com/spf13/cobra"
)

var port int
var cfgFile string

// serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs as a web server",
	Run: func(cmd *cobra.Command, args []string) {
		server.Start(port)
	},
}

// init
func init() {
	flags := serveCmd.Flags()
	flags.IntVarP(&port, "port", "p", 8080, "port for server to run")
	flags.StringVarP(&cfgFile, "config", "c", "", "config file")

	if cfgFile != "" {
		if err := conf.LoadEnvFromFile(cfgFile); err != nil {
			log.Fatalf("error occured while loading from config %v\n", err)
			return
		}
	} else {
		log.Println("no config file path set")
	}
}
