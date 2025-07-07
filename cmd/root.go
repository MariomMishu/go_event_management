package cmd

import (
	"fmt"
	"os"

	"ems/config"
	"ems/conn"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use: "app",
	}
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

// Execute executes the root command
func Execute() {
	// load config
	config.LoadConfig()

	// Initialize logger
	initLogger()

	conn.ConnectDb()
	conn.ConnectRedis()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initLogger() {
	fmt.Println("Initializing logger...")
}
