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
	Use:   "stream [file]",
	Short: "Stream data to your agent",
	Long:  `The stream command allows you to initiate a data stream to your agent.`,
	Run:   agentStream,
}

func init() {
	streamCmd.PersistentFlags().String(streamCmdFileFlag, "", "Parse agent responses from a file")
}

func agentStream(cmd *cobra.Command, args []string) {
	fmt.Println("stream command executed successfully")

	file := args[0]

	result, err := stream.ParseFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(result)
}
