package store

import (
	"archive/tar"
	"compress/bzip2"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cehbz/musicbrainz/internal/dumps"
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

// RunImport orchestrates the full import pipeline:
// read meta → guard sequence → create schema → load tarballs → load canonical CSVs →
// build indexes + FTS → enrich Discogs → set serving PRAGMAs → write meta → atomically repoint symlink.
func RunImport(dataRoot, dumpDir string) (Report, error) {
	var rep Report
	if err := os.MkdirAll(dataRoot, 0o755); err != nil {
		return rep, err
	}
	seq, ts, err := dumps.ReadMeta(filepath.Join(dumpDir, "mbdump.tar.bz2"))
	if err != nil {
		return rep, fmt.Errorf("read meta: %w", err)
	}
	if err := GuardSequence(seq); err != nil {
		return rep, err
	}
	dumpDate := strings.ReplaceAll(strings.SplitN(ts, " ", 2)[0], "-", "")
	dbPath := filepath.Join(dataRoot, "musicbrainz-"+dumpDate+".db")
	_ = os.Remove(dbPath)

	db, err := Open(dbPath, ModeImport)
	if err != nil {
		return rep, err
	}
	fail := func(e error) (Report, error) { db.Close(); return rep, e }
	if err := db.CreateSchema(); err != nil {
		return fail(err)
	}
	if err := db.CreateCanonicalSchema(); err != nil { // BEFORE BuildIndexes
		return fail(err)
	}

	counts := map[string]int{}
	malformed := map[string]int{}
	var skipped []string
	for _, name := range []string{"mbdump.tar.bz2", "mbdump-derived.tar.bz2"} {
		f, err := os.Open(filepath.Join(dumpDir, name))
		if err != nil {
			return fail(err)
		}
		res, err := LoadTarballBz2(db, f)
		f.Close()
		if err != nil {
			return fail(err)
		}
		for k, v := range res.Counts {
			counts[k] += v
		}
		for k, v := range res.Malformed {
			malformed[k] += v
		}
		skipped = append(skipped, res.Skipped...)
	}

	var canonCounts map[string]int
	// pattern is static, so Glob cannot error
	canonMatches, _ := filepath.Glob(filepath.Join(dumpDir, "musicbrainz-canonical-dump-*.tar.zst"))
	if len(canonMatches) == 0 {
		log.Printf("WARNING: no canonical dump (musicbrainz-canonical-dump-*.tar.zst) in %s; canonical_* tables will be empty", dumpDir)
	} else {
		sort.Strings(canonMatches)
		cf, err := os.Open(canonMatches[len(canonMatches)-1]) // newest
		if err != nil {
			return fail(err)
		}
		canonCounts, err = LoadCanonicalTarZst(db, cf)
		cf.Close()
		if err != nil {
			return fail(err)
		}
	}

	if err := db.BuildIndexes(); err != nil {
		return fail(err)
	}
	if err := db.BuildFTS(); err != nil {
		return fail(err)
	}
	cov, err := db.EnrichDiscogs()
	if err != nil {
		return fail(err)
	}
	// switch to serving PRAGMAs (WAL persists in the file) + optimize
	for _, p := range []string{"PRAGMA journal_mode=WAL", "PRAGMA synchronous=NORMAL", "ANALYZE"} {
		if _, err := db.SQL().Exec(p); err != nil {
			return fail(err)
		}
	}
	if err := db.WriteMeta(map[string]string{
		"schema_sequence":  fmt.Sprint(seq),
		"dump_date":        dumpDate,
		"importer_version": "0.1.0",
	}); err != nil {
		return fail(err)
	}
	db.Close()

	// atomically repoint the musicbrainz.db symlink to the verified new file
	link := filepath.Join(dataRoot, "musicbrainz.db")
	tmp := link + ".tmp"
	_ = os.Remove(tmp)
	if err := os.Symlink(filepath.Base(dbPath), tmp); err != nil {
		return rep, err
	}
	if err := os.Rename(tmp, link); err != nil {
		return rep, err
	}

	rep.Counts, rep.Malformed, rep.Skipped, rep.DiscogsCoverage, rep.Canonical = counts, malformed, skipped, cov, canonCounts
	return rep, nil
}
