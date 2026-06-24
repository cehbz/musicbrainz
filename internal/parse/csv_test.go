// internal/parse/csv_test.go
package parse

import (
	"io"
	"strings"
	"testing"
)

func TestCSVReaderHeaderMapped(t *testing.T) {
	in := "recording_mbid,canonical_recording_mbid,canonical_release_mbid\nA,B,C\n"
	r, err := NewCSVReader(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	row, err := r.Next()
	if err != nil {
		t.Fatal(err)
	}
	if row["canonical_release_mbid"] != "C" {
		t.Fatalf("row = %v", row)
	}
	if _, err := r.Next(); err != io.EOF {
		t.Fatalf("want EOF, got %v", err)
	}
}
