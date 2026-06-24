package store

import (
	"path/filepath"
	"testing"
)

func TestInserterBoolCoercionAndNull(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "t.db"), ModeImport)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateSchema(); err != nil {
		t.Fatal(err)
	}
	ins, err := db.NewInserter("artist", 2)
	if err != nil {
		t.Fatal(err)
	}
	// manifest column order for artist starts: id, gid, name, ... , ended(bool), ...
	cols := db.Manifest().Tables["artist"]
	row := make([]any, len(cols))
	for i, c := range cols {
		switch c.Name {
		case "id":
			row[i] = "42"
		case "gid":
			row[i] = "11111111-1111-1111-1111-111111111111"
		case "name", "sort_name":
			row[i] = "Test"
		case "ended":
			row[i] = "t"
		default:
			row[i] = nil
		}
	}
	if err := ins.Add(row); err != nil {
		t.Fatal(err)
	}
	if err := ins.Close(); err != nil {
		t.Fatal(err)
	}
	var ended int
	if err := db.SQL().QueryRow(`SELECT ended FROM artist WHERE id=42`).Scan(&ended); err != nil {
		t.Fatal(err)
	}
	if ended != 1 {
		t.Fatalf("ended = %d, want 1 (coerced from 't')", ended)
	}
}

func TestInserterColumnCountMismatch(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "t.db"), ModeImport)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateSchema(); err != nil {
		t.Fatal(err)
	}
	ins, err := db.NewInserter("artist", 2)
	if err != nil {
		t.Fatal(err)
	}
	// A row of the wrong length must return an error, not panic.
	if err := ins.Add([]any{"1", "g"}); err == nil {
		t.Fatal("Add with wrong-length row: got nil error, want non-nil")
	}
}
