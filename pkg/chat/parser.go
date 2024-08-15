package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/jclem/sseparser"
)

const (
	sseDataField  = "data"
	sseEventField = "event"
)

type dataEmitter func(data any)

// Parser is a parser for ServerSent Events (SSE).
type Parser struct {
	buf        io.Reader
	fn         dataEmitter
	eventCount int
}

// NewParser creates a new SSEParser.
func NewParser(buf io.Reader, fn dataEmitter) *Parser {
	return &Parser{
		buf:        buf,
		fn:         fn,
		eventCount: 0,
	}
}

// ParseAndEmit parses the SSE stream and emits the parsed events.
func (p *Parser) ParseAndEmit(ctx context.Context, debug string) error {
	scanner := sseparser.NewStreamScanner(p.buf)

	for {
		event, _, err := scanner.Next()
		if err != nil {
			if errors.Is(err, sseparser.ErrStreamEOF) {
				return nil
			}
			return fmt.Errorf("failed to read from stream: %w", err)
		}

		eventFields := map[string]string{}
		dataFields := []string{}
		for _, field := range event.Fields() {
			switch field.Name {
			case sseEventField:
				eventFields[field.Name] = field.Value
			case sseDataField:
				dataFields = append(dataFields, field.Value)
			default:
				return fmt.Errorf(red("only 'event' and 'data' fields are supported, found: %s\n\n"), field.Name)
			}
		}

		switch {
		case eventFields[sseEventField] == "copilot_confirmation":
			p.eventCount++
			err := emitConfirmation(dataFields, p.fn)
			if err != nil && shouldLog(debug, LEVEL_DEBUG) {
				return fmt.Errorf(red("\nAlas...The following is not a valid copilot confirmation:\n%v\n\nErrors:\n%v\n\n"), dataFields, err)
			}

		case eventFields[sseEventField] == "copilot_references":
			p.eventCount++
			err := emitReferences(dataFields, p.fn)
			if err != nil && shouldLog(debug, LEVEL_DEBUG) {
				return fmt.Errorf(red("\nAlas...The following is not a valid copilot reference:\n%v\n\nErrors:\n%v\n\n"), dataFields, err)
			}

		case eventFields[sseEventField] == "copilot_errors":
			p.eventCount++
			err := emitErrors(dataFields, p.fn)
			if err != nil && shouldLog(debug, LEVEL_DEBUG) {
				return fmt.Errorf(red("\nAlas...The following is not a valid a copilot error:\n%v\n\nErrors:\n%v\n\n"), dataFields, err)
			}

		case eventFields[sseEventField] != "":
			if shouldLog(debug, LEVEL_DEBUG) {
				err := fmt.Errorf("type not supported: %s", eventFields[sseEventField])
				return fmt.Errorf(red("\nAlas...The following is not a valid event:\n%v\n\nErrors:\n%v\n\n"), dataFields, err)
			}

		default:
			if _, ok := eventFields[sseEventField]; ok && shouldLog(debug, LEVEL_DEBUG) {
				err := fmt.Errorf("event field must have a type")
				return fmt.Errorf(red("\nAlas...The following is not a valid event:\n%v\n\nErrors:\n%v\n\n"), dataFields, err)
			}

			err := emitDatas(dataFields, p.fn)
			if err != nil && shouldLog(debug, LEVEL_DEBUG) {
				return fmt.Errorf(red("\nAlas...Failed to process data fields:\n%v\n\nErrors: %v\n\n"), dataFields, err)
			}
		}
	}
}

func (p *Parser) ValidEventCount() bool {
	return p.eventCount > 1
}

func emitErrors(data []string, fn dataEmitter) error {
	for _, d := range data {
		var errs []CopilotError
		if err := json.Unmarshal([]byte(d), &errs); err != nil {
			return fmt.Errorf("ensure data is an array of copilot_errors")
		}

		if len(errs) == 0 {
			return fmt.Errorf("no errors found")
		}

		var errMsg strings.Builder
		for i, err := range errs {
			if err.Type == "" {
				errMsg.WriteString(fmt.Sprintf("error %d is missing a type\n", i))
			}
			if err.Code == "" {
				errMsg.WriteString(fmt.Sprintf("error %d is missing a code\n", i))
			}
			if err.Message == "" {
				errMsg.WriteString(fmt.Sprintf("error %d is missing a message\n", i))
			}
			if err.Identifier == "" {
				errMsg.WriteString(fmt.Sprintf("error %d is missing an identifier\n", i))
			}
		}

		if errMsg.Len() > 0 {
			return fmt.Errorf(errMsg.String())
		}

		fn(errs)
	}

	return nil
}

