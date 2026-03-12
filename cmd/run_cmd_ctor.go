package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"postman-cli/internal/errs"
	"postman-cli/internal/http_executor"
	"postman-cli/internal/runner"
	"postman-cli/internal/storage"
)

// NewRunCmd constructs the `run` CLI command.
func NewRunCmd() *cobra.Command {
	var envFilePath string
	var noCookies bool
	var clearCookies bool

	c := &cobra.Command{
		Use:   "run [collection.json]",
		Short: "Execute a collection of requests",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			collectionPath := args[0]

			// Load Collection
			collBytes, err := storage.ReadJSONFile(collectionPath)
			if err != nil {
				return errs.Wrap(err, errs.KindInvalidInput, "could not read collection file")
			}

			coll, err := storage.ParseCollection(collBytes)
			if err != nil {
				return errs.Wrap(err, errs.KindInvalidInput, "could not parse collection JSON")
			}

			// Init Runtime Context
			ctx := runner.NewRuntimeContext()

			// Load Environment if provided
			if envFilePath != "" {
				envBytes, err := storage.ReadJSONFile(envFilePath)
				if err != nil {
					return errs.Wrap(err, errs.KindInvalidInput, "could not read environment file")
				}
				env, err := storage.ParseEnvironment(envBytes)
				if err != nil {
					return errs.Wrap(err, errs.KindInvalidInput, "could not parse environment JSON")
				}
				ctx.SetEnvironment(env)
			}

			// Build executor with cookie jar wired in
			exec := http_executor.NewDefaultExecutor()
			if noCookies {
				exec.DisableCookies()
			}

			// Run Collection
			engine := runner.NewCollectionRunner(exec, nil, nil)
			if clearCookies {
				engine.SetClearCookiesPerRequest(true)
			}

			err = engine.Run(coll, ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Execution Failed: %v\n", err)
				os.Exit(1)
			}

			return nil
		},
	}

	c.Flags().StringVarP(&envFilePath, "env", "e", "", "Path to the environment JSON file")
	c.Flags().BoolVar(&noCookies, "no-cookies", false, "Disable cookie persistence for this run")
	c.Flags().BoolVar(&clearCookies, "clear-cookies", false, "Clear cookie jar before each request")

	return c
}

