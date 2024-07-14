package workers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// ParseSSE is a function that reads from a Server-Sent Events (SSE) stream.
// It reads line by line from the provided io.Reader (body), tries to parse each line as JSON,
// and sends the parsed data to the provided dataChan.
// If an error occurs during reading or parsing, it sends the error to the provided errChan.
// The function stops reading when it encounters an EOF (End of File) or an error.
func ParseSSE[T any](body io.Reader, dataChan chan T, errChan chan error) {
	reader := bufio.NewReader(body) // Create a new reader

	for {
		// Read until the next newline
		line, err := reader.ReadBytes('\n')
		if err != nil {
			// If the error is EOF, send it to errChan and return
			if errors.Is(err, io.EOF) {
				errChan <- err
				return
			}
			// If it's another error, wrap it with a custom message, send it to errChan, and return
			errChan <- fmt.Errorf("create chat completions read response body err: %w", err)
			return
		}

		// If the line is longer than 6 bytes, try to parse it as JSON
		if len(line) > 6 {

			// Check if the line indicates the end of the SSE stream
			if bytes.HasPrefix(line[6:], []byte("[DONE]")) {
				errChan <- io.EOF
				return
			}
			var data T
			// The actual data starts from the 7th byte, so we slice the line from the 6th index
			err = json.Unmarshal(line[6:], &data)
			// If an error occurred during parsing, wrap it with a custom message, send it to errChan, and return
			if err != nil {
				errChan <- fmt.Errorf("create chat completions unmarshal response body err: %w, data: %s", err, string(line[6:]))
				return
			}

			// If the line was successfully parsed, send the parsed data to dataChan
			dataChan <- data
		}
	}
}
