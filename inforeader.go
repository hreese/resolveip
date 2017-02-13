package resolveip

import (
	"io"
)

type InfoReader struct {
	msg []byte
	out io.Writer
}

func (ir *InfoReader) Read(p []byte) (n int, err error) {
	ir.out.Write(ir.msg)
	// always signal EOF and no data
	return 0, io.EOF
}

func NewInfoWriter(msg string, out io.Writer) *InfoReader {
	return &InfoReader{
		[]byte(msg),
		out,
	}
}
