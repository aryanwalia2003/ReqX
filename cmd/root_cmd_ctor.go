package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd constructs the base CLI command.
func NewRootCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "postman-cli", //this is the name of the command that will be used to run the CLI
		Short: "A fast, scriptable API client for the terminal",
		Long: `🚀 postman-cli: A high-performance, developer-first API client.
Built for speed and simplicity, it allows you to run Postman-style collections, 
individual requests, and Socket.IO events directly from your command line.

Key Features:
- Run collections with environment-based variable injection
- Permanent collection management (add, move, list requests)
- Real-time Socket.IO event testing
- JavaScript scripting for pre-request logic and test assertions
- Temporary request injection during test runs`,
		Example: `  # Run a collection
  postman-cli run collection.json -e env.json
  
  # Run a specific request from a collection
  postman-cli run collection.json -f "Login"
  
  # Manage your collection file
  postman-cli collection list my-api.json
  postman-cli collection add my-api.json -n "New Req" -u "http://api.com/get"`,
	}

	c.AddCommand(NewRunCmd())
	c.AddCommand(NewSampleCmd())
	c.AddCommand(NewReqCmd())
	c.AddCommand(NewSioCmd())
	c.AddCommand(NewCollectionCmd())

	return c
}

// Execute is the main entrypoint called by main.go.
func Execute() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
