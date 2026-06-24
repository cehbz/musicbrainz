// internal/store/index_test.go
package store

import (
	"path/filepath"
	"testing"
)

func TestBuildFTSArtist(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "t.db"), ModeImport)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateSchema(); err != nil {
		t.Fatal(err)
	}
	// insert an artist with a diacritic in the name
	if _, err := db.SQL().Exec(`INSERT INTO artist(id, gid, name, sort_name) VALUES (1,'g','Antonín Dvořák','Dvořák, Antonín')`); err != nil {
		t.Fatal(err)
	}
	if err := db.BuildFTS(); err != nil {
		t.Fatal(err)
	}
	var id int
	// diacritic-insensitive match
	if err := db.SQL().QueryRow(`SELECT rowid FROM artist_fts WHERE artist_fts MATCH 'Antonin'`).Scan(&id); err != nil {
		t.Fatalf("fts match: %v", err)
	}
	if id != 1 {
		t.Fatalf("rowid = %d, want 1", id)
	}
}

func TestBuildFTSAllTablesValid(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "t.db"), ModeImport)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateSchema(); err != nil {
		t.Fatal(err)
	}
	// No data needed — just verify BuildFTS succeeds and all 6 tables exist.
	// A wrong column name in any SELECT would cause BuildFTS to return an error.
	if err := db.BuildFTS(); err != nil {
		t.Fatalf("BuildFTS: %v", err)
	}
	want := []string{
		"artist_fts", "label_fts", "work_fts",
		"release_group_fts", "release_fts", "recording_fts",
	}
	for _, name := range want {
		var n int
		if err := db.SQL().QueryRow(
			`SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?`, name,
		).Scan(&n); err != nil {
			t.Fatal(err)
		}
		if n != 1 {
			t.Errorf("FTS table %q not found in sqlite_master", name)
		}
	}
}

func TestBuildIndexes(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "t.db"), ModeImport)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateSchema(); err != nil {
		t.Fatal(err)
	}
	if err := db.BuildIndexes(); err != nil {
		t.Fatalf("BuildIndexes: %v", err)
	}
	var n int
	if err := db.SQL().QueryRow(
		`SELECT count(*) FROM sqlite_master WHERE type='index'`,
	).Scan(&n); err != nil {
		t.Fatal(err)
	}
	const expectedIndexes = 974
	if n != expectedIndexes {
		t.Fatalf("index count = %d, want %d; the generated index DDL changed or did not fully execute", n, expectedIndexes)
	}
}

func TestArtistNamesOrdered(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "t.db"), ModeImport)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateSchema(); err != nil {
		t.Fatal(err)
	}
	// One artist_credit (=1) with two names inserted OUT of position order.
	// artist_credit_name schema: artist_credit, position, artist, name, join_phrase.
	if _, err := db.SQL().Exec(`INSERT INTO artist(id, gid, name, sort_name) VALUES (1,'a','Alpha','Alpha'),(2,'b','Beta','Beta')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.SQL().Exec(
		`INSERT INTO artist_credit_name(artist_credit, position, artist, name, join_phrase) VALUES ` +
			`(1, 1, 2, 'Beta', ''), (1, 0, 1, 'Alpha', '')`); err != nil {
		t.Fatal(err)
	}
	var got string
	if err := db.SQL().QueryRow(
		`SELECT group_concat(name, ' ' ORDER BY position) FROM artist_credit_name WHERE artist_credit = 1`,
	).Scan(&got); err != nil {
		t.Fatalf("group_concat: %v", err)
	}
	if got != "Alpha Beta" {
		t.Fatalf("artist_names = %q, want %q (position order, not insert order)", got, "Alpha Beta")
	}
}
