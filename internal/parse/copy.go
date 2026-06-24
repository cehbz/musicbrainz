package parse

import (
	"bufio"
	"bytes"
	"io"
)

// CopyReader iterates PostgreSQL COPY text-format records.
type CopyReader struct {
	br *bufio.Reader
}

func NewCopyReader(r io.Reader) *CopyReader {
	return &CopyReader{br: bufio.NewReaderSize(r, 1<<20)}
}

// Next returns the next record. Each element is a string, or nil for SQL NULL.
func (c *CopyReader) Next() ([]any, error) {
	line, err := c.br.ReadBytes('\n')
	if len(line) == 0 && err != nil {
		return nil, err // io.EOF or read error
	}
	line = bytes.TrimRight(line, "\r\n")
	// Real tab bytes are field delimiters; in-data tabs are escaped as "\t".
	rawFields := bytes.Split(line, []byte{'\t'})
	row := make([]any, len(rawFields))
	for i, f := range rawFields {
		if len(f) == 2 && f[0] == '\\' && f[1] == 'N' {
			row[i] = nil // NULL
			continue
		}
		row[i] = unescape(f)
	}
	return row, nil
}

func unescape(f []byte) string {
	if bytes.IndexByte(f, '\\') < 0 {
		return string(f)
	}
	var b bytes.Buffer
	b.Grow(len(f))
	for i := 0; i < len(f); i++ {
		if f[i] != '\\' || i+1 >= len(f) {
			b.WriteByte(f[i])
			continue
		}
		i++
		switch f[i] {
		case 't':
			b.WriteByte('\t')
		case 'n':
			b.WriteByte('\n')
		case 'r':
			b.WriteByte('\r')
		case 'b':
			b.WriteByte('\b')
		case 'f':
			b.WriteByte('\f')
		case 'v':
			b.WriteByte('\v')
		case '\\':
			b.WriteByte('\\')
		default:
			b.WriteByte(f[i]) // backslash before unrecognized char: yields just that char (\q -> q)
		}
	}
	return b.String()
}
