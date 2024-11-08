package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/github-technology-partners/gh-debug-cli/pkg/chat"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	chatCmdURLFlag        = "url"
	chatCmdUsernameFlag   = "username"
	chatCmdLogLevelFlag   = "log-level"
	chatCmdTokenFlag      = "token"
	chatCmdPrivateKeyFlag = "private-key"
	chatCmdPublicKeyFlag  = "public-key"
)

var chatCmd = &cobra.Command{
	Use:              "chat",
	Short:            "Interact with your agent.",
	Long:             `This cli tool allows you to debug your agent by chatting with it locally.`,
	Run:              agentChat,
	TraverseChildren: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			optName := strings.ToUpper(f.Name)
			optName = strings.ReplaceAll(optName, "-", "_")
			if val, ok := os.LookupEnv(optName); !f.Changed && ok {
				fmt.Printf("Setting %s to %s\n", f.Name, val)
				err2 := f.Value.Set(val)
				if err2 != nil {
					err = fmt.Errorf("invalid environment variable %s: %w", optName, err2)
				}
			}
		})
		return err
	}}

func init() {
	chatCmd.CompletionOptions.DisableDefaultCmd = true

	chatCmd.PersistentFlags().String(chatCmdURLFlag, "http://localhost:8080", "url to chat with your agent")
	chatCmd.PersistentFlags().String(chatCmdUsernameFlag, "sparklyunicorn", "username to display in chat")
	chatCmd.PersistentFlags().String(chatCmdTokenFlag, "", "GitHub token for chat authentication (optional)")
	chatCmd.PersistentFlags().String(chatCmdLogLevelFlag, "DEBUG", "Log level to help debug events. Supported types are `DEBUG`, `TRACE`, `NONE`. `DEBUG` returns general logs. `TRACE` prints the raw http response.")
	chatCmd.PersistentFlags().String(chatCmdPrivateKeyFlag, "", "Private key for payload verification")
	chatCmd.PersistentFlags().String(chatCmdPublicKeyFlag, "", "Public key for payload verification")

}

func agentChat(cmd *cobra.Command, args []string) {

	url, _ := cmd.Flags().GetString(chatCmdURLFlag)
	if url == "" {
		fmt.Println("a url is required to chat with your agent")
	}

	username, _ := cmd.Flags().GetString(chatCmdUsernameFlag)

	token, _ := cmd.Flags().GetString(chatCmdTokenFlag)

	debug, _ := cmd.Flags().GetString(chatCmdLogLevelFlag)
	debug = strings.ToUpper(debug)
	if debug != chat.LEVEL_NONE && debug != chat.LEVEL_DEBUG && debug != chat.LEVEL_TRACE {
		fmt.Println("debug mode must be either `DEBUG`, `TRACE`, or `NONE`")
	}

	err := chat.Chat(url, username, token, debug)
	if err != nil {
		fmt.Println(err)
	}
}
