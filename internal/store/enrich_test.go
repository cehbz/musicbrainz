package store

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func mustExec(t *testing.T, db *sql.DB, q string, args ...any) {
	t.Helper()
	if _, err := db.Exec(q, args...); err != nil {
		t.Fatalf("exec %q: %v", q, err)
	}
}

func TestParseDiscogsID(t *testing.T) {
	cases := []struct {
		url, typ string
		id       int
		ok       bool
	}{
		{"https://www.discogs.com/artist/12345", "artist", 12345, true},
		{"https://www.discogs.com/artist/12345-Some-Name", "artist", 12345, true},
		{"http://www.discogs.com/label/678", "label", 678, true},
		{"https://www.discogs.com/master/42-Foo", "master", 42, true},
		{"https://www.discogs.com/release/99", "release", 99, true},
		{"https://example.com/artist/1", "artist", 0, false},
		{"https://www.discogs.com/artist/12345", "label", 0, false}, // wrong type
		{"https://www.discogs.com/fr/master/42", "master", 42, true},
		{"https://www.discogs.com/fr/artist/12345", "label", 0, false}, // locale + wrong type
		{"https://www.discogs.com/master/42", "master", 42, true},      // plain (no locale) still works
	}
	for _, c := range cases {
		id, ok := ParseDiscogsID(c.url, c.typ)
		if id != c.id || ok != c.ok {
			t.Errorf("ParseDiscogsID(%q,%q) = (%d,%v) want (%d,%v)", c.url, c.typ, id, ok, c.id, c.ok)
		}
	}
}

func TestEnrichDiscogsArtist(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "t.db"), ModeImport)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateSchema(); err != nil {
		t.Fatal(err)
	}
	x := db.SQL()
	mustExec(t, x, `INSERT INTO artist(id,gid,name,sort_name) VALUES (5,'g','A','A')`)
	mustExec(t, x, `INSERT INTO url(id,gid,url) VALUES (100,'u','https://www.discogs.com/artist/777-X')`)
	mustExec(t, x, `INSERT INTO l_artist_url(id,link,entity0,entity1) VALUES (1,1,5,100)`)

	cov, err := db.EnrichDiscogs()
	if err != nil {
		t.Fatal(err)
	}
	if cov["artist"] != 1 {
		t.Fatalf("artist coverage = %d, want 1", cov["artist"])
	}
	var got int
	if err := x.QueryRow(`SELECT discogs_artist_id FROM artist WHERE id=5`).Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got != 777 {
		t.Fatalf("discogs_artist_id = %d, want 777", got)
	}
}
