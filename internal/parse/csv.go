// internal/parse/csv.go
package parse

import (
	"encoding/csv"
	"io"
)

type CSVReader struct {
	r      *csv.Reader
	header []string
}

func NewCSVReader(r io.Reader) (*CSVReader, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.ReuseRecord = true
	head, err := cr.Read()
	if err != nil {
		return nil, err
	}
	header := make([]string, len(head))
	copy(header, head)
	return &CSVReader{r: cr, header: header}, nil
}

func (c *CSVReader) Next() (map[string]string, error) {
	rec, err := c.r.Read()
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, len(c.header))
	for i, h := range c.header {
		if i < len(rec) {
			m[h] = rec[i]
		}
	}
	return m, nil
}
