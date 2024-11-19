// stream.go
package cmd

import (
	"fmt"
	"os"

	"github.com/github-technology-partners/gh-debug-cli/pkg/stream"
	"github.com/spf13/cobra"
)

const (
	streamCmdFileFlag = "file"
)

// streamCmd represents the new command for streaming functionality
var streamCmd = &cobra.Command{
	Use:   "stream --file [filename]",
	Short: "Parse stream data from agent",
	Long:  `Allows you to parse a data stream to your agent response.`,
	Run:   agentStream,
}

func init() {
	streamCmd.PersistentFlags().String(streamCmdFileFlag, "", "Parse agent responses from a file")
	rootCmd.AddCommand(streamCmd)
}

func agentStream(cmd *cobra.Command, args []string) {
	fmt.Println("stream command executed successfully")

	file, _ := cmd.Flags().GetString(streamCmdFileFlag)
	if file == "" {
		fmt.Fprintln(os.Stderr, "Error: --file [file] is required")
		os.Exit(1)
	}

	result, err := stream.ParseFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(result)
}
