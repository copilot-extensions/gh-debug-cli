// rootCmd.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "A CLI tool for debugging",
	Long:  `This CLI tool allows you to debug your agent by chatting with it locally.`,
	Run: func(cmd *cobra.Command, args []string) {
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
