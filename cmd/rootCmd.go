// rootCmd.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	// "github.com/github-technology-partners/gh-debug-cli/cmd/stream"
)

var rootCmd = &cobra.Command{
	Use:   "gh-debug-cli",
	Short: "A CLI tool for debugging",
	Long:  `This CLI tool allows you to debug your agent by chatting with it locally.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'gh-debug-cli --help' to see available commands")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands to rootCmd
	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(streamCmd)
}