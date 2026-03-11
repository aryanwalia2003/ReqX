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
		Short: "A fast, scriptable API client for the command line",
		Long:  "postman-cli is a lightweight developer tool for running requests and debugging APIs directly from the terminal.",
		RunE: func(cmd *cobra.Command, args []string) error { 
			return cmd.Help()
		}, //this is the function that will be called when the command is run for eg. when we type postman-cli in the terminal it will print the help message 
	}
	
	c.AddCommand(NewRunCmd())
	
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
