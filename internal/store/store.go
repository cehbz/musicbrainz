// internal/store/store.go
package store

import (
	_ "embed"
	"database/sql"
	"encoding/json"

	"github.com/cehbz/musicbrainz/internal/schema"
	_ "modernc.org/sqlite"
)

const PinnedSequence = 31

//go:embed schema_seq31.sql
var schemaSQL string

//go:embed indexes_seq31.sql
var indexesSQL string

//go:embed manifest_seq31.json
var manifestJSON []byte

type Mode int

const (
	ModeImport Mode = iota
	ModeServe
)

type DB struct {
	db       *sql.DB
	manifest schema.Manifest
}

func Open(path string, mode Mode) (*DB, error) {
	sdb, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	sdb.SetMaxOpenConns(1) // single writer; deterministic PRAGMA scope
	pragmas := []string{}
	switch mode {
	case ModeImport:
		pragmas = []string{
			"PRAGMA page_size=16384", "PRAGMA journal_mode=OFF", "PRAGMA synchronous=OFF",
			"PRAGMA foreign_keys=OFF", "PRAGMA cache_size=-2000000", "PRAGMA temp_store=MEMORY",
		}
	case ModeServe:
		pragmas = []string{"PRAGMA journal_mode=WAL", "PRAGMA synchronous=NORMAL"}
	}
	for _, p := range pragmas {
		if _, err := sdb.Exec(p); err != nil {
			sdb.Close()
			return nil, err
		}
	}
	var m schema.Manifest
	if err := json.Unmarshal(manifestJSON, &m); err != nil {
		sdb.Close()
		return nil, err
	}
	return &DB{db: sdb, manifest: m}, nil
}

func (d *DB) CreateSchema() error {
	_, err := d.db.Exec(schemaSQL)
	return err
}

func (d *DB) Manifest() schema.Manifest { return d.manifest }
func (d *DB) SQL() *sql.DB              { return d.db }
func (d *DB) Close() error              { return d.db.Close() }

// LoadEmbeddedManifest returns the build's pinned column manifest (parsed from the
// embedded manifest_seq31.json) without opening a database.
func LoadEmbeddedManifest() schema.Manifest {
	var m schema.Manifest
	_ = json.Unmarshal(manifestJSON, &m)
	return m
}
