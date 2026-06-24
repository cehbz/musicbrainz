// internal/parse/copy_test.go
package parse

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestCopyReaderBasic(t *testing.T) {
	// three fields: NULL, empty-string, escaped tab+newline+backslash
	in := "1\t\\N\t\\thello\\nworld\\\\\n2\tx\ty\n"
	r := NewCopyReader(strings.NewReader(in))

	row, err := r.Next()
	if err != nil {
		t.Fatalf("row1: %v", err)
	}
	if len(row) != 3 {
		t.Fatalf("row1 len = %d, want 3", len(row))
	}
	if row[0] != "1" {
		t.Fatalf("row1[0] = %v, want \"1\"", row[0])
	}
	if row[1] != nil {
		t.Fatalf("row1[1] = %v, want nil (NULL)", row[1])
	}
	if row[2] != "\thello\nworld\\" {
		t.Fatalf("row1[2] = %q, want unescaped tab/newline/backslash", row[2])
	}

	if _, err := r.Next(); err != nil {
		t.Fatalf("row2: %v", err)
	}
	if _, err := r.Next(); !errors.Is(err, io.EOF) {
		t.Fatalf("want io.EOF, got %v", err)
	}
}

func TestCopyReaderEmptyStringNotNull(t *testing.T) {
	r := NewCopyReader(strings.NewReader("\t\\N\n"))
	row, err := r.Next()
	if err != nil {
		t.Fatal(err)
	}
	if row[0] != "" {
		t.Fatalf("field0 = %v, want empty string", row[0])
	}
	if row[1] != nil {
		t.Fatalf("field1 = %v, want nil", row[1])
	}
}

func TestCopyReaderBackslashNInLongerFieldNotNull(t *testing.T) {
	// "x\N" is NOT NULL; only the exact two-byte field \N is NULL.
	r := NewCopyReader(strings.NewReader("x\\N\n"))
	row, err := r.Next()
	if err != nil {
		t.Fatal(err)
	}
	if len(row) != 1 {
		t.Fatalf("len = %d, want 1", len(row))
	}
	if row[0] != "xN" {
		t.Fatalf("field0 = %v, want string \"xN\" (not nil)", row[0])
	}
}

func TestCopyReaderEmptyInputEOF(t *testing.T) {
	r := NewCopyReader(strings.NewReader(""))
	if _, err := r.Next(); !errors.Is(err, io.EOF) {
		t.Fatalf("want io.EOF on empty input, got %v", err)
	}
}

func TestCopyReaderRemainingEscapes(t *testing.T) {
	r := NewCopyReader(strings.NewReader("\\r\\b\\f\\v\n"))
	row, err := r.Next()
	if err != nil {
		t.Fatal(err)
	}
	if row[0] != "\r\b\f\v" {
		t.Fatalf("field0 = %q, want \\r\\b\\f\\v control bytes", row[0])
	}
}

func TestCopyReaderLastLineNoTrailingNewline(t *testing.T) {
	r := NewCopyReader(strings.NewReader("a\tb"))
	row, err := r.Next()
	if err != nil {
		t.Fatalf("row1: %v", err)
	}
	if len(row) != 2 || row[0] != "a" || row[1] != "b" {
		t.Fatalf("row1 = %v, want [a b]", row)
	}
	if _, err := r.Next(); !errors.Is(err, io.EOF) {
		t.Fatalf("want io.EOF, got %v", err)
	}
}
