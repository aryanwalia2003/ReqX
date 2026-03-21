package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"reqx/internal/history"
	"reqx/internal/ui"
)

func NewUICmd() *cobra.Command {
	var port int

	c := &cobra.Command{
		Use:   "ui",
		Short: "Open the local test history dashboard in your browser",
		RunE: func(_ *cobra.Command, _ []string) error {
			db, err := history.Open()
			if err != nil {
				color.Red("✗ Could not open history database: %v\n", err)
				os.Exit(1)
			}
			defer db.Close()

			fmt.Println()
			color.Cyan("📊 Starting ReqX UI...")
			server := ui.NewServer(db, port)
			return server.Start()
		},
	}

	c.Flags().IntVarP(&port, "port", "p", 8090, "Port to serve the dashboard on")
	return c
}
