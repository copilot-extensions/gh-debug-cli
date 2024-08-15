package chat

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChat(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		username      string
		token         string
		expectedError error
	}{
		{
			name:          "happy_path",
			url:           "http://localhost:8080",
			username:      "username",
			token:         "token",
			expectedError: nil,
		},
		{
			name:          "failure_agent_url_empty",
			url:           "",
			username:      "username",
			token:         "token",
			expectedError: fmt.Errorf("agent url is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualError := Chat(tt.url, tt.username, tt.token, LEVEL_NONE)
			assert.Equal(t, tt.expectedError, actualError)
		})
	}
}

func TestOutput_String(t *testing.T) {
	tests := []struct {
		name           string
		output         *Output
		expectedString string
	}{
		{
			name: "happy_path_function_call",
			output: &Output{
				Message: &Message{
					FunctionCall: &ChatMessageFunctionCall{
						Name:      "test",
						Arguments: "args",
					},
				},
				LogLevel: LEVEL_DEBUG,
			},
			expectedString: "\x1b[32m\nHuzzah! You successfully received a function call!\n\x1b[37m\x1b[32m╔═════════════╤════════╗\n║     Key     │ Value  ║\n╟━━━━━━━━━━━━━┼━━━━━━━━╢\n║ role        │        ║\n║ name        │ test   ║\n║ arguments   │ args   ║\n╟━━━━━━━━━━━━━┼━━━━━━━━╢\n║ Parsed function data ║\n╚═════════════╧════════╝\x1b[37m\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualString := tt.output.String()
			assert.Equal(t, tt.expectedString, actualString)
		})
	}
}
