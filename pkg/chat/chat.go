package chat

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/alexeyco/simpletable"
)

func Chat(url string, username string, token string, logLevel string) error {
	if url == "" {
		return fmt.Errorf("agent url is required")
	}

	ctx := context.Background()
	var history []Message

	if _, err := fmt.Fprintf(os.Stdout, "\nStart typing to chat with your assistant...\n%s: ", magenta(username)); err != nil {
		return fmt.Errorf("error writing to stdout: %w", err)
	}

	// Read full message from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		userMessage := Message{
			Role:    "user",
			Content: scanner.Text(),
		}
		history = append(history, userMessage)

		msgs, err := invokeAgent(ctx, url, token, history, logLevel)
		if err != nil {
			return fmt.Errorf(red("error creating message: %w"), err)
		}

		for _, msg := range msgs {
			fmt.Fprint(os.Stdout, &Output{
				Message:  msg,
				LogLevel: logLevel,
			})

			chatMsg := Message{
				Role:    msg.Role,
				Content: msg.Content,
			}
			if msg.FunctionCall != nil {
				chatMsg.FunctionCall = &ChatMessageFunctionCall{
					Name:      msg.FunctionCall.Name,
					Arguments: msg.FunctionCall.Arguments,
				}
			}

			history = append(history, chatMsg)
		}

		if _, err := fmt.Fprintf(os.Stdout, "%s: ", magenta(username)); err != nil {
			return fmt.Errorf("error writing to stdout: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading from stdin: %w", err)
	}

	return nil
}

