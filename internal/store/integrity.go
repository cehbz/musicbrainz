package store

import "fmt"

type Report struct {
	Counts          map[string]int
	Orphans         map[string]int
	DiscogsCoverage map[string]int
	Canonical       map[string]int // per-table canonical row counts; empty if canonical absent
	Skipped         []string       // mbdump entries with no matching table
	Malformed       map[string]int // table -> rows skipped (bad column count)
}

func (d *DB) TableCounts(tables []string) (map[string]int, error) {
	out := map[string]int{}
	for _, t := range tables {
		var n int
		if err := d.db.QueryRow(fmt.Sprintf("SELECT count(*) FROM %s", t)).Scan(&n); err != nil {
			return out, err
		}
		out[t] = n
	}
	return out, nil
}

func (d *DB) OrphanCount(childTable, childCol, parentTable string) (int, error) {
	q := fmt.Sprintf(
		`SELECT count(*) FROM %s c WHERE c.%s IS NOT NULL AND NOT EXISTS (SELECT 1 FROM %s p WHERE p.id = c.%s)`,
		childTable, childCol, parentTable, childCol)
	var n int
	err := d.db.QueryRow(q).Scan(&n)
	return n, err
}

// orphanChecks are the FK relationships the import verifies (all FK columns are indexed,
// so each is an indexed anti-join). Keyed "<child>.<col>" in the report.
var orphanChecks = []struct{ child, col, parent string }{
	{"artist_credit_name", "artist", "artist"},
	{"artist_credit_name", "artist_credit", "artist_credit"},
	{"isrc", "recording", "recording"},
	{"release", "release_group", "release_group"},
	{"release", "artist_credit", "artist_credit"},
	{"recording", "artist_credit", "artist_credit"},
	{"medium", "release", "release"},
	{"track", "medium", "medium"},
}

// OrphanPass runs the curated orphan checks, returning "<child>.<col>" -> dangling count.
func (d *DB) OrphanPass() (map[string]int, error) {
	out := map[string]int{}
	for _, c := range orphanChecks {
		n, err := d.OrphanCount(c.child, c.col, c.parent)
		if err != nil {
			return out, fmt.Errorf("orphan check %s.%s: %w", c.child, c.col, err)
		}
		out[c.child+"."+c.col] = n
	}
	return out, nil
}
