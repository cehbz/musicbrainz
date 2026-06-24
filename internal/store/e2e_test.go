package store

import (
	"archive/tar"
	"bytes"
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// writeTarBz2 builds a .tar.bz2 file by shelling to the bzip2 CLI.
// It skips the test gracefully if bzip2 is unavailable.
func writeTarBz2(t *testing.T, path string, files map[string]string) {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for n, b := range files {
		tw.WriteHeader(&tar.Header{Name: n, Mode: 0o644, Size: int64(len(b))})
		tw.Write([]byte(b))
	}
	tw.Close()
	cmd := exec.Command("bzip2", "-c")
	cmd.Stdin = &buf
	out, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		t.Skipf("bzip2 CLI unavailable: %v", err)
	}
}

func TestRunImportEndToEnd(t *testing.T) {
	dumpDir := t.TempDir()
	dataRoot := t.TempDir()

	// Build a minimal artist row sized from the embedded manifest.
	cols := LoadEmbeddedManifest().Tables["artist"]
	fields := make([]string, len(cols))
	for i, c := range cols {
		switch c.Name {
		case "id":
			fields[i] = "1"
		case "gid":
			fields[i] = "33333333-3333-3333-3333-333333333333"
		case "name", "sort_name":
			fields[i] = "E2E"
		default:
			fields[i] = `\N`
		}
	}
	writeTarBz2(t, filepath.Join(dumpDir, "mbdump.tar.bz2"), map[string]string{
		"SCHEMA_SEQUENCE": "31\n",
		"TIMESTAMP":       "2026-06-20 00:20:52\n",
		"mbdump/artist":   joinTabs(fields) + "\n",
	})
	writeTarBz2(t, filepath.Join(dumpDir, "mbdump-derived.tar.bz2"), map[string]string{
		"SCHEMA_SEQUENCE": "31\n",
	})
	os.WriteFile(filepath.Join(dumpDir, "canonical_recording_redirect.csv"),
		[]byte("recording_mbid,canonical_recording_mbid,canonical_release_mbid\nr1,rc,rel\n"), 0o644)
	os.WriteFile(filepath.Join(dumpDir, "canonical_release_redirect.csv"),
		[]byte("release_mbid,canonical_release_mbid,release_group_mbid\nrel,rel,rg\n"), 0o644)
	os.WriteFile(filepath.Join(dumpDir, "canonical_musicbrainz_data.csv"),
		[]byte("id,artist_credit_id,artist_mbids,artist_credit_name,release_mbid,release_name,recording_mbid,recording_name,combined_lookup,score\n1,1,a,A,rel,R,r1,Rec,lk,10\n"), 0o644)

	rep, err := RunImport(dataRoot, dumpDir)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Counts["artist"] != 1 {
		t.Fatalf("artist count = %d, want 1", rep.Counts["artist"])
	}

	// symlink must exist and point at a real, queryable DB
	link := filepath.Join(dataRoot, "musicbrainz.db")
	if _, err := os.Stat(link); err != nil {
		t.Fatalf("symlink missing: %v", err)
	}
	db, err := sql.Open("sqlite", link)
	if err != nil {
		t.Fatalf("open symlink DB: %v", err)
	}
	defer db.Close()
	var name string
	if err := db.QueryRow(`SELECT name FROM artist WHERE id=1`).Scan(&name); err != nil {
		t.Fatalf("query artist: %v", err)
	}
	if name != "E2E" {
		t.Fatalf("artist name = %q, want E2E", name)
	}
}