func (o *Output) String() string {

	m := o.Message

	var msg strings.Builder
	if m.FunctionCall != nil {
		if shouldLog(o.LogLevel, LEVEL_DEBUG) {
			msg.WriteString(green("\nHuzzah! You successfully received a function call!\n"))

			table := simpletable.New()
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "Key"},
					{Align: simpletable.AlignCenter, Text: "Value"},
				},
			}
			cells := [][]*simpletable.Cell{
				{{Text: "role"}, {Text: m.Role}},
				{{Text: "name"}, {Text: m.FunctionCall.Name}},
				{{Text: "arguments"}, {Text: m.FunctionCall.Arguments}},
			}
			table.Body = &simpletable.Body{Cells: cells}

			table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Span: 2, Text: "Parsed function data"},
			}}

			table.SetStyle(simpletable.StyleUnicode)
			msg.WriteString(fmt.Sprintf("%s\n", green(table.String())))
		}

	} else {
		if m.Role != "" && m.Content != "" {
			msg.WriteString(fmt.Sprintf("%s: %s\n", cyan(m.Role), m.Content))

			if shouldLog(o.LogLevel, LEVEL_DEBUG) {
				msg.WriteString(fmt.Sprintf("\n%s\n", green("Huzzah! You successfully received a message!")))

				table := simpletable.New()
				table.Header = &simpletable.Header{
					Cells: []*simpletable.Cell{
						{Align: simpletable.AlignCenter, Text: "Role"},
						{Align: simpletable.AlignCenter, Text: "Content"},
					},
				}
				cells := [][]*simpletable.Cell{
					{{Text: m.Role}, {Text: fmt.Sprintf("[condensed] %.50s", m.Content)}},
				}
				table.Body = &simpletable.Body{Cells: cells}

				table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
					{Align: simpletable.AlignRight, Span: 2, Text: "Parsed message data"},
				}}

				table.SetStyle(simpletable.StyleUnicode)
				msg.WriteString(fmt.Sprintf("%s\n", green(table.String())))
			}
		}
	}

	if m.Confirmation != nil {
		if shouldLog(o.LogLevel, LEVEL_DEBUG) {
			msg.WriteString(green("\nHuzzah! You successfully received a confirmation!\n"))

			table := simpletable.New()
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignLeft, Text: "Key"},
					{Align: simpletable.AlignLeft, Text: "Value"},
				},
			}

			cells := [][]*simpletable.Cell{
				{{Text: "type"}, {Text: m.Confirmation.Type}},
				{{Text: "title"}, {Text: m.Confirmation.Title}},
				{{Text: "message"}, {Text: m.Confirmation.Message}},
				{{Text: "confirmation"}, {Text: fmt.Sprintf("%s", m.Confirmation.Confirmation)}},
			}
			table.Body = &simpletable.Body{Cells: cells}

			table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
				{Align: simpletable.AlignRight, Span: 2, Text: "Parsed confirmation data"},
			}}

			table.SetStyle(simpletable.StyleUnicode)
			msg.WriteString(fmt.Sprintf("%s\n", green(table.String())))
		}
		msg.WriteString(cyan(fmt.Sprintf("\n%s\n  %s\nReply: [y/N]\n", m.Confirmation.Title, m.Confirmation.Message)))
	}

	if len(m.References) > 0 {
		// When debug mode is turned off, the refrerences are not explicitly displayed
		if shouldLog(o.LogLevel, LEVEL_DEBUG) {
			msg.WriteString(green("\nHuzzah! You successfully received some references!\n"))

			table := simpletable.New()
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignLeft, Text: "index"},
					{Align: simpletable.AlignLeft, Text: "id"},
					{Align: simpletable.AlignLeft, Text: "type"},
					{Align: simpletable.AlignLeft, Text: "data"},
					{Align: simpletable.AlignLeft, Text: "display_icon"},
					{Align: simpletable.AlignLeft, Text: "display_name"},
					{Align: simpletable.AlignLeft, Text: "display_url"},
				},
			}

			var cells [][]*simpletable.Cell
			for i, reference := range m.References {
				cells = append(cells, []*simpletable.Cell{
					{Text: fmt.Sprintf("%d", i)},
					{Text: reference.ID},
					{Text: reference.Type},
					{Text: fmt.Sprintf("[condensed] %.20s", reference.Data)},
					{Text: reference.Metadata.DisplayIcon},
					{Text: reference.Metadata.DisplayName},
					{Text: reference.Metadata.DisplayURL},
				})

			}
			table.Body = &simpletable.Body{Cells: cells}

			table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
				{Align: simpletable.AlignRight, Span: 7, Text: "Parsed references data"},
			}}

			table.SetStyle(simpletable.StyleUnicode)
			msg.WriteString(fmt.Sprintf("%s\n", green(table.String())))
		}

		for i, reference := range m.References {
			msg.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, reference.ID, reference.Metadata.DisplayName))
		}
	}

	if len(m.Errors) > 0 {
		table := simpletable.New()

		if shouldLog(o.LogLevel, LEVEL_DEBUG) {
			msg.WriteString(green("\nHuzzah! You successfully received some errors!\n"))

			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignLeft, Text: "index"},
					{Align: simpletable.AlignLeft, Text: "message"},
					{Align: simpletable.AlignLeft, Text: "type"},
					{Align: simpletable.AlignLeft, Text: "code"},
					{Align: simpletable.AlignLeft, Text: "identifier"},
				},
			}

			var cells [][]*simpletable.Cell
			for i, error := range m.Errors {
				cells = append(cells, []*simpletable.Cell{
					{Text: fmt.Sprintf("%d", i)},
					{Text: error.Message},
					{Text: error.Type},
					{Text: error.Code},
					{Text: error.Identifier},
				})
			}
			table.Body = &simpletable.Body{Cells: cells}

			table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
				{Align: simpletable.AlignRight, Span: 5, Text: "Parsed error data"},
			}}

			table.SetStyle(simpletable.StyleUnicode)
			msg.WriteString(fmt.Sprintf("%s\n", green(table.String())))
		}

		for i, error := range m.Errors {
			msg.WriteString(fmt.Sprintf("%d. %s error: %s\n", i+1, error.Type, error.Message))
		}
	}

	return msg.String()
}
