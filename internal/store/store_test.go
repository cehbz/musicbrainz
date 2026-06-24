// internal/store/store_test.go
package store

import (
	"path/filepath"
	"testing"
)

func TestCreateSchema(t *testing.T) {
	p := filepath.Join(t.TempDir(), "t.db")
	db, err := Open(p, ModeImport)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateSchema(); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := db.SQL().QueryRow(
		`SELECT count(*) FROM sqlite_master WHERE type='table' AND name='artist'`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("artist table not created")
	}
	if PinnedSequence != 31 {
		t.Fatalf("PinnedSequence = %d, want 31", PinnedSequence)
	}
	if cols := db.Manifest().Tables["artist"]; len(cols) == 0 {
		t.Fatalf("manifest empty for artist")
	}

	// CRITICAL: confirm all 371 CREATE TABLE statements were executed
	var tableCount int
	if err := db.SQL().QueryRow(
		`SELECT count(*) FROM sqlite_master WHERE type='table'`).Scan(&tableCount); err != nil {
		t.Fatal(err)
	}
	const expectedTables = 371
	if tableCount != expectedTables {
		t.Fatalf("CreateSchema created %d tables, want %d; either modernc truncated the multi-statement Exec or the generated schema changed (regenerate via gen-schema)", tableCount, expectedTables)
	}
}
