package store

import (
	"path/filepath"
	"testing"
)

func TestOrphanCountAndMeta(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "t.db"), ModeImport)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateSchema(); err != nil {
		t.Fatal(err)
	}
	x := db.SQL()
	mustExec(t, x, `INSERT INTO artist(id,gid,name,sort_name) VALUES (1,'g','A','A')`)
	mustExec(t, x, `INSERT INTO artist_credit(id,name,artist_count) VALUES (1,'A',1)`)
	// one good child (artist 1 exists) and one orphan (artist 999 missing)
	mustExec(t, x, `INSERT INTO artist_credit_name(artist_credit,position,artist,name) VALUES (1,0,1,'A')`)
	mustExec(t, x, `INSERT INTO artist_credit_name(artist_credit,position,artist,name) VALUES (1,1,999,'ghost')`)

	n, err := db.OrphanCount("artist_credit_name", "artist", "artist")
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("orphans = %d, want 1", n)
	}

	if err := db.WriteMeta(map[string]string{"schema_sequence": "31", "dump_date": "20260620"}); err != nil {
		t.Fatal(err)
	}
	var v string
	if err := x.QueryRow(`SELECT value FROM meta WHERE key='dump_date'`).Scan(&v); err != nil {
		t.Fatal(err)
	}
	if v != "20260620" {
		t.Fatalf("meta dump_date = %q", v)
	}
}
