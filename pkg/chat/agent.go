package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/google/uuid"
)

func invokeAgent(ctx context.Context, url string, token string, history []Message, debugMode string) ([]*Message, error) {
	copilotThreadID := uuid.New().String()
	body := Request{
		Messages:        history,
		CopilotThreadID: copilotThreadID,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Github-Public-Key-Signature", "")
	req.Header.Set("Github-Public-Key-Identifier", "")

	if token != "" {
		req.Header.Set("X-GitHub-Token", token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	var buf messageBuffer
	fn := func(data any) {
		switch v := data.(type) {
		case Completion:
			buf.WriteChatMessage(v)

		case Confirmation:
			buf.WriteConfirmation(v)

		case []Reference:
			buf.WriteReferences(v)

		case []CopilotError:
			buf.WriteErrors(v)

		default:
			fmt.Printf("Invalid data type: %T\n", v)
		}
	}

	if shouldLog(debugMode, LEVEL_TRACE) {
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(yellow("Raw Response\n" + string(respDump) + "\n\n"))
	}

	parser := NewParser(resp.Body, fn)
	if err := parser.ParseAndEmit(ctx, debugMode); err != nil {
		fmt.Println(err)
	}

	if parser.ValidEventCount() {
		return nil, fmt.Errorf("cannot have more than one event type in an invocation, found %d", parser.eventCount)
	}

	return buf, nil
}