func emitReferences(data []string, fn dataEmitter) error {
	for _, d := range data {
		var refs []Reference
		if err := json.Unmarshal([]byte(d), &refs); err != nil {
			return fmt.Errorf("ensure data is an array of copilot_references")
		}

		if len(refs) == 0 {
			return fmt.Errorf("no references found")
		}

		var errMsg strings.Builder
		for i, ref := range refs {
			if ref.Type == "" {
				errMsg.WriteString(fmt.Sprintf("ref %d is missing a type\n", i))
			}
			if ref.ID == "" {
				errMsg.WriteString(fmt.Sprintf("ref %d is missing an id\n", i))
			}
			if ref.Metadata.DisplayName == "" {
				errMsg.WriteString(fmt.Sprintf("ref %d is missing a metadata display name\n", i))
			}
		}

		if errMsg.Len() > 0 {
			return fmt.Errorf(errMsg.String())
		}

		fn(refs)
	}

	return nil
}

func emitConfirmation(data []string, fn dataEmitter) error {
	for _, d := range data {
		var confirmation Confirmation
		if err := json.Unmarshal([]byte(d), &confirmation); err != nil {
			return fmt.Errorf("ensure data is of type copilot_confirmation")
		}

		var errMsg strings.Builder
		if confirmation.Type == "" {
			errMsg.WriteString("confirmation is missing a type\n")
		}
		if confirmation.Title == "" {
			errMsg.WriteString("confirmation is missing a title\n")
		}
		if confirmation.Message == "" {
			errMsg.WriteString("confirmation is missing a message\n")
		}

		if errMsg.Len() > 0 {
			return fmt.Errorf(errMsg.String())
		}

		fn(confirmation)
	}

	return nil
}

func emitDatas(datas []string, fn dataEmitter) error {
	for _, data := range datas {
		if data == "" || data == "[DONE]" {
			continue
		}

		var message Message
		if err := json.Unmarshal([]byte(data), &message); err == nil {
			var errMsg strings.Builder

			if message.Confirmation != nil {
				errMsg.WriteString("setting confirmation in a message payload is not supported\n")
			}

			if message.Errors != nil {
				errMsg.WriteString("setting errors in a message payload is not supported\n")
			}

			if message.References != nil {
				errMsg.WriteString("setting references in a message payload is not supported\n")
			}

			if errMsg.Len() > 0 {
				return fmt.Errorf(errMsg.String())
			}
		}

		var chatMessage Completion
		if err := json.Unmarshal([]byte(data), &chatMessage); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		fn(chatMessage)
	}

	return nil
}

type messageBuffer []*Message

func (mb *messageBuffer) lastMessage() *Message {
	if len(*mb) == 0 {
		*mb = append(*mb, new(Message))
	}

	buf := *mb
	return buf[len(buf)-1]
}

func (mb *messageBuffer) WriteConfirmation(c Confirmation) {
	if last := mb.lastMessage(); last != nil {
		last.Confirmation = &c
	}
}

func (mb *messageBuffer) WriteReferences(r []Reference) {
	if last := mb.lastMessage(); last != nil {
		last.References = r
	}
}

func (mb *messageBuffer) WriteErrors(e []CopilotError) {
	if last := mb.lastMessage(); last != nil {
		last.Errors = e
	}
}

func (mb *messageBuffer) WriteChatMessage(m Completion) {
	if len(m.Choices) > 0 {
		choice := m.Choices[0]
		lastmsg := mb.lastMessage()

		// ensure that the last message in the buffer has the same role as the choice
		// this will help us group delta messages by role
		if lastmsg.Role != "" && lastmsg.Role != choice.Delta.Role && choice.Delta.Role != "" {
			lastmsg = &Message{Role: choice.Delta.Role}
			*mb = append(*mb, lastmsg)
		}

		// ensure the first time we see a delta message, we set the role of the last message
		if choice.Delta.Role != "" {
			lastmsg.Role = choice.Delta.Role
		}

		lastmsg.Content += choice.Delta.Content
		lastmsg.FunctionCall = choice.Delta.FunctionCall
	}
}
