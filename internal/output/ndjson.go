package output

import (
	"encoding/json"
	"fmt"
	"io"
)

func PrintNDJSON(w io.Writer, items []json.RawMessage) error {
	for _, item := range items {
		if _, err := fmt.Fprintln(w, string(item)); err != nil {
			return err
		}
	}
	return nil
}

func WriteNDJSONLine(w io.Writer, item json.RawMessage) error {
	_, err := fmt.Fprintln(w, string(item))
	return err
}
