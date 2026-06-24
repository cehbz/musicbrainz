package store

import (
	"archive/tar"
	"compress/bzip2"
	"fmt"
	"io"
	"strings"

	"github.com/cehbz/musicbrainz/internal/parse"
)

const insertBatch = 5000

// GuardSequence returns an error unless seq == PinnedSequence.
func GuardSequence(seq int) error {
	if seq != PinnedSequence {
		return fmt.Errorf("dump SCHEMA_SEQUENCE=%d but this build is pinned to %d; re-run gen-schema for the new schema", seq, PinnedSequence)
	}
	return nil
}

// LoadResult reports what a load did, including what it deliberately skipped —
// no silent drops (so an unexpectedly-dumped table or malformed input is visible).
type LoadResult struct {
	Counts    map[string]int // table -> rows inserted
	Skipped   []string       // mbdump/ entry names with no matching manifest table
	Malformed map[string]int // table -> rows skipped due to wrong column count
}

// loadTarStream walks an already-opened tar, routing mbdump/<table> files into the DB.
func loadTarStream(db *DB, tr *tar.Reader) (LoadResult, error) {
	res := LoadResult{Counts: map[string]int{}, Malformed: map[string]int{}}
	man := db.Manifest()
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return res, nil
		}
		if err != nil {
			return res, err
		}
		if h.Typeflag != tar.TypeReg {
			continue
		}
		name := strings.TrimPrefix(h.Name, "mbdump/")
		if name == h.Name {
			continue // not under mbdump/
		}
		if _, ok := man.Tables[name]; !ok {
			res.Skipped = append(res.Skipped, h.Name) // record, do NOT silently ignore
			continue
		}
		inserted, malformed, err := loadTable(db, name, tr)
		if err != nil {
			return res, fmt.Errorf("load %s: %w", name, err)
		}
		res.Counts[name] = inserted
		if malformed > 0 {
			res.Malformed[name] = malformed
		}
	}
}

// loadTable streams one COPY file into table. A row whose field count != the table's
// column count is counted as malformed and SKIPPED (not fatal — MusicBrainz data is
// clean but we don't abort a multi-hour import over one bad line). A parse error or a
// real insert error (DB failure) aborts.
//
// On an abort the Inserter needs no explicit Close: between batches it holds no open
// transaction, statement, or connection (flush is self-contained and self-rolls-back on
// failure), and the caller discards the partial DB on any load error.
func loadTable(db *DB, table string, r io.Reader) (inserted, malformed int, err error) {
	want := len(db.Manifest().Tables[table])
	ins, err := db.NewInserter(table, insertBatch)
	if err != nil {
		return 0, 0, err
	}
	cr := parse.NewCopyReader(r)
	for {
		row, e := cr.Next()
		if e == io.EOF {
			break
		}
		if e != nil {
			return inserted, malformed, e
		}
		if len(row) != want {
			malformed++
			continue
		}
		if e := ins.Add(row); e != nil {
			return inserted, malformed, e // real DB error -> abort
		}
		inserted++
	}
	if e := ins.Close(); e != nil {
		return inserted, malformed, e
	}
	return inserted, malformed, nil
}

// LoadTarballBz2 is the entry point: decompress a .tar.bz2 stream and load its mbdump/ tables.
func LoadTarballBz2(db *DB, r io.Reader) (LoadResult, error) {
	return loadTarStream(db, tar.NewReader(bzip2.NewReader(r)))
}
