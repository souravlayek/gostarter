package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	webServerFlag    bool
	includeTestsFlag bool
	includeAirFlag   bool
	middlewaresFlag  string
	prefixFlag       string
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "gostarter",
		Short: "A CLI tool to setup fresh Go projects",
		Long:  `gostarter is a CLI application that helps you setup fresh Go projects from scratch.`,
	}

	var createCmd = &cobra.Command{
		Use:   "create [project-name]",
		Short: "Create a new Go project",
		Long:  `Create command creates a new Go project with a standard directory structure and files.`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// For now, always include tests when using CLI
			createProject(args, webServerFlag, includeTestsFlag, includeAirFlag, middlewaresFlag, prefixFlag)
		},
	}

	// Add flags to the create command
	createCmd.Flags().BoolVar(&webServerFlag, "web-server", false, "Create a web server project")
	createCmd.Flags().BoolVar(&includeTestsFlag, "include-tests", true, "Include test files in the project")
	createCmd.Flags().BoolVar(&includeAirFlag, "air", false, "Setup Air for auto-reloading (web server projects only)")
	createCmd.Flags().StringVar(&middlewaresFlag, "middlewares", "", "Comma-separated middleware list (realIP,logger,ratelimit)")
	createCmd.Flags().StringVar(&prefixFlag, "prefix", "", "Module prefix (e.g., github.com/username/)")

	var addPkgCmd = &cobra.Command{
		Use:   "add-pkg [package-name]",
		Short: "Add a package to the current project",
		Long:  `Add-pkg command adds a predefined package to the current project's pkg directory.`,
		Args:  cobra.ArbitraryArgs, // Changed to allow zero args for the list flag
		Run: func(cmd *cobra.Command, args []string) {
			list, _ := cmd.Flags().GetBool("list")
			if list {
				listPackages()
				return
			}

			if len(args) == 0 {
				fmt.Println("Error: package name is required")
				cmd.Help()
				return
			}

			addPackage(args[0])
		},
	}

	addPkgCmd.Flags().BoolP("list", "l", false, "List all available packages")

	var addMiddlewareCmd = &cobra.Command{
		Use:   "add-middleware [middleware-name]",
		Short: "Add a middleware to the current project",
		Long:  `Add-middleware command adds a predefined middleware to the current project's internal/middleware directory.`,
		Args:  cobra.ArbitraryArgs, // Changed to allow zero args for the list flag
		Run: func(cmd *cobra.Command, args []string) {
			list, _ := cmd.Flags().GetBool("list")
			if list {
				listMiddlewares()
				return
			}

			if len(args) == 0 {
				fmt.Println("Error: middleware name is required")
				cmd.Help()
				return
			}

			addMiddleware(args[0])
		},
	}

	addMiddlewareCmd.Flags().BoolP("list", "l", false, "List all available middleware")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(addPkgCmd)
	rootCmd.AddCommand(addMiddlewareCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
