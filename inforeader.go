// Package resolveip provides utilities for finding IP addresses in text
// and resolving them to hostnames using the system's DNS resolver.
package resolveip

import (
	"fmt"
	"io"
)

// InfoReader is an [io.Reader] that writes a one-time informational message
// to a separate writer and immediately returns [io.EOF]. It is used to display
// a prompt (e.g., on Windows) before reading from the next reader in a
// [io.MultiReader] chain.
type InfoReader struct {
	msg []byte
	out io.Writer
}

// Read writes the informational message to the configured writer and
// returns 0 bytes read with [io.EOF], signaling that no data is available
// from this reader.
func (ir *InfoReader) Read(_ []byte) (n int, err error) {
	if _, werr := ir.out.Write(ir.msg); werr != nil {
		fmt.Fprintf(ir.out, "warning: failed to write info message: %v\n", werr)
	}
	return 0, io.EOF
}

// NewInfoReader creates a new [InfoReader] that writes msg to out when read.
func NewInfoReader(msg string, out io.Writer) *InfoReader {
	return &InfoReader{
		msg: []byte(msg),
		out: out,
	}
}
