package store

import (
	"archive/tar"
	"bytes"
	"path/filepath"
	"testing"
)

func tarWith(t *testing.T, files map[string]string) *bytes.Reader {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for n, b := range files {
		tw.WriteHeader(&tar.Header{Name: n, Mode: 0o644, Size: int64(len(b))})
		tw.Write([]byte(b))
	}
	tw.Close()
	return bytes.NewReader(buf.Bytes())
}

func joinTabs(f []string) string {
	out := ""
	for i, s := range f {
		if i > 0 {
			out += "\t"
		}
		out += s
	}
	return out
}

func TestGuardSequence(t *testing.T) {
	if err := GuardSequence(31); err != nil {
		t.Fatalf("31 should pass: %v", err)
	}
	if GuardSequence(30) == nil {
		t.Fatalf("30 should fail")
	}
}

func TestLoadTarballArtist(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "t.db"), ModeImport)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateSchema(); err != nil {
		t.Fatal(err)
	}
	// build a single artist row matching the manifest column count
	cols := db.Manifest().Tables["artist"]
	fields := make([]string, len(cols))
	for i, c := range cols {
		switch c.Name {
		case "id":
			fields[i] = "7"
		case "gid":
			fields[i] = "22222222-2222-2222-2222-222222222222"
		case "name", "sort_name":
			fields[i] = "Loaded"
		case "ended":
			fields[i] = "f"
		default:
			fields[i] = `\N`
		}
	}
	line := joinTabs(fields) + "\n"
	tb := tarWith(t, map[string]string{"mbdump/artist": line, "mbdump/unknown_table": "x\ty\n"})

	res, err := loadTarStream(db, tar.NewReader(tb))
	if err != nil {
		t.Fatal(err)
	}
	if res.Counts["artist"] != 1 {
		t.Fatalf("artist count = %d, want 1", res.Counts["artist"])
	}
	var name string
	if err := db.SQL().QueryRow(`SELECT name FROM artist WHERE id=7`).Scan(&name); err != nil {
		t.Fatal(err)
	}
	if name != "Loaded" {
		t.Fatalf("name = %q", name)
	}
	// unknown_table must be recorded in Skipped, not silently ignored
	found := false
	for _, s := range res.Skipped {
		if s == "mbdump/unknown_table" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("mbdump/unknown_table not in Skipped: %v", res.Skipped)
	}
}

func TestLoadTarballMalformedRow(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "t.db"), ModeImport)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateSchema(); err != nil {
		t.Fatal(err)
	}
	cols := db.Manifest().Tables["artist"]
	// good row
	fields := make([]string, len(cols))
	for i, c := range cols {
		switch c.Name {
		case "id":
			fields[i] = "9"
		case "gid":
			fields[i] = "33333333-3333-3333-3333-333333333333"
		case "name", "sort_name":
			fields[i] = "Good"
		case "ended":
			fields[i] = "f"
		default:
			fields[i] = `\N`
		}
	}
	goodLine := joinTabs(fields) + "\n"
	// short row — fewer fields than expected
	shortLine := "1\t2\n"

	content := goodLine + shortLine
	tb := tarWith(t, map[string]string{"mbdump/artist": content})

	res, err := loadTarStream(db, tar.NewReader(tb))
	if err != nil {
		t.Fatal(err)
	}
	if res.Counts["artist"] != 1 {
		t.Fatalf("inserted = %d, want 1", res.Counts["artist"])
	}
	if res.Malformed["artist"] != 1 {
		t.Fatalf("malformed = %d, want 1", res.Malformed["artist"])
	}
}
