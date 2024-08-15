package chat

type Output struct {
	Message  *Message `json:"message"`
	LogLevel string   `json:"log_level"`
}

type Message struct {
	Role         string                   `json:"role"`
	Content      string                   `json:"content"`
	Name         string                   `json:"name,omitempty"`
	FunctionCall *ChatMessageFunctionCall `json:"function_call,omitempty"`
	Confirmation *Confirmation            `json:"copilot_confirmation"`
	References   []Reference              `json:"copilot_references"`
	Errors       []CopilotError           `json:"copilot_errors"`
}

type Completion struct {
	Choices []CompletionChoice `json:"choices"`
}

type CompletionChoice struct {
	Delta Message `json:"delta"`
}

type Request struct {
	CopilotThreadID string    `json:"copilot_thread_id"`
	Messages        []Message `json:"messages"`
	Agent           string    `json:"agent"`
}

type ChatMessageFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Confirmation struct {
	Type         string `json:"type"`
	Title        string `json:"title"`
	Message      string `json:"message"`
	Confirmation any    `json:"confirmation"`
}

type Reference struct {
	Type     string            `json:"type"`
	ID       string            `json:"id"`
	Data     any               `json:"data"`
	Metadata ReferenceMetadata `json:"metadata"`
}

type ReferenceMetadata struct {
	DisplayName string `json:"display_name"`
	DisplayIcon string `json:"display_icon"`
	DisplayURL  string `json:"display_url"`
}

type CopilotError struct {
	Type       string `json:"type"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Identifier string `json:"identifier"`
}
