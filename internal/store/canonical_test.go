// internal/store/canonical_test.go
package store

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadCanonicalRedirectChain(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "t.db"), ModeImport)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateCanonicalSchema(); err != nil {
		t.Fatal(err)
	}
	recCSV := "recording_mbid,canonical_recording_mbid,canonical_release_mbid\nrec1,recC,relX\n"
	relCSV := "release_mbid,canonical_release_mbid,release_group_mbid\nrelX,relX,rgZ\n"
	if _, err := db.LoadCanonical("canonical_recording_redirect", strings.NewReader(recCSV)); err != nil {
		t.Fatal(err)
	}
	if _, err := db.LoadCanonical("canonical_release_redirect", strings.NewReader(relCSV)); err != nil {
		t.Fatal(err)
	}
	// the resolution chain a consumer runs: recording -> release-group
	var rg string
	err = db.SQL().QueryRow(`
		SELECT rel.release_group_mbid
		FROM canonical_recording_redirect rec
		JOIN canonical_release_redirect rel ON rel.release_mbid = rec.canonical_release_mbid
		WHERE rec.recording_mbid = 'rec1'`).Scan(&rg)
	if err != nil {
		t.Fatal(err)
	}
	if rg != "rgZ" {
		t.Fatalf("release_group = %q, want rgZ", rg)
	}
}
