package chat

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAndEmit(t *testing.T) {
	testCases := []struct {
		name          string
		stream        string
		expectedAny   []any
		expectedTypes []interface{}
		expectedError error
	}{
		{
			name: "happy_path",
			stream: `event: copilot_confirmation
data: {"type": "confirm", "title": "Test Confirmation", "message": "This is a test confirmation"}

event: copilot_references
data: [{"type": "ref", "id": "1", "metadata": {"display_name": "Test Reference"}}]

event: copilot_errors
data: [{"type": "error", "code": "E1", "message": "Test Error", "identifier": "1"}]

data: {"choices":[{"delta":{"content":"ahoy there"}}]}

				`,
			expectedTypes: []interface{}{
				Confirmation{},
				[]Reference{},
				[]CopilotError{},
				Completion{},
			},
			expectedAny: []any{
				Confirmation{
					Type:    "confirm",
					Title:   "Test Confirmation",
					Message: "This is a test confirmation",
				},
				[]Reference{
					{
						Type: "ref",
						ID:   "1",
						Metadata: ReferenceMetadata{
							DisplayName: "Test Reference",
						},
					},
				},
				[]CopilotError{
					{
						Type:       "error",
						Code:       "E1",
						Message:    "Test Error",
						Identifier: "1",
					},
				},
				Completion{
					Choices: []CompletionChoice{
						{
							Delta: Message{
								Content: "ahoy there",
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "happy_path_arrays",
			stream: `event: copilot_confirmation
data: {"type": "confirm", "title": "Test Confirmation", "message": "This is a test confirmation"}

event: copilot_references
data: [{"type": "ref", "id": "1", "metadata": {"display_name": "Test Reference"}}, {"type": "ref", "id": "2", "metadata": {"display_name": "Test Reference"}}]

event: copilot_errors
data: [{"type": "error", "code": "E1", "message": "Test Error", "identifier": "1"}, {"type": "error", "code": "E1", "message": "Test Error", "identifier": "2"}]

data: {"choices":[{"delta":{"content":"ahoy there"}}]}

				`,
			expectedTypes: []interface{}{
				Confirmation{},
				[]Reference{},
				[]CopilotError{},
				Completion{},
			},
			expectedAny: []any{
				Confirmation{
					Type:    "confirm",
					Title:   "Test Confirmation",
					Message: "This is a test confirmation",
				},
				[]Reference{
					{
						Type: "ref",
						ID:   "1",
						Metadata: ReferenceMetadata{
							DisplayName: "Test Reference",
						},
					},
					{
						Type: "ref",
						ID:   "2",
						Metadata: ReferenceMetadata{
							DisplayName: "Test Reference",
						},
					},
				},
				[]CopilotError{
					{
						Type:       "error",
						Code:       "E1",
						Message:    "Test Error",
						Identifier: "1",
					},
					{
						Type:       "error",
						Code:       "E1",
						Message:    "Test Error",
						Identifier: "2",
					},
				},
				Completion{
					Choices: []CompletionChoice{
						{
							Delta: Message{
								Content: "ahoy there",
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "failure_mismatched_event_types",
			stream: `event: copilot_references
data: {"type": "confirm", "title": "Test Confirmation", "message": "This is a test confirmation"}

		`,
			expectedError: fmt.Errorf("\x1b[31m\nAlas...The following is not a valid copilot reference:\n[{\"type\": \"confirm\", \"title\": \"Test Confirmation\", \"message\": \"This is a test confirmation\"}]\n\nErrors:\nensure data is an array of copilot_references\n\n\x1b[37m"),
		},
		{
			name: "failure_invalid_event_type",
			stream: `retry: copilot_references

		`,
			expectedError: fmt.Errorf("\x1b[31monly 'event' and 'data' fields are supported, found: retry\n\n\x1b[37m"),
		},
		{
			name: "failure_missing_copilot_reference_metadata",
			stream: `event: copilot_references
data: [{"type": "", "id": "", "metadata": {"display_name": ""}}]

				`,
			expectedError: fmt.Errorf("\x1b[31m\nAlas...The following is not a valid copilot reference:\n[[{\"type\": \"\", \"id\": \"\", \"metadata\": {\"display_name\": \"\"}}]]\n\nErrors:\nref 0 is missing a type\nref 0 is missing an id\nref 0 is missing a metadata display name\n\n\n\x1b[37m"),
		},
		{
			name: "failure_missing_data_field",
			stream: `event: copilot_references
[{"type": "", "id": "", "metadata": {"display_name": ""}}]

				`,
			expectedError: fmt.Errorf("\x1b[31monly 'event' and 'data' fields are supported, found: [{\"type\"\n\n\x1b[37m"),
		},
		{
			name: "failure_references_not_array",
			stream: `event: copilot_references
data: {"type": "ref", "id": "1", "metadata": {"display_name": "Test Reference"}}

				`,
			expectedError: fmt.Errorf("\x1b[31m\nAlas...The following is not a valid copilot reference:\n[{\"type\": \"ref\", \"id\": \"1\", \"metadata\": {\"display_name\": \"Test Reference\"}}]\n\nErrors:\nensure data is an array of copilot_references\n\n\x1b[37m"),
		},
		{
			name: "failure_invalid_message_chunk",
			stream: `data: {"copilot_confirmation": {"type":"action","title":"Turn off feature flag","message":"Are you sure you wish to turn off the feature flag?","confirmation":{"id":"id-123"}}}

				`,
			expectedError: fmt.Errorf("\x1b[31m\nAlas...Failed to process data fields:\n[{\"copilot_confirmation\": {\"type\":\"action\",\"title\":\"Turn off feature flag\",\"message\":\"Are you sure you wish to turn off the feature flag?\",\"confirmation\":{\"id\":\"id-123\"}}}]\n\nErrors: setting confirmation in a message payload is not supported\n\n\n\x1b[37m"),
		},
		{
			name: "failure_invalid_error_type",
			stream: `event: error
data: {"type":"function","code":"foo","message":"A function error occurred","identifier":"fn123"}

		`,
			expectedError: fmt.Errorf("\x1b[31m\nAlas...The following is not a valid event:\n[{\"type\":\"function\",\"code\":\"foo\",\"message\":\"A function error occurred\",\"identifier\":\"fn123\"}]\n\nErrors:\ntype not supported: error\n\n\x1b[37m"),
		},
		{
			name: "failure_invalid_error_must_be_array",
			stream: `event: copilot_errors
data: {"type":"function","code":"foo","message":"A function error occurred","identifier":"fn123"}

		`,
			expectedError: fmt.Errorf("\x1b[31m\nAlas...The following is not a valid a copilot error:\n[{\"type\":\"function\",\"code\":\"foo\",\"message\":\"A function error occurred\",\"identifier\":\"fn123\"}]\n\nErrors:\nensure data is an array of copilot_errors\n\n\x1b[37m"),
		},
		{
			name: "failure_extra_double_quotes",
			stream: `event: copilot_errors
data: [{"type":"reference","code":"foo","message":"A reference error occurred","identifier":"ref123"},{"type":"function","code":"foo","message":"A function error occurred","identifier":"fn123"},{"type":"agent","code":"foo","message":"An agent error occurred","identifier":"agt123"}]

"data: [DONE]"

		`,
			expectedError: fmt.Errorf("\x1b[31monly 'event' and 'data' fields are supported, found: \"data\n\n\x1b[37m"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var emittedData []any
			dataEmitter := func(data any) {
				emittedData = append(emittedData, data)
			}

			p := NewParser(bytes.NewBufferString(tc.stream), dataEmitter)
			err := p.ParseAndEmit(context.Background(), LEVEL_DEBUG)

			assert.Equal(t, tc.expectedError, err)

			if tc.expectedAny != nil {
				if emittedData != nil {
					assert.Equal(t, len(tc.expectedAny), len(emittedData))
				} else {
					assert.Fail(t, "emittedData is nil")
				}
			}

			for i, expectedType := range tc.expectedTypes {
				switch v := expectedType.(type) {
				case Confirmation:
					actualConfirmation, ok := emittedData[i].(Confirmation)
					assert.True(t, ok)

					expectedConfirmation, ok := tc.expectedAny[i].(Confirmation)
					assert.True(t, ok)

					assert.Equal(t, expectedConfirmation.Type, actualConfirmation.Type)
					assert.Equal(t, expectedConfirmation.Title, actualConfirmation.Title)
					assert.Equal(t, expectedConfirmation.Message, actualConfirmation.Message)
				case []Reference:
					actualReferences, ok := emittedData[i].([]Reference)
					assert.True(t, ok)

					expectedReferences, ok := tc.expectedAny[i].([]Reference)
					assert.True(t, ok)

					assert.Equal(t, len(actualReferences), len(expectedReferences))

					for j := range actualReferences {
						assert.Equal(t, actualReferences[j].Type, expectedReferences[j].Type)
						assert.Equal(t, actualReferences[j].ID, expectedReferences[j].ID)
						assert.Equal(t, actualReferences[j].Metadata.DisplayName, expectedReferences[j].Metadata.DisplayName)
					}
				case []CopilotError:
					actualCopilotErrors, ok := emittedData[i].([]CopilotError)
					assert.True(t, ok)

					expectedCopilotErrors, ok := tc.expectedAny[i].([]CopilotError)
					assert.True(t, ok)

					assert.Equal(t, len(actualCopilotErrors), len(expectedCopilotErrors))

					for j := range actualCopilotErrors {
						assert.Equal(t, actualCopilotErrors[j].Type, actualCopilotErrors[j].Type)
						assert.Equal(t, actualCopilotErrors[j].Code, actualCopilotErrors[j].Code)
						assert.Equal(t, actualCopilotErrors[j].Message, actualCopilotErrors[j].Message)
						assert.Equal(t, actualCopilotErrors[j].Identifier, actualCopilotErrors[j].Identifier)
					}
				case Completion:
					actualChat, ok := emittedData[i].(Completion)
					assert.True(t, ok)

					expectedChat, ok := tc.expectedAny[i].(Completion)
					assert.True(t, ok)

					for j, actualChoice := range actualChat.Choices {
						assert.Equal(t, expectedChat.Choices[j].Delta.Role, actualChoice.Delta.Role)
						assert.Equal(t, expectedChat.Choices[j].Delta.Content, actualChoice.Delta.Content)
						assert.Equal(t, expectedChat.Choices[j].Delta.Name, actualChoice.Delta.Name)
					}

				default:
					assert.Fail(t, fmt.Sprintf("unexpected type %T", v))
				}
			}
		})
	}

}
