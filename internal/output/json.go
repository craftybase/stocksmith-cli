package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func PrintJSON(w io.Writer, body []byte) error {
	var v interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return fmt.Errorf("parse JSON: %w", err)
	}

	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("format JSON: %w", err)
	}

	_, err = fmt.Fprintln(w, string(out))
	return err
}

func ReadBody(r io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
