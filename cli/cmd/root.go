package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "kikplate",
	Short: "Kikplate command line interface",
	Long: `Kikplate CLI lets you discover, add, and manage plates from a Kikplate server.

Use it to browse public plates, add them locally, and manage your local plate list.`,
	Example: `  kikplate config init
  kikplate plates add owner/repo
  kikplate plates list
  kikplate help`,
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "Path to config file (default: ~/.kikplate/config.yaml)")
}
func Execute() error {
	return rootCmd.Execute()
}
