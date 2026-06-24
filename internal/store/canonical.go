// internal/store/canonical.go
package store

import (
	"fmt"
	"io"
	"strings"

	"github.com/cehbz/musicbrainz/internal/parse"
)

const canonicalDDL = `
CREATE TABLE canonical_musicbrainz_data (
  id INTEGER PRIMARY KEY, artist_credit_id INTEGER, artist_mbids TEXT,
  artist_credit_name TEXT, release_mbid TEXT, release_name TEXT,
  recording_mbid TEXT, recording_name TEXT, combined_lookup TEXT, score INTEGER
);
CREATE TABLE canonical_recording_redirect (
  recording_mbid TEXT, canonical_recording_mbid TEXT, canonical_release_mbid TEXT
);
CREATE TABLE canonical_release_redirect (
  release_mbid TEXT, canonical_release_mbid TEXT, release_group_mbid TEXT
);
`

// canonicalCols lists the destination columns per canonical table, in insert order.
var canonicalCols = map[string][]string{
	"canonical_musicbrainz_data":   {"id", "artist_credit_id", "artist_mbids", "artist_credit_name", "release_mbid", "release_name", "recording_mbid", "recording_name", "combined_lookup", "score"},
	"canonical_recording_redirect": {"recording_mbid", "canonical_recording_mbid", "canonical_release_mbid"},
	"canonical_release_redirect":   {"release_mbid", "canonical_release_mbid", "release_group_mbid"},
}

func (d *DB) CreateCanonicalSchema() error {
	_, err := d.db.Exec(canonicalDDL)
	return err
}

func (d *DB) LoadCanonical(table string, r io.Reader) (int, error) {
	cols, ok := canonicalCols[table]
	if !ok {
		return 0, fmt.Errorf("unknown canonical table %q", table)
	}
	cr, err := parse.NewCSVReader(r)
	if err != nil {
		return 0, err
	}
	ph := strings.TrimSuffix(strings.Repeat("?,", len(cols)), ",")
	stmtSQL := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", table, strings.Join(cols, ","), ph)

	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare(stmtSQL)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	n := 0
	for {
		row, err := cr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			stmt.Close()
			tx.Rollback()
			return n, err
		}
		vals := make([]any, len(cols))
		for i, c := range cols {
			if v, ok := row[c]; ok && v != "" {
				vals[i] = v
			}
		}
		if _, err := stmt.Exec(vals...); err != nil {
			stmt.Close()
			tx.Rollback()
			return n, err
		}
		n++
	}
	stmt.Close()
	return n, tx.Commit()
}
