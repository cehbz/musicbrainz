package store

import (
	"database/sql"
	"fmt"
	"strings"
)

type Inserter struct {
	db        *sql.DB
	batch     int
	boolCol   []bool
	stmtSQL   string
	pending   [][]any
}

func (d *DB) NewInserter(table string, batch int) (*Inserter, error) {
	cols, ok := d.manifest.Tables[table]
	if !ok {
		return nil, fmt.Errorf("no manifest for table %q", table)
	}
	names := make([]string, len(cols))
	ph := make([]string, len(cols))
	boolCol := make([]bool, len(cols))
	for i, c := range cols {
		names[i] = c.Name
		ph[i] = "?"
		boolCol[i] = c.Bool
	}
	stmtSQL := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", table, strings.Join(names, ","), strings.Join(ph, ","))
	return &Inserter{db: d.db, batch: batch, boolCol: boolCol, stmtSQL: stmtSQL}, nil
}

// Add queues one row, coercing Bool-flagged columns ("t"/"f" -> 1/0) IN PLACE and
// committing each full batch. Because coercion mutates row, the caller must pass a
// slice it does not reuse; the COPY parser allocates a fresh []any per record, which
// satisfies this and avoids a per-row copy on the bulk-load hot path.
func (in *Inserter) Add(row []any) error {
	if len(row) != len(in.boolCol) {
		return fmt.Errorf("inserter: row has %d fields, want %d", len(row), len(in.boolCol))
	}
	for i := range row {
		if in.boolCol[i] {
			switch row[i] {
			case "t":
				row[i] = 1
			case "f":
				row[i] = 0
			}
		}
	}
	in.pending = append(in.pending, row)
	if len(in.pending) >= in.batch {
		return in.flush()
	}
	return nil
}

func (in *Inserter) flush() error {
	if len(in.pending) == 0 {
		return nil
	}
	tx, err := in.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(in.stmtSQL)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, r := range in.pending {
		if _, err := stmt.Exec(r...); err != nil {
			stmt.Close()
			tx.Rollback()
			return err
		}
	}
	stmt.Close()
	if err := tx.Commit(); err != nil {
		return err
	}
	in.pending = in.pending[:0]
	return nil
}

func (in *Inserter) Close() error { return in.flush() }
